package main

import (
	"fmt"
	"log"
	"net/http"

	"hd-go/db"
)

func main() {
	// Инициализация базы данных
	if err := db.InitDB(); err != nil {
		log.Fatalf("Не удалось инициализировать БД: %v", err)
	}
	fmt.Println("База данных SQLite успешно подключена")

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "ok", "message": "Helpdesk API is running"}`))
	})

	fmt.Println("Сервер запущен на http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
