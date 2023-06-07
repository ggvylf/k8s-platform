package service

import (
	"time"

	corev1 "k8s.io/api/core/v1"
)

type pod struct{}

var Pod pod

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
