package repository

import messageModel "smtp/pkg/model/message"

type Mailer interface {
	SendMail(message *messageModel.MessageInputModel) (bool, error)
	GetMail(output *messageModel.OutputModel) (*messageModel.MessagesModel, error)
}

type Repository struct {
	Mailer
}

func NewRepository() *Repository {
	return &Repository{
		Mailer: NewMailerRepository(),
	}
}
