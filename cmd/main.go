package main

import (
	"fmt"
	"log"
	"net/http"

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

	fmt.Println("Сервер запущен на http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
