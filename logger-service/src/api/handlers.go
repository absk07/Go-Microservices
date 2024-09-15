package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/logger-service/data"
)

type logRequest struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(ctx *gin.Context) {
	var req logRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"message": err,
		})
		return
	}

	event := data.LogEntry{
		Name: req.Name,
		Data: req.Data,
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"message": err,
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"error": false,
		"message": "logs saved",
	})
}