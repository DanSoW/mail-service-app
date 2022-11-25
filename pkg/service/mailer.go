package service

import (
	messageModel "smtp/pkg/model/message"
	repository "smtp/pkg/repository"
)

type MailerService struct {
	repo repository.Mailer
}

func NewMailerService(repo repository.Mailer) *MailerService {
	return &MailerService{
		repo: repo,
	}
}

func (s *MailerService) SendMail(message *messageModel.MessageInputModel) (bool, error) {
	return s.repo.SendMail(message)
}

func (s *MailerService) GetMail(output *messageModel.OutputModel) (*messageModel.MessagesModel, error) {
	return s.repo.GetMail(output)
}
