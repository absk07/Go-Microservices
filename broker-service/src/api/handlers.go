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
	Action   string          `json:"action"`
	Register RegisterPayload `json:"register,omitempty"`
	Login    LoginPayload    `json:"login,omitempty"`
	// Log      LogPayload      `json:"log,omitempty"`
	// Mail     MailPayload     `json:"mail,omitempty"`
}

type RegisterPayload struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Password  string `json:"password"`
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// type MailPayload struct {
// 	From    string `json:"from"`
// 	To      string `json:"to"`
// 	Subject string `json:"subject"`
// 	Message string `json:"message"`
// }

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
	case "register":
		app.register(ctx, req.Register)
	case "login":
		app.login(ctx, req.Login)
	// case "log":
	// 	app.logEventViaRabbitMQ(ctx, req.Log)
	// case "mail":
	// 	app.sendMail(ctx, req.Mail)
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "unknown action",
		})
	}
}

func (app *Config) register(ctx *gin.Context, registerPayload RegisterPayload) {
	jsonData, _ := json.Marshal(registerPayload)

	request, err := http.NewRequest("POST", "http://auth-service:9090/register", bytes.NewBuffer(jsonData))

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
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "invalid credentials!",
		})
		return
	} else if response.StatusCode == http.StatusNotAcceptable {
		ctx.JSON(http.StatusNotAcceptable, gin.H{
			"error":   true,
			"message": "User already exists!",
		})
		return
	} else if response.StatusCode != http.StatusAccepted || response.StatusCode == http.StatusInternalServerError {
		ctx.JSON(http.StatusInternalServerError, gin.H{
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": data.Message,
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"error":   data.Error,
		"message": data.Message,
	})
}

func (app *Config) login(ctx *gin.Context, loginPayload LoginPayload) {
	jsonData, _ := json.Marshal(loginPayload)

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
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "invalid credentials!",
		})
		return
	} else if response.StatusCode != http.StatusAccepted || response.StatusCode == http.StatusInternalServerError {
		ctx.JSON(http.StatusInternalServerError, gin.H{
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": data.Message,
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"error":   data.Error,
		"message": data.Message,
		"data":    data.Data,
	})
}

/** func (app *Config) logger(ctx *gin.Context, logPayload LogPayload) {
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
		ctx.JSON(http.StatusInternalServerError, gin.H{
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

	ctx.JSON(http.StatusAccepted, gin.H{
		"error":   false,
		"message": "Logged!",
	})
} */

/** func (app *Config) sendMail(ctx *gin.Context, mailPayload MailPayload) {
	jsonData, _ := json.Marshal(mailPayload)

	request, err := http.NewRequest("POST", "http://mail-service:6060/send", bytes.NewBuffer(jsonData))

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

	if response.StatusCode != http.StatusOK {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "error calling mail service!",
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"error":   false,
		"message": "Email sent",
	})
} */