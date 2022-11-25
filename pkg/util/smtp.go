package util

import (
	"os"
	"smtp/pkg/model/email"
	"strconv"

	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"

	fileModel "smtp/pkg/model/file"
)

/**
 * Функция создания строки сообщения по почтовым данным
 * @param {email.Mail} mail Почтовые данные
 * @returns {string} Результирующая строка, содержащая все необходимые данные
 */
func BuildMessage(mail *email.MailModel, files []*fileModel.FileModel) *gomail.Message {
	message := gomail.NewMessage()
	message.SetHeader("From", mail.Sender)
	message.SetHeader("To", mail.To...)
	// message.SetAddressHeader("Cc", mail.Sender, "Daniil")
	message.SetHeader("Subject", mail.Subject)
	message.SetBody("text/html", mail.Body)

	for _, element := range files {
		message.Attach(element.Filepath)
	}

	return message
}

/**
 * Функция отправки сообщения получателю
 * @param {*gomail.Message} message Сообщение для отправки
 * @returns {error} Сообщение об ошибке или значение nil
 */
func SendMessage(message *gomail.Message) error {
	port, _ := strconv.Atoi(viper.GetString("smtp.port"))
	dialer := gomail.NewDialer(viper.GetString("smtp.host"), port, viper.GetString("smtp.email"), os.Getenv("SMTP_PASSWORD"))
	// dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err := dialer.DialAndSend(message)

	return err
}
