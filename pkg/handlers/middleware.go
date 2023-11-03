package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	// валидируем, что не пустой
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "пустой хедер ауф")
		return
	}

	// разделяем строку по пробелу
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "невалидный ауф хедер")
		return
	}

	// парсим токен
	userId, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.Set(userCtx, userId)
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "юзер айди не найден")
		return 0, errors.New("юзер айди не найден")
	}

	idInt, ok := id.(int)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "юзер айди invalid type")
		return 0, errors.New("юзер айди не найден")
	}

	return idInt, nil
}
