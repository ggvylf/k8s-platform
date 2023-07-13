package main

import (
	"k8s-platform/config"
	"k8s-platform/controller"
	"k8s-platform/db"
	"k8s-platform/middle"
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

	// 加载中间件

	r.Use(middle.Cors())

	// 登录相关接口
	controller.Router.LoginRouter(r)

	r.Use(middle.JWTAuth())

	// 加载业务路由
	controller.Router.InitApiRouter(r)

	r.Run(config.ListenAddr)

	//关闭db连接池
	db.Close()
}
