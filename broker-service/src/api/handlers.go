package main

import (
	"bytes"
	"encoding/json"
	// "log"
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
}

type AuthPayload struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (app *Config) HandleSubmission(ctx *gin.Context) {
	var req RequestPayload
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err,
		})
		return
	}

	switch req.Action {
	case "auth":
		app.authenticate(ctx, req.Auth)
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "unknown action",
		})
	}
}

func (app *Config) authenticate(ctx *gin.Context, authPayload AuthPayload) {
	jsonData, _ := json.MarshalIndent(authPayload, "", "\t")

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
