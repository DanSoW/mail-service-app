package smtp

import (
	"context"
	"net/http"
	"time"
)

/**
 * Структура сервера
 */
type Server struct {
	httpServer *http.Server
}

/**
 * Метод для запуска сервера
 * @param {string} port Порт для запуска сервера
 * @param {http.Handler} handler Обработчики событий
 * @returns {error} Сообщение об ошибке (или nil)
 */
func (s *Server) Run(port string, handler http.Handler) error {

	// Создание экземпляра объекта сервера
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 10 << 20, // 10 MB
		ReadTimeout:    10 * time.Minute,
		WriteTimeout:   10 * time.Minute,
	}

	// Возвращение результата запуска процесса прослушивания сообщений сервера
	return s.httpServer.ListenAndServe()
}

/**
 * Метод для завершения работы сервера
 * @param {context.Context} ctx Контекст завершения обработки сообщений сервером
 * @returns {error} Сообщение об ошибке (или nil)
 */
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
