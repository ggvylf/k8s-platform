package controller

import (
	"k8s-platform/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
)

type workflow struct{}

var Workflow workflow

func (w *workflow) GetList(ctx *gin.Context) {
	// 请求参数
	params := new(struct {
		Name      string `form:"name"`
		Namespace string `form:"namespace"`
		Page      int    `form:"page"`
		Limit     int    `form:"limit"`
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

	data, err := service.Workflow.GetList(params.Name, params.Namespace, params.Page, params.Limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return

	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取workflow成功",
		"data": data,
	})

}

func (w *workflow) GetById(ctx *gin.Context) {
	// 请求参数
	params := new(struct {
		ID int `form:"id"`
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

	data, err := service.Workflow.GetById(params.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return

	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取workflow成功",
		"data": data,
	})

}

func (w *workflow) Create(ctx *gin.Context) {

	// 请求参数
	params := &service.WorkflowCreate{}

	// 参数绑定
	if err := ctx.Bind(params); err != nil {
		logger.Error("参数绑定失败：" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "参数绑定失败：" + err.Error(),
			"data": nil,
		})
		return
	}

	err := service.Workflow.CreateWorkFlow(params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return

	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "创建workflow成功",
		"data": nil,
	})
}

func (w *workflow) Delete(ctx *gin.Context) {
	// 请求参数
	params := new(struct {
		ID int `form:"id"`
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

	err := service.Workflow.DeleteWorkflowById(params.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return

	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "删除workflow成功",
		"data": nil,
	})
}
