package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type loginRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
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

	err = app.logRequest("login", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"error": true,
			"message": "unable to connect to logger service",
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"error": false,
		"message": "User logged in successfully",
		"data": user,
	})
}

func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.Marshal(entry)

	request, err := http.NewRequest("POST", "http://logger-service:7070/log", bytes.NewBuffer(jsonData))

	if err != nil {
		// log.Println("http req err", err)
		return err
	}

	client := &http.Client{}
	_, err = client.Do(request)

	if err != nil {
		// log.Println("http res err", err)
		return err
	}

	return nil
}