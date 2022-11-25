package main

import (
	"context"
	"os"
	"os/signal"
	mainserver "smtp"
	handler "smtp/pkg/handler"
	repository "smtp/pkg/repository"
	"smtp/pkg/service"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"

	"github.com/spf13/viper"
)

// @title Почтовый сервис
// @version 1.0
// description Проект для рассылки почты

// @host localhost:5000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	// Инициализация конфигурационного файла
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variable: %s", err.Error())
	}

	// Инициализация логгера
	logrus.SetFormatter(new(logrus.JSONFormatter))

	// Открытие файла для записи в него логов
	fileError, err := os.OpenFile(viper.GetString("paths.logs.error"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		// Прикрепление хука к logrus, для записи логов ошибки в конкретный файл
		logrus.AddHook(&writer.Hook{
			Writer: fileError,
			LogLevels: []logrus.Level{
				logrus.ErrorLevel,
			},
		})
	} else {
		logrus.SetOutput(os.Stderr)
		logrus.Error("Не удалось загрузить файл, используется stderr по умолчанию")
	}

	defer fileError.Close()

	fileInfo, err := os.OpenFile(viper.GetString("paths.logs.info"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logrus.AddHook(&writer.Hook{
			Writer: fileInfo,
			LogLevels: []logrus.Level{
				logrus.InfoLevel,
				logrus.DebugLevel,
			},
		})
	} else {
		logrus.SetOutput(os.Stderr)
		logrus.Error("Не удалось загрузить файл, используется stderr по умолчанию")
	}

	defer fileInfo.Close()

	fileWarn, err := os.OpenFile(viper.GetString("paths.logs.warn"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logrus.AddHook(&writer.Hook{
			Writer: fileWarn,
			LogLevels: []logrus.Level{
				logrus.WarnLevel,
			},
		})
	} else {
		logrus.SetOutput(os.Stderr)
		logrus.Error("Не удалось загрузить файл, используется stderr по умолчанию")
	}

	defer fileWarn.Close()

	fileFatal, err := os.OpenFile(viper.GetString("paths.logs.fatal"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logrus.AddHook(&writer.Hook{
			Writer: fileFatal,
			LogLevels: []logrus.Level{
				logrus.FatalLevel,
			},
		})
	} else {
		logrus.SetOutput(os.Stderr)
		logrus.Error("Не удалось загрузить файл, используется stderr по умолчанию")
	}

	defer fileFatal.Close()

	// Реализация паттерна dependency injection
	repos := repository.NewRepository()
	service := service.NewService(repos)
	handlers := handler.NewHandler(service)

	// Создание нового экземпляра сервера
	srv := new(mainserver.Server)

	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("Произошла ошибка при запуске http-сервера: %s", err.Error())
		}
	}()

	logrus.Print("Сервер почтового сервера стартовал")

	// Реализация Graceful Shutdown
	// Блокировка функции main с помощью канала os.Signal
	quit := make(chan os.Signal, 1)

	// Запись в канал, если процесс, в котором выполняется приложение
	// получит сигнал SIGTERM или SIGINT
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	// Чтение из канала, блокирующая выполнение функции main
	<-quit

	logrus.Print("Сервер почтового сервера остановился")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("Произошла ошибка при завершении работы сервера: %s", err.Error())
	}
}

/**
 * Функция для инициализации объекта viper для получения информации из конфигурационных файлов
 */
func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
