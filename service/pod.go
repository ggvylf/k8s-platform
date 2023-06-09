package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"k8s-platform/config"
	"time"

	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type pod struct{}

var Pod pod

// 从k8s获取的pod信息
type PodResp struct {
	Items []corev1.Pod
	Total int
}

// 指定namespace下的pod个数
type PodsNp struct {
	Namespace string `json:"namespace"`
	PodNum    int    `json:"pod_num"`
}

// 实现DataCell的接口
type podCell corev1.Pod

func (p podCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time

}

func (p podCell) GetName() string {
	return p.Name

}

// 类型转换
// corev1.Pod转换成DataCell
func (p *pod) toCells(pods []corev1.Pod) []DataCell {
	cells := make([]DataCell, len(pods))
	for i := range pods {
		cells[i] = podCell(pods[i])
	}
	return cells

}

// DataCell转换成corev1.Pod
func (p *pod) toPods(cells []DataCell) []corev1.Pod {
	pods := make([]corev1.Pod, len(cells))
	for i := range cells {
		pods[i] = corev1.Pod(cells[i].(podCell))
	}

	return pods
}

// 获取pod列表，支持过滤 排序 分页
func (p *pod) GetPods(filterName, namespace string, limit, page int) (podresp *PodResp, err error) {
	podList, err := K8s.ClientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		logger.Error(errors.New("获取pod列表失败：" + err.Error()))
		return nil, errors.New("获取pod列表失败：" + err.Error())
	}

	//实例化dataSelector对象
	ds := &dataSelector{
		GenericDataList: p.toCells(podList.Items),
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
	pods := p.toPods(data.GenericDataList)

	return &PodResp{
		Items: pods,
		Total: total,
	}, nil

}

// pod详情
func (p *pod) GetPodDetail(podName, namespace string) (pod *corev1.Pod, err error) {
	pod, err = K8s.ClientSet.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取pod详情失败：" + err.Error()))
		return nil, errors.New("获取pod详情失败：" + err.Error())
	}

	return pod, nil
}

// 删除pod
func (p *pod) DeletePod(podName, namespace string) (err error) {
	err = K8s.ClientSet.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除pod失败：" + err.Error()))
		return errors.New("删除pod失败：" + err.Error())
	}

	return nil
}

// 更新pod
func (p *pod) UpdatePod(podName, namespace, updateinfo string) (err error) {
	pod := &corev1.Pod{}
	err = json.Unmarshal([]byte(updateinfo), pod)
	if err != nil {
		logger.Error(errors.New("解析pod更新数据失败：" + err.Error()))
		return errors.New("解析pod更新数据失败：" + err.Error())
	}

	_, err = K8s.ClientSet.CoreV1().Pods(namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新pod失败：" + err.Error()))
		return errors.New("更新pod失败：" + err.Error())
	}

	return nil
}

// 获取pod中的container的名称列表
func (p *pod) GetContainerNameList(podName, namespace string) (cnames []string, err error) {

	pod, err := p.GetPodDetail(podName, namespace)

	if err != nil {
		return nil, err
	}

	for _, c := range pod.Spec.Containers {
		cnames = append(cnames, c.Name)
	}
	return cnames, nil
}

// 获取container的日志
func (p *pod) GetPodLog(containerName, podName, namespace string) (log string, err error) {
	// 日志行数
	lineLimint := int64(config.PodLogTailLline)
	logopt := &corev1.PodLogOptions{
		Container: containerName,
		TailLines: &lineLimint,
	}

	// 创建*rest.Request的实例
	logreq := K8s.ClientSet.CoreV1().Pods(namespace).GetLogs(podName, logopt)

	// 使用stream的方式返回响应
	podlogs, err := logreq.Stream(context.TODO())
	if err != nil {
		logger.Error(errors.New("获取pod log失败：" + err.Error()))
		return "", errors.New("获取pod log失败：" + err.Error())
	}
	// io类型，用完要关闭
	defer podlogs.Close()

	// 用bytes.Buffer来缓冲log stream
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podlogs)
	if err != nil {
		logger.Error(errors.New("复制pod log失败：" + err.Error()))
		return "", errors.New("复制pod log失败：" + err.Error())
	}

	return buf.String(), nil

}

// 获取指定ns下的pod数量
func (p *pod) GetPodSumPerNp() (podsNps []*PodsNp, err error) {

	// 获取所有ns
	nslist, err := K8s.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取ns列表失败：" + err.Error()))
		return nil, errors.New("获取ns列表失败：" + err.Error())
	}

	// 遍历ns下的pod数量
	for _, ns := range nslist.Items {
		podlist, err := K8s.ClientSet.CoreV1().Pods(ns.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			logger.Error(errors.New("获取pod列表失败：" + err.Error()))
			return nil, errors.New("获取pod列表失败：" + err.Error())
		}

		podNp := &PodsNp{
			Namespace: ns.Name,
			PodNum:    len(podlist.Items),
		}

		podsNps = append(podsNps, podNp)
	}

	return podsNps, nil
}
