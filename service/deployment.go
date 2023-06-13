package service

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
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
	Namespace     string `json:"namespace"`
	DeploymentNum int    `json:"deployment_num"`
}

// 创建deployment需要的参数
type DeployCreate struct {
	Name          string            `json:"name"`
	Namespace     string            `json:"namespace"`
	Replicas      int32             `json:"replicas"`
	Image         string            `json:"image"`
	Label         map[string]string `json:"label"`
	Cpu           string            `json:"cpu"`
	Memory        string            `json:"memory"`
	ContainerPort int32             `json:"container_port"`
	HealthCheck   bool              `json:"health_check"`
	HealthPath    string            `json:"health_path"`
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
func (d *deployment) GetDeploymentDetail(name, namespace string) (deploy *appsv1.Deployment, err error) {

	deploy, err = K8s.ClientSet.AppsV1().Deployments(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取deployment详情失败：" + err.Error()))
		return nil, errors.New("获取deployment详情失败：" + err.Error())
	}
	return deploy, nil
}

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
func (d *deployment) CreateDeployment(data *DeployCreate) (err error) {

	// 填充deployment数据
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
			Labels:    data.Label,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &data.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: data.Label,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   data.Name,
					Labels: data.Label,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  data.Name,
							Image: data.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: data.ContainerPort,
								},
							},
						},
					},
				},
			},
		},
		Status: appsv1.DeploymentStatus{},
	}

	// 健康检查相关
	if data.HealthCheck {
		deploy.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: data.HealthPath,
					Port: intstr.IntOrString{
						Type:   0,
						IntVal: data.ContainerPort,
					},
				},
			},
			InitialDelaySeconds: 5,
			TimeoutSeconds:      5,
			PeriodSeconds:       5,
		}
		deploy.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: data.HealthPath,
					Port: intstr.IntOrString{
						Type:   0,
						IntVal: data.ContainerPort,
					},
				},
			},
			InitialDelaySeconds: 5,
			TimeoutSeconds:      5,
			PeriodSeconds:       5,
		}
	}

	// 资源限制相关
	deploy.Spec.Template.Spec.Containers[0].Resources.Limits = map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU:    resource.MustParse(data.Cpu),
		corev1.ResourceMemory: resource.MustParse(data.Memory),
	}

	deploy.Spec.Template.Spec.Containers[0].Resources.Requests = map[corev1.ResourceName]resource.Quantity{
		corev1.ResourceCPU:    resource.MustParse(data.Cpu),
		corev1.ResourceMemory: resource.MustParse(data.Memory),
	}

	// 创建deployment
	_, err = K8s.ClientSet.AppsV1().Deployments(data.Namespace).Create(context.TODO(), deploy, metav1.CreateOptions{})
	if err != nil {
		logger.Error(errors.New("创建deployment失败：" + err.Error()))
		return errors.New("创建deployment失败：" + err.Error())
	}

	return nil

}

// 删除deployment
func (d *deployment) DeleteDeployment(name, namespace string) (err error) {
	err = K8s.ClientSet.AppsV1().Deployments(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除deployment失败：" + err.Error()))
		return errors.New("删除deployment失败：" + err.Error())
	}

	return nil
}

// 更新deployment

func (d *deployment) UpdateDeployment(namespace, data string) (err error) {

	deploy := &appsv1.Deployment{}

	err = json.Unmarshal([]byte(data), deploy)
	if err != nil {
		logger.Error(errors.New("反序列化deployment失败：" + err.Error()))
		return errors.New("反序列化deployment失败：" + err.Error())
	}

	_, err = K8s.ClientSet.AppsV1().Deployments(namespace).Update(context.TODO(), deploy, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("删除deployment失败：" + err.Error()))
		return errors.New("删除deployment失败：" + err.Error())
	}

	return nil
}

// 重启deployment
// 修改无关的值 会触发deployment的重启
func (d *deployment) RestartDeployment(name, namespace string) (err error) {

	// 拼接patch数据
	patchData := map[string]interface{}{
		"spec": map[string]interface{}{
			"template:": map[string]interface{}{
				"spec": map[string]interface{}{
					"containers": []map[string]interface{}{
						{
							"name": name,
							"env": []map[string]interface{}{
								{
									"name":  "RESTART_",
									"value": strconv.FormatInt(time.Now().Unix(), 10),
								},
							},
						},
					},
				},
			},
		},
	}

	// 序列化
	patchBytes, err := json.Marshal(patchData)
	if err != nil {
		logger.Error(errors.New("序列化patch失败：" + err.Error()))
		return errors.New("序列化patch失败：" + err.Error())
	}

	// deployment的patch操作
	_, err = K8s.ClientSet.AppsV1().Deployments(namespace).Patch(context.TODO(), name, "application/strategic-merge-patch+json", patchBytes, metav1.PatchOptions{})
	if err != nil {
		logger.Error(errors.New("重启deployment失败：" + err.Error()))
		return errors.New("重启deployment失败失败：" + err.Error())
	}

	return nil

}

// 获取对应ns下的deployment数量
func (p *pod) GetDeploymentSumPerNp() (deployNps []*deploymentsNp, err error) {

	// 获取所有ns
	nslist, err := K8s.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取ns列表失败：" + err.Error()))
		return nil, errors.New("获取ns列表失败：" + err.Error())
	}

	// 遍历ns下的pod数量
	for _, ns := range nslist.Items {
		deploylist, err := K8s.ClientSet.AppsV1().Deployments(ns.Namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			logger.Error(errors.New("获取deployment列表失败：" + err.Error()))
			return nil, errors.New("获取deployment列表失败：" + err.Error())
		}

		deployNp := &deploymentsNp{
			Namespace:     ns.Name,
			DeploymentNum: len(deploylist.Items),
		}

		deployNps = append(deployNps, deployNp)
	}

	return deployNps, nil
}
