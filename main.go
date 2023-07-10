package main

import (
	"k8s-platform/config"
	"k8s-platform/controller"
	"k8s-platform/db"
	"k8s-platform/service"

	"github.com/gin-gonic/gin"
)

func main() {

	// 初始化db
	db.Init()

	// 初始化k8s的clientset
	service.K8s.Init()

	// 初始化路由
	r := gin.Default()
	controller.Router.InitApiRouter(r)
	r.Run(config.ListenAddr)

	//关闭db连接池
	db.Close()
}
