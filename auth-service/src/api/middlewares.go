package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func jsonLoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		logs := make(map[string]any)

		logs["status_code"] = ctx.Writer.Status()
		logs["path"] = ctx.Request.URL.Path
		logs["method"] = ctx.Request.Method
		logs["start_time"] = ctx.Request.Header.Get("Date")
		logs["remote_addr"] = ctx.ClientIP()
		logs["response_time"] = ctx.Writer.Header().Get("X-Response-Time")

		s, _ := json.Marshal(logs)

		log.Println("****LOGS****", string(s)+"\n")

		err := logRequest("name", string(s))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"error":   true,
				"message": "unable to connect to logger service",
			})
			return
		}
	}
}

func logRequest(name, data string) error {
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