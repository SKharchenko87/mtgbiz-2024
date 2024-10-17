package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	mtgbiz_2024 "github.com/SKharchenko87/mtgbiz-2024"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

func handleConnections1(w http.ResponseWriter, r *http.Request) {
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
	db, err := sql.Open("postgres", mtgbiz_2024.GetDSN())
	if err != nil {
		log.Fatalf("Ошибка при получения соединения до БД: %v", err)
		return
	}
	defer db.Close()

	// Предварительный разбор запроса
	rows, err := db.Query("SELECT id, n, code, data, create_dttm FROM table1")
	if err != nil {
		log.Fatalf("Ошибка при разборе запроса: %v", err)
		return
	}

	//Построчное чтение и отправка данных
	for rows.Next() {
		t := mtgbiz_2024.Table1{}
		err = rows.Scan(&t.Id, &t.N, &t.Code, &t.Data, &t.CreateDttm)
		if err != nil {
			log.Printf("Ошибка при чтение строки: %v", err)
			continue
		}

		b, err := json.Marshal(t)
		if err != nil {
			log.Printf("Ошибка при формировании сообщения: %v", err)
			continue
		}

		err = ws.WriteMessage(websocket.TextMessage, b)
		if err != nil {
			log.Printf("Ошибка при отправке сообщения: %v", err)
			continue
		}
	}
	fmt.Println("xxxxxxx")
	err = ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Printf("Ошибка при отправке сообщения: %v", err)
	}
	time.Sleep(1 * time.Second)

}

func handleConnections2(w http.ResponseWriter, r *http.Request) {
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
	db, err := sql.Open("postgres", mtgbiz_2024.GetDSN())
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
	pathParam := os.Getenv("SERVER_PATH_PARAM_CLIENT1")
	http.HandleFunc("/"+pathParam, handleConnections1)

	// Обработчик второго клиента
	pathParam = os.Getenv("SERVER_PATH_PARAM_CLIENT2")
	http.HandleFunc("/"+pathParam, handleConnections2)

	host := os.Getenv("SERVER_HOST")
	port := os.Getenv("SERVER_PORT")
	addr := fmt.Sprintf("%s:%s", host, port)

	// Запуск сервера
	log.Printf("Сервер запущен на порту %s", port)
	log.Fatal(http.ListenAndServe(addr, nil))
}
