package utils

import "github.com/gin-gonic/gin"

type Responses struct {
	StatusCode int         `json:"statusCode"`
	Method     string      `json:"method"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

func APIResponse(ctx *gin.Context, message string, statusCode int, method string, data interface{}) {
	jsonResponse := Responses{
		StatusCode: statusCode,
		Method:     method,
		Message:    message,
		Data:       data,
	}

	if statusCode >= 400 {
		ctx.JSON(statusCode, jsonResponse)
		defer ctx.AbortWithStatus(statusCode)
	} else {
		ctx.JSON(statusCode, jsonResponse)
	}
}
