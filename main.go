package main

import (
	"k8s-platform/config"
	"k8s-platform/controller"
	"k8s-platform/db"
	"k8s-platform/middle"
	"k8s-platform/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	// 初始化db
	db.Init()
	//关闭db连接池
	defer db.Close()

	// 初始化k8s的clientset
	service.K8s.Init()

	// 初始化路由
	r := gin.Default()

	// 加载中间件

	r.Use(middle.Cors())

	// 登录相关接口
	controller.Router.LoginRouter(r)
	// 前端
	controller.Router.WebRouter(r)

	r.Use(middle.JWTAuth())

	// 加载业务路由
	controller.Router.InitApiRouter(r)

	// 监听websocket
	go func() {
		http.HandleFunc("/ws", service.Terminal.WsHander)
		http.ListenAndServe(config.WsAddr, nil)
	}()

	r.Run(config.ListenAddr)

}
