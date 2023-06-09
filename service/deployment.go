package service

import (
	"context"
	"errors"
	"time"

	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type deployment struct{}

var Deployment deployment

// 从k8s获取的pod信息
type deploymentResp struct {
	Items []appsv1.Deployment
	Total int
}

// 指定namespace下的pod个数
type deploymentsNp struct {
	Namespace string `json:"namespace"`
	PodNum    int    `json:"deployment_num"`
}

// 实现DataCell的接口
type deplopymentCell appsv1.Deployment

func (d deplopymentCell) GetCreation() time.Time {
	return d.CreationTimestamp.Time

}

func (d deplopymentCell) GetName() string {
	return d.Name

}

// 类型转换
// appsv1.Deployment转换成DataCell
func (d *deployment) toCells(deployments []appsv1.Deployment) []DataCell {
	cells := make([]DataCell, len(deployments))
	for i := range deployments {
		cells[i] = deplopymentCell(deployments[i])
	}
	return cells

}

// DataCell转换成appsv1.Deployment
func (p *deployment) toDeployments(cells []DataCell) []appsv1.Deployment {
	deployments := make([]appsv1.Deployment, len(cells))
	for i := range cells {
		deployments[i] = appsv1.Deployment(cells[i].(deplopymentCell))
	}

	return deployments
}

// TODO
// 获取deployment列表，支持过滤 排序 分页
func (p *deployment) Getdeployments(filterName, namespace string, limit, page int) (deploymentresp *deploymentResp, err error) {
	deploymentList, err := K8s.ClientSet.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取deployment列表失败：" + err.Error()))
		return nil, errors.New("获取deployment列表失败：" + err.Error())
	}

	//实例化dataSelector对象
	ds := &dataSelector{
		GenericDataList: p.toCells(deploymentList.Items),
		dataSelectQuery: &DataSelectQuery{
			FilterQuery: &FilterQuery{
				Name: filterName,
			},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}

	// 先过滤
	filtered := ds.Filter()
	total := len(filtered.GenericDataList)

	// 再排序和分页
	data := filtered.Sort().Paginate()

	// 把结果做类型转换
	deployments := p.toDeployments(data.GenericDataList)

	return &deploymentResp{
		Items: deployments,
		Total: total,
	}, nil

}

// 获取deployment详情
// 修改deployment副本个数
func (d *deployment) ScaleDeployment(deploymentName, namespace string, scaleNum int) (replica int32, err error) {

	//获取当前的scale状态
	scale, err := K8s.ClientSet.AppsV1().Deployments(namespace).GetScale(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取deployment副本个数失败：" + err.Error()))
		return 0, errors.New("获取deployment副本个数失败：" + err.Error())
	}

	// 修改副本数
	scale.Spec.Replicas = int32(scaleNum)

	// 应用更新
	// 这里使用UpdateScale，没有使用ApplyScale
	newscale, err := K8s.ClientSet.AppsV1().Deployments(namespace).UpdateScale(context.TODO(), deploymentName, scale, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新deployment副本个数失败：" + err.Error()))
		return 0, errors.New("更新deployment副本个数失败：" + err.Error())
	}

	return newscale.Spec.Replicas, nil
}

// 创建deployment
func (d *deployemnt) CreateDeploymen(data *deployment) (err error) {

	// 填充数据
	deploy := &appsv1.Deployment{}

	// 创建
	_, err = K8s.ClientSet.AppsV1().Deployments(data.namespace).Create(context.TODO(), deploy, metav1.CreateOptions{})
	if err != nil {
		logger.Error(errors.New("创建deployment失败：" + err.Error()))
		return errors.New("创建deployment失败：" + err.Error())
	}

	return nil

}

// 删除deployment
// 更新deployment
// 重启deployment
// 获取对应ns下的deployment数量
