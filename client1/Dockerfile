FROM golang:alpine


# Сборка приложения Go
WORKDIR /app
COPY . /app/

# Установка зависимостей Go
RUN go mod download

RUN go build -o main ./main.go

# Открыавем порт 8080
EXPOSE 8091

# Название образа
LABEL name="mtgbiz-client1"

CMD ["./main"]