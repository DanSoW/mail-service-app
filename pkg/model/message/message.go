package message

import (
	fileModel "smtp/pkg/model/file"
)

/**
 * Структура, описывающая полное содержимое сообщения пользователя
 */
type MessageInputModel struct {
	Receivers []string               `json:"receiver" binding:"required"` // Получатели сообщения
	Subject   string                 `json:"subject" binding:"required"`  // Тема сообщения
	Message   string                 `json:"message" binding:"required"`  // Тело сообщения
	Files     []*fileModel.FileModel `json:"files"`                       // Прикреплённые файлы
}

/**
 * Структура, описывающая полное содержимое сообщения пользователя
 */
type MessageOutputModel struct {
	Sender  string                 `json:"sender" binding:"required"`  // Отправитель сообщения
	Subject string                 `json:"subject" binding:"required"` // Тема сообщения
	Message string                 `json:"message" binding:"required"` // Тело сообщения
	Files   []*fileModel.FileModel `json:"files"`                      // Прикреплённые файлы
}

/**
 * Структура, описывающая множество сообщений
 */
type MessagesModel struct {
	Messages []MessageOutputModel // Множество сообщений
}

/**
 * Структура, описывающая входные данные для чтения множества сообщений
 */
type OutputModel struct {
	Count uint32 `json:"count" binding:"required"` // Количество считываемых сообщений
}
