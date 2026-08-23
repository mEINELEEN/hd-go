package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"hd-go/db"
)

type TicketRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      int    `json:"user_id"`
}

type StatusRequest struct {
	Status        string `json:"status"`
	AdminResponse string `json:"admin_response"`
}

func CreateTicketHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req TicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Некорректный JSON"})
		return
	}

	if req.Title == "" || req.Description == "" || req.UserID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Заполните все поля"})
		return
	}

	query := `INSERT INTO tickets (title, description, status, user_id) VALUES (?, ?, 'new', ?)`
	_, err := db.DB.Exec(query, req.Title, req.Description, req.UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Ошибка создания заявки"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Заявка успешно создана"})
}

func GetTicketsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	role := r.URL.Query().Get("role")
	userIDStr := r.URL.Query().Get("user_id")

	var rows *sql.Rows
	var err error

	if role == "admin" {
		rows, err = db.DB.Query(`
			SELECT t.id, t.title, t.description, t.status, COALESCE(t.admin_response, ''), t.user_id, u.username, t.created_at 
			FROM tickets t JOIN users u ON t.user_id = u.id ORDER BY t.created_at DESC`)
	} else {
		userID, _ := strconv.Atoi(userIDStr)
		rows, err = db.DB.Query(`
			SELECT t.id, t.title, t.description, t.status, t.user_id, u.username, t.created_at 
			FROM tickets t JOIN users u ON t.user_id = u.id WHERE t.user_id = ? ORDER BY t.created_at DESC`, userID)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Ошибка получения заявок"})
		return
	}
	defer rows.Close()

	type TicketResponse struct {
		ID            int    `json:"id"`
		Title         string `json:"title"`
		Description   string `json:"description"`
		Status        string `json:"status"`
		UserID        int    `json:"user_id"`
		Username      string `json:"username"`
		CreatedAt     string `json:"created_at"`
		AdminResponse string `json:"admin_response"`
	}

	var tickets []TicketResponse
	for rows.Next() {
		var t TicketResponse
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.AdminResponse, &t.UserID, &t.Username, &t.CreatedAt); err != nil {
			continue
		}
		tickets = append(tickets, t)
	}

	json.NewEncoder(w).Encode(tickets)
}

func UpdateTicketStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPatch {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ticketIDStr := r.URL.Query().Get("id")
	ticketID, err := strconv.Atoi(ticketIDStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Некорректный ID заявки"})
		return
	}

	var req StatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Некорректный JSON"})
		return
	}

	_, err = db.DB.Exec("UPDATE tickets SET status = ? admin_response = ? WHERE id = ?", req.Status, ticketID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Ошибка обновления статуса"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Статус заявки успешно обновлен"})
}
