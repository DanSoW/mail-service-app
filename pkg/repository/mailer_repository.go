package repository

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"smtp/pkg/model/email"
	fileModel "smtp/pkg/model/file"
	messageModel "smtp/pkg/model/message"
	util "smtp/pkg/util"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message"
	"github.com/emersion/go-message/mail"
	"github.com/spf13/viper"
)

type MailerRepository struct{}

func NewMailerRepository() *MailerRepository {
	return &MailerRepository{}
}

/**
 * Метод репозитория для отправки сообщения
 */
func (r *MailerRepository) SendMail(message *messageModel.MessageInputModel) (bool, error) {

	msg := util.BuildMessage(&email.MailModel{
		To:      message.Receivers,
		Sender:  viper.GetString("smtp.email"),
		Subject: message.Subject,
		Body:    message.Message,
	},
		message.Files,
	)

	err := util.SendMessage(msg)

	if err != nil {
		return false, err
	}

	// Удаление файлов (после отправки файлы не сохраняются, а удаляются сразу)
	for _, element := range message.Files {
		os.Remove(element.Filepath)
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

/**
 * Метод репозитория для получения сообщений
 */
func (r *MailerRepository) GetMail(output *messageModel.OutputModel) (*messageModel.MessagesModel, error) {
	// Подключение к серверу
	connect, err := client.DialTLS(viper.GetString("imap.host")+":993", nil)
	if err != nil {
		return nil, err
	}

	// Завершение подключения при освобождении ресурсов
	defer connect.Logout()

	// Авторизация пользователя
	if err := connect.Login(viper.GetString("imap.email"), os.Getenv("SMTP_PASSWORD")); err != nil {
		return nil, err
	}

	// Список mailbox (категорий писем)
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- connect.List("", "*", mailboxes)
	}()

	log.Println("Mailboxes:")
	for m := range mailboxes {
		log.Println("* " + m.Name)
	}

	if err := <-done; err != nil {
		return nil, err
	}

	// Получение сообщений
	mbox, err := connect.Select(viper.GetString("imap.inbox"), false)
	if err != nil {
		return nil, err
	}

	seqset := new(imap.SeqSet)

	if output.Count > mbox.Messages {
		output.Count = mbox.Messages
	}

	seqset.AddRange(1, output.Count)

	var section imap.BodySectionName
	items := []imap.FetchItem{section.FetchItem(), imap.FetchEnvelope, imap.FetchRFC822}

	// Сообщения
	var messages messageModel.MessagesModel

	msgs := make(chan *imap.Message, 10)
	done = make(chan error, 1)
	go func() {
		done <- connect.Fetch(seqset, items, msgs)
	}()

	// Цикл прохода по всем сообщениям
	for msg := range msgs {
		// Формирование одного сообщения
		var oneMessage messageModel.MessageOutputModel
		oneMessage.Sender = msg.Envelope.From[len(msg.Envelope.From)-1].Address()
		oneMessage.Subject = msg.Envelope.Subject

		r := msg.GetBody(&section)
		if r == nil {
			log.Fatal("Сервер не вернул тело сообщения")
		}

		// Создание нового потока чтения
		messageReader, err := mail.CreateReader(r)
		if err != nil {
			return nil, err
		}

		// Формирование текста сообщения
		var text string = ""

		for {
			p, err := messageReader.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}

			switch p.Header.(type) {
			case *mail.InlineHeader:
				b, _ := ioutil.ReadAll(p.Body)
				text += string(b)
			}
		}

		oneMessage.Message = text

		// Загрузка файлов из сообщения (если таковые имеются)
		for _, r := range msg.Body {
			entity, err := message.Read(r)
			if err != nil {
				return nil, err
			}

			multiPartReader := entity.MultipartReader()

			// Если файлы отсутствуют - пропускаем процесс их загрузки на сервер
			if multiPartReader == nil {
				continue
			}

			for e, err := multiPartReader.NextPart(); err != io.EOF && multiPartReader != nil; e, err = multiPartReader.NextPart() {
				kind, params, cErr := e.Header.ContentType()
				if cErr != nil {
					return nil, cErr
				}

				if kind != "image/png" && kind != "image/gif" && kind != "image/jpeg" {
					continue
				}

				data, rErr := ioutil.ReadAll(e.Body)
				if rErr != nil {
					return nil, cErr
				}

				if fErr := ioutil.WriteFile("./public/"+params["name"], data, 0777); fErr != nil {
					return nil, fErr
				}

				// Добавление ссылки на загруженный файл
				oneMessage.Files = append(oneMessage.Files, &fileModel.FileModel{
					Filename: params["name"],
					Filepath: viper.GetString("api_url") + "/public/" + params["name"],
				})
			}

		}

		// Добавление информации о новом сообщении
		messages.Messages = append(messages.Messages, oneMessage)
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	return &messages, err
}
