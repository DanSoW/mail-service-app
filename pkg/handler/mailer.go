package handler

import (
	"net/http"
	fileModel "smtp/pkg/model/file"
	messageModel "smtp/pkg/model/message"
	"strings"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// @Summary SendMail
// @Tags mailer
// @Description Отправка сообщения пользователю
// @ID mailer-send-mail
// @Accept  json
// @Produce  json
// @Param files formData []string true "Файлы"
// @Param receiver formData string true "Адрес получателя сообщения"
// @Param message formData string true "Сообщение для отправки"
// @Success 200 {object} bool "data"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /mailer/send [post]
func (h *Handler) sendMail(c *gin.Context) {
	// Получение данных в формате multipart-form/data
	form, err := c.MultipartForm()

	// Обработка ошибок
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Получение информации из multipart-form/data
	files := form.File["file"]
	receivers := c.PostFormArray("receiver")
	message := c.PostForm("message")
	subject := c.PostForm("subject")

	// Массив для хранения информации о файлах локально
	var fileInfo []*fileModel.FileModel

	for _, element := range files {
		// Генерация нового имени файла и его пути
		filename := uuid.NewV4().String()

		headers := strings.Split(element.Header["Content-Disposition"][0], ";")
		fileFormat := strings.Split(headers[2], ".")[1]
		fileFormat = strings.Trim(fileFormat, "\"")

		filepath := "public/" + filename + "." + fileFormat

		fileInfo = append(fileInfo, &fileModel.FileModel{
			Filename: filename,
			Filepath: filepath,
		})

		// Сохранение файла локально
		c.SaveUploadedFile(element, filepath)
	}

	// Отправка сообщения
	data, err := h.services.SendMail(&messageModel.MessageInputModel{
		Receivers: receivers,
		Message:   message,
		Subject:   subject,
		Files:     fileInfo,
	})

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, data)
}

// @Summary GetMail
// @Tags mailer
// @Description Получение сообщений из почты пользователя
// @ID mailer-get-mail
// @Accept  json
// @Produce  json
// @Param input body messageModel.OutputModel true "Настройки"
// @Success 200 {object} messageModel.MessagesModel "messages"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /mailer/get [post]
func (h *Handler) getMail(c *gin.Context) {
	var input messageModel.OutputModel

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	data, err := h.services.GetMail(&input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, data)
}
