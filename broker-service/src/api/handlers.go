package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *Config) Broker(ctx *gin.Context) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	ctx.JSON(http.StatusOK, payload)
}
