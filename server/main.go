package main

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

type message struct {
	id  uuid.UUID
	msg []byte
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	id := uuid.New()
	fmt.Printf("%s клиент подключился\n", id)
	defer fmt.Printf("%s клиент отключился\n", id)
	// соединения до WebSocket
	upgrader := websocket.Upgrader{}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Ошибка при получения соединения до WebSocket: %v", err)
		return
	}
	defer ws.Close()

	// соединение до БД
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)
	fmt.Println(dsn)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Ошибка при получения соединения до БД: %v", err)
		return
	}
	defer db.Close()

	// Предварительный разбор запроса
	stmt, err := db.Prepare("INSERT INTO table2(id, data) VALUES ($1, $2)")
	if err != nil {
		log.Fatalf("Ошибка при разборе запроса: %v", err)
		return
	}
	defer stmt.Close()

	// Бесконечный цикл для получения сообщений от клиента.
	for {
		// читаем из сокета
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Ошибка при чтении сообщения: %v", err)
			break
		}
		fmt.Printf("%s пришли новые данные: %s\n", id, msg)
		// декодируем из Base64
		dst := make([]byte, base64.StdEncoding.DecodedLen(len(msg)))
		realLen, err := base64.StdEncoding.Decode(dst, msg)
		if err != nil {
			log.Printf("Ошибка при декодировании сообщения: %v", err)
			continue
		}
		// пишем в DB
		_, err = stmt.Exec(id.String(), dst[:realLen])
		if err != nil {
			log.Printf("Ошибка при записи в БД: %v", err)
			continue
		}

	}
}

func main() {
	// Обработчик второго клиента
	pathParam := os.Getenv("SERVER_PATH_PARAM_CLIENT2")
	http.HandleFunc("/"+pathParam, handleConnections)

	host := os.Getenv("SERVER_HOST")
	port := os.Getenv("SERVER_PORT")
	addr := fmt.Sprintf("%s:%s", host, port)

	// Запуск сервера
	log.Printf("Сервер запущен на порту %s", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}
