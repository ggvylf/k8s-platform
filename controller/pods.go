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

	// 参数绑定
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

// 获取pod详情
func (p *pod) GetPodDetail(ctx *gin.Context) {

	// 请求参数
	params := new(struct {
		PodName   string `form:"pod_name"`
		Namespace string `form:"namespace"`
	})

	// 参数绑定
	if err := ctx.Bind(params); err != nil {
		logger.Error("参数绑定失败：" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "参数绑定失败：" + err.Error(),
			"data": nil,
		})
		return
	}

	data, err := service.Pod.GetPodDetail(params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return

	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取pod详情成功",
		"data": data,
	})
}

// 删除pod
func (p *pod) DeletePod(ctx *gin.Context) {

	// 请求参数
	params := new(struct {
		PodName   string `json:"pod_name"`
		Namespace string `json:"namespace"`
	})

	// 参数绑定
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("参数绑定失败：" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "参数绑定失败：" + err.Error(),
			"data": nil,
		})
		return
	}

	err := service.Pod.DeletePod(params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return

	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "删除pod成功",
		"data": nil,
	})
}

// 更新pod
// 获取podlist
func (p *pod) UpdatePod(ctx *gin.Context) {

	// 请求参数
	params := new(struct {
		Name       string `form:"name"`
		Namespace  string `form:"namespace"`
		UpdateInfo string `form: "updateinfo"`
	})

	// 参数绑定
	if err := ctx.Bind(params); err != nil {
		logger.Error("参数绑定失败：" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "参数绑定失败：" + err.Error(),
			"data": nil,
		})
		return
	}

	err := service.Pod.UpdatePod(params.Name, params.Namespace, params.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return

	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "更新pod成功",
		"data": nil,
	})
}
