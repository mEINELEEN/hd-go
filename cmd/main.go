package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"hd-go/db"
	"hd-go/handlers"
)

func main() {
	if err := db.InitDB(); err != nil {
		log.Fatalf("Не удалось инициализировать БД: %v", err)
	}
	fmt.Println("База данных SQLite успешно подключена")

	// Раздача статического фронтенда
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	// Роуты авторизации
	http.HandleFunc("/api/register", handlers.RegisterHandler)
	http.HandleFunc("/api/login", handlers.LoginHandler)

	// Роуты заявок
	http.HandleFunc("/api/tickets", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			handlers.CreateTicketHandler(w, r)
		} else if r.Method == http.MethodGet {
			handlers.GetTicketsHandler(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/api/tickets/status", handlers.UpdateTicketStatusHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Сервер запущен на http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
