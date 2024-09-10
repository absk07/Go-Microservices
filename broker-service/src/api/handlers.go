package main

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type jsonResponse struct{
	Error bool `json:"error"`
	Message string `json:"message"`
	Data any `json:"data,omitempty"`
}

func (app *Config) Broker(ctx *gin.Context) {
	payload := jsonResponse{
		Error: false,
		Message: "Hit the broker",
	}

	res, err := json.Marshal(payload)
	
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	ctx.JSON(http.StatusOK, res)
}