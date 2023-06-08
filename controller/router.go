package controller

import (
	"github.com/gin-gonic/gin"
)

type router struct{}

var Router router

func (r *router) InitApiRouter(router *gin.Engine) {
	router.GET("/api/k8s/pods", Pod.GetPods)
}
