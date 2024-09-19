package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/auth-service/data"
	"github.com/gin-gonic/gin"
)

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type registerRequest struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Password  string `json:"password"`
}

func (app *Config) Register(ctx *gin.Context) {
	var req registerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err,
		})
		return
	}

	existingUser, _ := app.Models.User.GetByEmail(req.Email)

	if existingUser != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{
			"error":   true,
			"message": "User already exists!",
		})
		return
	}

	user := data.User{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  req.Password,
		Active:    1,
	}

	_, err := app.Models.User.Insert(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": err,
		})
		return
	}

	err = sendMail(MailPayload{
		From:    "admin@gmail.com",
		To:      user.Email,
		Subject: "Registration Successfull!",
		Message: "Welcome " + user.FirstName + " " + user.LastName,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": err,
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"error":   false,
		"message": "User registered successfully",
	})
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Login(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": err,
		})
		return
	}

	user, err := app.Models.User.GetByEmail(req.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Invalid credentials",
		})
		return
	}

	isValid, err := user.PasswordMatches(req.Password)
	if err != nil || !isValid {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error":   true,
			"message": "Invalid credentials",
		})
		return
	}

	ctx.JSON(http.StatusAccepted, gin.H{
		"error":   false,
		"message": "User logged in successfully",
		"data":    user,
	})
}

func sendMail(mailPayload MailPayload) error {
	jsonData, _ := json.Marshal(mailPayload)

	request, err := http.NewRequest("POST", "http://mail-service:6060/send", bytes.NewBuffer(jsonData))

	if err != nil {
		// log.Println("http req err", err)
		return err
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		// log.Println("http res err", err)
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return err
	}

	return nil
}