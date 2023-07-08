package service

import (
	"context"
	"errors"

	"github.com/wonderivan/logger"
	nwv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ingress struct{}

var Ingress ingress

type IngressCreate struct {
	Name      string                 `json:"name"`
	Namespace string                 `json:"namespace"`
	Label     map[string]string      `json:"label"`
	Hosts     map[string][]*HttpPath `json:"hosts"`
}

type HttpPath struct {
	Path        string `json:"path"`
	PathType    string `json:"pathtype"`
	ServiceName string `json:"servicename"`
	ServicePort int32  `json:"serviceport"`
}

func (i *ingress) CreateIngress(data *IngressCreate) (err error) {

	var ingressRules []nwv1.IngressRule
	var httpIngressPaths []nwv1.HTTPIngressPath

	// 拼接ingress
	ingress := &nwv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
			Labels:    data.Label,
		},
		Status: nwv1.IngressStatus{},
	}

	// 遍历host
	for host, paths := range data.Hosts {
		ir := nwv1.IngressRule{
			Host: host,
			IngressRuleValue: nwv1.IngressRuleValue{
				HTTP: &nwv1.HTTPIngressRuleValue{
					Paths: nil,
				},
			},
		}

		// 遍历host下的每个rule
		for _, path := range paths {
			hip := nwv1.HTTPIngressPath{
				Path:     path.Path,
				PathType: (*nwv1.PathType)(&path.PathType),
				Backend: nwv1.IngressBackend{
					Service: &nwv1.IngressServiceBackend{
						Name: path.ServiceName,
						Port: nwv1.ServiceBackendPort{
							Number: path.ServicePort,
						},
					},
				},
			}

			httpIngressPaths = append(httpIngressPaths, hip)
		}

		// 填充ingress rule
		ir.IngressRuleValue.HTTP.Paths = httpIngressPaths
		ingressRules = append(ingressRules, ir)

	}
	ingress.Spec.Rules = ingressRules

	// 创建ingress
	_, err = K8s.ClientSet.NetworkingV1().Ingresses(data.Namespace).Create(context.TODO(), ingress, metav1.CreateOptions{})
	if err != nil {
		logger.Error(errors.New("创建ingress失败：" + err.Error()))
		return errors.New("创建ingress失败：" + err.Error())
	}

	return nil
}
