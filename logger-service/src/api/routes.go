package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

func (app *Config) Routes() *gin.Engine {
	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "pong!",
		})
	})

	router.POST("/log", app.WriteLog)

	return router
}