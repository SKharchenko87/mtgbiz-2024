package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	mtgbiz_2024 "github.com/SKharchenko87/mtgbiz-2024"

	"github.com/gorilla/websocket"
)

// task1 - получает данные в одном потоке и передаём для записи в файл в другой поток. Данные из потока в поток передаются построчно.
func task1(ch chan mtgbiz_2024.Table1) {
	log.Println("Клиент запущен")
	defer log.Println("Клиент остановлен")
	// Подключение к серверу
	conn, _, err := websocket.DefaultDialer.Dial(mtgbiz_2024.GetServerURL("CLIENT1"), nil)
	if err != nil {
		log.Printf("Ошибка при подключении к серверу: %v", err)
		return
	}
	defer conn.Close()

	// Получаем сообщения от сервера.
	for {
		// читаем из сокета
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Printf("Ошибка при чтении сообщения: %v", err)
			}
			break
		}
		fmt.Println(string(msg))

		// десереализуем
		t := mtgbiz_2024.Table1{}
		err = json.Unmarshal(msg, &t)
		if err != nil {
			log.Printf("Ошибка при декодировании сообщения: %v", err)
			continue
		}

		// пишем в канал на запись в файл
		ch <- t
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

func main() {

	ch := make(chan mtgbiz_2024.Table1) // канал для записи в файл

	http.Handle("GET /", http.HandlerFunc(mainPage))
	http.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		go task1(ch) // Нажатие на кнопку не должно замораживать клиент
		mainPage(w, r)
	})

	// Читаем из канала и отправляем в файл
	go func() {
		f, err := os.Create("/client1.out")
		if err != nil {
			log.Printf("Ошибка открытия файла данных: %v", err)
			return
		}

		for t := range ch {
			_, err = f.WriteString(fmt.Sprint(t))
			if err != nil {
				log.Printf("Ошибка записи в файл: %v", err)
				continue
			}
		}
	}()

	host := os.Getenv("CLIENT1_HOST")
	port := os.Getenv("CLIENT1_PORT")
	addr := fmt.Sprintf("%s:%s", host, port)
	log.Fatal(http.ListenAndServe(addr, nil))
}
