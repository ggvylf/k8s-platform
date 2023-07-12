package controller

import (
	"github.com/gin-gonic/gin"
)

type router struct{}

var Router router

func (r *router) InitApiRouter(router *gin.Engine) {

	// pod
	router.GET("/api/k8s/pods", Pod.GetPods)
	router.GET("/api/k8s/poddetail", Pod.GetPodDetail)
	router.GET("/api/k8s/podsnp", Pod.GetPodSumPerNp)
	router.GET("/api/k8s/containers", Pod.GetContainerNameList)
	router.GET("/api/k8s/podlog", Pod.GetPodLog)
	router.DELETE("/api/k8s/deletepod", Pod.DeletePod)
	router.PUT("/api/k8s/updatepod", Pod.UpdatePod)

	// deployment
	router.GET("/api/k8s/deployments", Deployment.GetDeployments)
	router.GET("/api/k8s/deploymentdetail", Deployment.GetDeploymentDetail)
	router.GET("/api/k8s/scaledeployment", Deployment.ScaleDeployment)
	router.GET("/api/k8s/createdeployment", Deployment.CreateDeployment)
	router.GET("/api/k8s/restartdeployment", Deployment.RestartDeployment)
	router.GET("/api/k8s/deploymentsnp", Deployment.GetDeploymentSumPerNp)

	router.PUT("/api/k8s/updatedeployment", Deployment.UpdateDeployment)
	router.DELETE("/api/k8s/deletedeployment", Deployment.DeleteDeployment)

	// workflow
	router.GET("/api/workflow/list", Workflow.GetList)
	router.GET("/api/workflow/getdetail", Workflow.GetById)
	router.GET("/api/workflow/create", Workflow.Create)
	router.DELETE("/api/workflow/delete", Workflow.Delete)
}
