package main

import (
	"encoding/base64"
	"fmt"
	mtgbiz_2024 "github.com/SKharchenko87/mtgbiz-2024"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	minBlockSize = 1000         // minBlockSize - минимальный размер блока данных в байтах
	maxBlockSize = 10000        // minBlockSize - максимальный размер блока данных в байтах
	srcData      = "1234567890" // srcData - вид данных
)

// send - шлет бесконечно из разных потоков (numberOfThreads - число потоков)
// данные на сервер вида 1234567890 (блоками от 1000 до 10000 байт рандомного размера),
// но перед этим кодирует их в Base64.
func send(numberOfThreads int) {
	// Подключение к серверу
	conn, _, err := websocket.DefaultDialer.Dial(mtgbiz_2024.GetServerURL("CLIENT2"), nil)
	if err != nil {
		log.Fatalf("Ошибка при подключении к серверу: %v", err)
	}
	defer conn.Close()

	ch := make(chan []byte)
	for i := 0; i < numberOfThreads; i++ {
		go func(ch chan []byte, n int) {
			for {
				// формируем данные блоками от 1000 до 10000 байт рандомного размера
				size := rand.Intn(maxBlockSize-minBlockSize) + minBlockSize
				data := make([]byte, size)
				for j := 0; j < size; j++ {
					data[j] = srcData[j%10]
				}
				// кодирует в Base64
				dst := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
				base64.StdEncoding.Encode(dst, data)
				ch <- dst
				time.Sleep(1000 * time.Millisecond) // ToDo что бы не забить быстро БД
			}
		}(ch, i)
	}

	// Читаем из канала и отправляем данные
	for bytes := range ch {
		err = conn.WriteMessage(websocket.TextMessage, bytes)
		if err != nil {
			log.Printf("Ошибка при отправке сообщения: %v", err)
			return
		}
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

// Запуск отправки данных
func start(w http.ResponseWriter, r *http.Request) {
	numberOfThreadsStr := r.FormValue("numberOfThreads")
	numberOfThreads, err := strconv.Atoi(numberOfThreadsStr)
	if err != nil {
		log.Fatalf("Ошибка числа потоков: %v", err)
	}
	go send(numberOfThreads) // Нажатие на кнопку не должно замораживать клиент
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
	http.Handle("POST /", http.HandlerFunc(start))
	host := os.Getenv("CLIENT2_HOST")
	port := os.Getenv("CLIENT2_PORT")
	addr := fmt.Sprintf("%s:%s", host, port)
	log.Fatal(http.ListenAndServe(addr, nil))

}
