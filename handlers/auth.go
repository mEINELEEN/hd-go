package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"hd-go/db"

	"golang.org/x/crypto/bcrypt"
)

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"` // 'user' или 'admin'
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Регистрация нового пользователя
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, `{"message":"Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Status: "error", Message: "Некорректный JSON"})
		return
	}

	if req.Username == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Status: "error", Message: "Заполните логин и пароль"})
		return
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Status: "error", Message: "Ошибка обработки пароля"})
		return
	}

	// По умолчанию роль 'user', если не указана 'admin'
	role := "user"
	if req.Role == "admin" {
		role = "admin"
	}

	// Запись в БД
	query := `INSERT INTO users (username, password, role) VALUES (?, ?, ?)`
	_, err = db.DB.Exec(query, req.Username, string(hashedPassword), role)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(Response{Status: "error", Message: "Пользователь уже существует"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{Status: "success", Message: "Пользователь успешно зарегистрирован"})
}

// Вход пользователя
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, `{"message":"Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	var req UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Status: "error", Message: "Некорректный JSON"})
		return
	}

	var storedHash string
	var role string
	var userID int

	query := `SELECT id, password, role FROM users WHERE username = ?`
	err := db.DB.QueryRow(query, req.Username).Scan(&userID, &storedHash, &role)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{Status: "error", Message: "Неверный логин или пароль"})
		return
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(req.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Response{Status: "error", Message: "Неверный логин или пароль"})
		return
	}

	// Возвращаем успешный ответ и роль пользователя
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "success",
		"message": "Успешный вход",
		"user": map[string]interface{}{
			"id":       userID,
			"username": req.Username,
			"role":     role,
		},
	})
}
