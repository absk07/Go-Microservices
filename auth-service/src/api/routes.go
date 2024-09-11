package main

import (
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

func (app *Config) Routes() *gin.Engine {
	router := gin.Default()

	router.Use(cors.Default())

	return router
}
