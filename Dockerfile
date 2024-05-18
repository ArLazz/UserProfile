# Указываем базовый образ с Go
FROM golang:latest

COPY ./ ./

# Сборка Go приложения
RUN go mod download
RUN go build -o /userprofile-server cmd/userprofile-server/main.go

# Указываем команду для запуска приложения
CMD ["/userprofile-server", "--host=0.0.0.0", "--port=8080"]

# Указываем порт, который будет использоваться
