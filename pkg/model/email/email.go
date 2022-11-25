package email

/**
 * Структура, описывающая содержимое сообщения
 */
type MailModel struct {
	Sender  string   // Отправитель
	To      []string // Получатели
	Subject string   // Тема сообщения
	Body    string   // Тело сообщения
}
