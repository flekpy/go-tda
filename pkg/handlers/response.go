package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Message string `json:"message"`
}

type statusResponse struct {
	Status string `json:"status"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	// блокирует выполнение последующих обработчиков + записывает ответ в формате json
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
