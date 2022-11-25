# Запуск сервера
run:
	go run cmd/main.go

# Генерация документации
generate-doc:
	swag init -g cmd/main.go

# Генерация docker-образа
docker-build: 
	docker build -t mailer-api-image:latest .

# Запуск контейнера на основе сгенерированного образа
docker-run:
	docker run -d -p 5000:5000 --env-file .env --rm --name mailer-api-app mailer-api-image:latest

# Запуск контейнера в режиме разработки
docker-run-dev:
	docker run -d -p 5000:5000 -v "C:\Projects\DevelopmentProjects\rental-housing\server-app-main:/server-app-main" --env-file .env --rm --name mailer-api-app mailer-api-image:latest 

# Остановка контейнера
docker-stop:
	docker stop mailer-api-app

# Запуск контейнера
docker-start:
	docker start mailer-api-app