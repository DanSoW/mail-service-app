FROM golang

RUN go version
ENV GOPATH=/

COPY ./ ./

# Обновление пакетов
RUN apt-get update

# Установка postgresql-client
RUN apt-get -y install postgresql-client

# Установка всех зависимостей
RUN go mod download

# Сборка go-приложения
RUN go build -o smtp ./cmd/main.go

# Запуск приложения
CMD ["./smtp"]