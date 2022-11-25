package handler

import (
	"smtp/pkg/service"

	_ "smtp/docs"

	"github.com/gin-gonic/gin"
	_ "github.com/swaggo/files"
	swaggerFiles "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"
	ginSwagger "github.com/swaggo/gin-swagger"
)

/**
 * Структура Handler, описывающая слой обработки запросов
 */
type Handler struct {
	services *service.Service
}

/**
 * Функция для создания нового объекта структуры Handler
 * @params {*service.Service} Указатель на объект структуры Service
 * @retuns {*Handler} Указатель на объект структуры Handler
 */
func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

/**
 * Метод, для инициализации маршрутов сервера
 * @returns {*gin.Engine} Экземпляр объекта gin.Engige
 */
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	// Установка максимального размера multipart-form/data
	router.MaxMultipartMemory = 50 << 20 // 50 MiB

	// Установка папки для получения доступа к статическим данным
	router.Static("/public", "./public")

	// Инициализация маршрута для просмотра документации REST API
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// URL: /mailer
	mailer := router.Group("/mailer")
	{
		// URL: /mailer/send
		mailer.POST("/send", h.sendMail)

		// URL: /mailer/get
		mailer.POST("/get", h.getMail)
	}

	return router
}
