package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type router struct{}

var Router router

func (r *router) InitApiRouter(router *gin.Engine) {
	router.GET("/testapi", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":  "testapi success!",
			"data": nil,
		})
	})
}
