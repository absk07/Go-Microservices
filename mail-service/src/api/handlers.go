package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type mailMessage struct {
	From string `json:"from"`
	To string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (app *Config) SendMail(ctx *gin.Context) {
	var req mailMessage
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"message": err,
		})
		return
	}

	msg := Message {
		From: req.From,
		To: req.To,
		Subject: req.Subject,
		Data: req.Message,
	}

	err := app.Mailer.SendSMTPMessage(msg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": true,
			"message": err,
		})
		return
	}

	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": false,
		"message": "email sent to " + req.To,
	})
}