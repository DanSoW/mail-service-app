package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

/**
 * Структура errorResponse, описывающая сообщение об ошибке
 */
type errorResponse struct {
	Message string `json:"message"`
}

/**
 * Структура successResponse, описывающая сообщение об удачной обработке запроса
 */
type successResponse struct {
	Message string `json:"message"`
}

/**
 * Структура statusResponse, описывающая статус запроса
 */
type statusResponse struct {
	Status string `json:"status"`
}

/**
 * Структура BooleanResponse, описывающая значение bool запроса
 */
type BooleanResponse struct {
	Value bool `jsong:"value"`
}

/**
 * Метод создания ошибочного ответа на запрос
 * @param {*gin.Context} c Контекст, в рамках которого будет происходить генерации сообщения об ошибке
 * @param {int} statusCode Статус код ответа
 * @param {string} message Сообщение об ошибке
 */
func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}

/**
 * Метод создания ошибочного ответа на запрос
 * @param {*gin.Context} c Контекст, в рамках которого будет происходить генерации сообщения об ошибке
 * @param {int} statusCode Статус код ответа
 * @param {string} message Сообщение об ошибке
 */
func NewErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
