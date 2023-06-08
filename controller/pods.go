package controller

import (
	"k8s-platform/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
)

type pod struct{}

var Pod pod

// 获取podlist
func (p *pod) GetPods(ctx *gin.Context) {

	// 请求参数
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
	})

	if err := ctx.Bind(params); err != nil {
		logger.Error("参数绑定失败：" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "参数绑定失败：" + err.Error(),
			"data": nil,
		})
		return
	}

	data, err := service.Pod.GetPods(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return

	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取pod列表成功",
		"data": data,
	})
}
