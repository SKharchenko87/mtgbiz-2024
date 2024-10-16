package main

import (
	"encoding/json"
	"fmt"
	"github.com/SKharchenko87/mtgbiz-2024/common"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"os"
)

func getServerURL() string {
	protocol := os.Getenv("SERVER_PROTOCOL_CLIENT1")
	host := os.Getenv("SERVER_HOST_CLIENT1")
	port := os.Getenv("SERVER_PORT_CLIENT1")
	pathParam := os.Getenv("SERVER_PATH_PARAM_CLIENT1")
	serverURL := fmt.Sprintf("%s://%s:%s/%s", protocol, host, port, pathParam)
	return serverURL
}

// ToDO
func send() {
	log.Println("Клиент запущен")
	// Подключение к серверу
	conn, _, err := websocket.DefaultDialer.Dial(getServerURL(), nil)
	if err != nil {
		log.Fatalf("Ошибка при подключении к серверу: %v", err)
	}
	defer conn.Close()

	// Получаем сообщения от сервера.
	for {
		// читаем из сокета
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Ошибка при чтении сообщения: %v", err)
			break
		}
		if mt == websocket.CloseMessage {
			break
		}
		fmt.Printf("%s\n", msg) //ToDo
		// десереализуем
		t := common.Table1{}
		err = json.Unmarshal(msg, t)
		if err != nil {
			log.Printf("Ошибка при декодировании сообщения: %v", err)
			continue
		}
		// пишем в файл ToDO
		fmt.Printf("%s\n", t) //ToDo

	}

}

// Обработчик главной страницы клиента
func mainPage(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("index.html")
	if err != nil {
		log.Printf("Ошибка открытия шаблона страницы: %v", err)
		return
	}

	data, err := io.ReadAll(f)
	if err != nil {
		log.Printf("Ошибка при чтении шаблона страницы: %v", err)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		log.Printf("Ошибка передачи данных на страницу: %v", err)
		return
	}

}

// Запрос данных
func getData(w http.ResponseWriter, r *http.Request) {

	go send() // Нажатие на кнопку не должно замораживать клиент
	mainPage(w, r)
}

type serverParam struct {
	protocol  string
	host      string
	port      string
	pathParam string
}

func main() {
	http.Handle("GET /", http.HandlerFunc(mainPage))
	http.Handle("POST /", http.HandlerFunc(getData))
	host := os.Getenv("CLIENT2_HOST")
	port := os.Getenv("CLIENT2_PORT")
	addr := fmt.Sprintf("%s:%s", host, port)
	log.Fatal(http.ListenAndServe(addr, nil))

}
