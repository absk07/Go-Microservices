package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/logger-service/data"
)

type logRequest struct {
	Method       string `json:"method"`
	Path         string `json:"path"`
	RemoteAddr   string `json:"remote_addr"`
	ResponseTime string `json:"response_time"`
	StartTime    string `json:"start_time"`
	StatusCode   string `json:"status_code"`
}

func (app *Config) WriteLog(ctx *gin.Context) {
	var req logRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err,
		})
		return
	}

	event := data.LogEntry{
		Method:       req.Method,
		Path:         req.Path,
		RemoteAddr:   req.RemoteAddr,
		ResponseTime: req.ResponseTime,
		StartTime:    req.StartTime,
		StatusCode:   req.StatusCode,
	}

	err := app.Models.LogEntry.Insert(event)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": err,
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"error":   false,
		"message": "logs saved",
	})
}
