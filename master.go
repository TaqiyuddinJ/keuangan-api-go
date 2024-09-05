package main

import (
	"github.com/gin-gonic/gin"
)

func MasterRoute(router *gin.Engine) {
	group := router.Group("/keuangan/master")
	{
		group.GET("/akun", func(context *gin.Context) {
			context.JSON(200, gin.H{
				"success": true,
			})
		})
	}
}
