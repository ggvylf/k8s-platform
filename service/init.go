package service

import (
	"k8s-platform/config"

	"github.com/wonderivan/logger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type k8s struct {
	ClientSet *kubernetes.Clientset
}

var K8s k8s

func (k *k8s) Init() {
	config, err := clientcmd.BuildConfigFromFlags("", config.Kubeconfig)
	if err != nil {
		logger.Error("创建kubeconfig失败：", err.Error())
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Error("创建k8s clientSet失败：", err.Error())
	}
	logger.Info("创建k8s clientSet成功")

	k.ClientSet = clientSet
}
