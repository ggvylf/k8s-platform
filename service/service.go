package service

import (
	"context"
	"errors"

	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type service struct{}

var Service service

type ServiceCreate struct {
	Name          string            `json:"name"`
	Namespace     string            `json:"namespace"`
	Type          string            `json:"type"`
	ContainerPort int32             `json:"container_port"`
	Port          int32             `json:"port"`
	NodePort      int32             `json:"nodeport"`
	Label         map[string]string `json:"label"`
}

func (s *service) CreateService(data *ServiceCreate) (err error) {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
			Labels:    data.Label,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceType(data.Type),
			Ports: []corev1.ServicePort{
				{
					Name:     "http",
					Port:     data.Port,
					Protocol: "TCP",
					TargetPort: intstr.IntOrString{
						Type:   0,
						IntVal: data.ContainerPort,
					},
				},
			},
			Selector: data.Label,
		},
	}

	// 判断是否有NodePort，默认是ClusterIP
	// 不明确指定Port，会随即分配
	if data.NodePort != 0 && data.Type == "NodePort" {
		service.Spec.Ports[0].NodePort = data.NodePort
	}

	_, err = K8s.ClientSet.CoreV1().Services(data.Namespace).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		logger.Error(errors.New("创建service失败：" + err.Error()))
		return errors.New("创建service失败：" + err.Error())
	}

	return nil

}
