package main

import (
	"bytes"
	"encoding/json"

	"log"
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

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LogPayload struct {
	Name string `json:"name" binding:"required"`
	Data string `json:"data" binding:"required"`
}

func (app *Config) HandleSubmission(ctx *gin.Context) {
	log.Println("before req body")
	var req RequestPayload
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err,
		})
		return
	}

	log.Println("req body", req)

	switch req.Action {
	case "auth":
		app.authenticate(ctx, req.Auth)
	case "log":
		app.logger(ctx, req.Log)
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "unknown action",
		})
	}
}

func (app *Config) authenticate(ctx *gin.Context, authPayload AuthPayload) {
	jsonData, _ := json.Marshal(authPayload)

	request, err := http.NewRequest("POST", "http://auth-service:9090/login", bytes.NewBuffer(jsonData))

	if err != nil {
		// log.Println("http req err", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err,
		})
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		// log.Println("http res err", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err,
		})
		return
	}
	defer response.Body.Close()

	// log.Println("response", response)
	// log.Println("res code", response.StatusCode)
	// log.Println("http code", http.StatusAccepted)

	if response.StatusCode == http.StatusUnauthorized {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "invalid credentials!",
		})
		return
	} else if response.StatusCode != http.StatusAccepted {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "error calling auth service!",
		})
		return
	}

	var data jsonResponse
	err = json.NewDecoder(response.Body).Decode(&data)

	if err != nil {
		// log.Println("json decoder err", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err,
		})
		return
	}

	if data.Error {
		// log.Println("data err", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err,
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"error":   false,
		"message": "Authenticated!",
		"data":    data.Data,
	})
}

func (app *Config) logger(ctx *gin.Context, logPayload LogPayload) {
	jsonData, _ := json.Marshal(logPayload)

	request, err := http.NewRequest("POST", "http://logger-service:7070/log", bytes.NewBuffer(jsonData))

	if err != nil {
		// log.Println("http req err", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err,
		})
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		// log.Println("http res err", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err,
		})
		return
	}
	defer response.Body.Close()

	// log.Println("response", response)
	// log.Println("res code", response.StatusCode)
	// log.Println("http code", http.StatusAccepted)

	if response.StatusCode != http.StatusAccepted {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "error calling log service!",
		})
		return
	}

	var data jsonResponse
	err = json.NewDecoder(response.Body).Decode(&data)

	if err != nil {
		// log.Println("json decoder err", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err,
		})
		return
	}

	if data.Error {
		// log.Println("data err", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err,
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"error":   false,
		"message": "Logged!",
		"data":    data.Data,
	})
}
