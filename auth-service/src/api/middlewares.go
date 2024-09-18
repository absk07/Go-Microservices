package main

import (
	"bytes"
	"encoding/json"

	// "log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func jsonLoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		status_code := strconv.Itoa(ctx.Writer.Status())
		path := ctx.Request.URL.Path
		method := ctx.Request.Method
		start_time := ctx.Request.Header.Get("Date")
		remote_addr := ctx.ClientIP()
		response_time := ctx.Writer.Header().Get("X-Response-Time")

		err := logRequest(method, path, remote_addr, response_time, start_time, status_code)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"error":   true,
				"message": "unable to connect to logger service",
			})
			return
		}
	}
}

func logRequest(method, path, remote_addr, response_time, start_time, status_code string) error {
	var entry struct {
		Method       string `json:"method"`
		Path         string `json:"path"`
		RemoteAddr   string `json:"remote_addr"`
		ResponseTime string `json:"response_time"`
		StartTime    string `json:"start_time"`
		StatusCode   string `json:"status_code"`
	}

	entry.Method = method
	entry.Path = path
	entry.RemoteAddr = remote_addr
	entry.ResponseTime = response_time
	entry.StartTime = start_time
	entry.StatusCode = status_code

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
