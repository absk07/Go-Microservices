package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type loginRequest struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (app *Config) Login(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": true,
			"message": err,
		})
		return
	}

	user, err := app.Models.User.GetByEmail(req.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": true,
			"message": "Invalid credentials",
		})
		return
	}

	isValid, err := user.PasswordMatches(req.Password)
	if err != nil || !isValid {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": true,
			"message": "Invalid credentials",
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"error": false,
		"message": "User logged in successfully",
		"data": user,
	})
}