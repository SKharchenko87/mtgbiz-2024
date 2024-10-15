FROM golang:alpine


# Сборка приложения Go
WORKDIR /app
COPY . /app/

# Установка зависимостей Go
RUN go mod download

RUN go build -o main ./main.go

# Открыавем порт 8080
EXPOSE 8090

# Название образа
LABEL name="mtgbiz-server"

CMD ["./main"]