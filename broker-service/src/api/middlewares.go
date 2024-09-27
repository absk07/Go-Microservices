package main

import (
	"encoding/json"
	"time"

	// "log"
	"net/http"
	"strconv"

	"github.com/broker-service/events"
	"github.com/gin-gonic/gin"
)

const (
	SeverityInfo    = "INFO"
	SeverityWarning = "WARNING"
	SeverityError   = "ERROR"
)

type Log struct {
	Method       string `json:"method"`
	Path         string `json:"path"`
	RemoteAddr   string `json:"remote_addr"`
	ResponseTime string `json:"response_time"`
	StartTime    string `json:"start_time"`
	StatusCode   string `json:"status_code"`
	Severity     string `json:"severity"`
}

func jsonLoggerMiddleware(app *Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()
		
		ctx.Next()

		severity := getSeverity(ctx.Writer.Status())

		endTime := time.Now()
		duration := endTime.Sub(startTime)

		ctx.Writer.Header().Set("X-Response-Time", strconv.FormatFloat(duration.Seconds()*1000, 'f', 2, 64)+"ms")

		payload := Log{
			StatusCode: strconv.Itoa(ctx.Writer.Status()),
			Path: ctx.Request.URL.Path,
			Method: ctx.Request.Method,
			StartTime: startTime.Format(time.RFC3339),
			RemoteAddr: ctx.ClientIP(),
			ResponseTime: ctx.Writer.Header().Get("X-Response-Time"),
			Severity: severity,
		}


		err := app.logEventViaRabbitMQ(payload)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": err,
			})
			return
		}
	}
}

func getSeverity(statusCode int) string {
	switch {
	case statusCode >= 500:
		return SeverityError
	case statusCode >= 400:
		return SeverityWarning
	default:
		return SeverityInfo
	}
}

func (app *Config) logEventViaRabbitMQ(logPayload Log) error {
	err := app.PushToQueue(logPayload)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) PushToQueue(logPayload Log) error {
	emitter, err := events.NewEventEmitter(app.RabbitConn)
	if err != nil {
		return err
	}

	payload := Log{
		Method: logPayload.Method,
		Path: logPayload.Path,
		RemoteAddr: logPayload.RemoteAddr,
		ResponseTime: logPayload.ResponseTime,
		StartTime: logPayload.StartTime,
		StatusCode: logPayload.StatusCode,
		Severity: logPayload.Severity,
	}

	data, _ := json.Marshal(&payload)

	if payload.Severity == SeverityError {
		if err := emitter.Push(string(data), "log.ERROR"); err != nil {
			return err
		}
	} else if payload.Severity == SeverityWarning {
		if err := emitter.Push(string(data), "log.WARNING"); err != nil {
			return err
		}
	} else {
		if err := emitter.Push(string(data), "log.INFO"); err != nil {
			return err
		}
	}

	return nil
}
