package service

import (
	messageModel "smtp/pkg/model/message"
	"smtp/pkg/repository"
)

type Mailer interface {
	SendMail(message *messageModel.MessageInputModel) (bool, error)
	GetMail(output *messageModel.OutputModel) (*messageModel.MessagesModel, error)
}

type Service struct {
	Mailer
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Mailer: NewMailerService(repos.Mailer),
	}
}
