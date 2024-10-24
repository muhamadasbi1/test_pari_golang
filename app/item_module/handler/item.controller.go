package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"test_kerja/app_core/utils"
	"test_kerja/item_module/dto"
	models "test_kerja/item_module/model"
	"time"
)

var table_item = "tr_items"

type PaginatedResponse struct {
	Items       []models.Item `json:"items"`
	TotalItems  int           `json:"total_items"`
	TotalPages  int           `json:"total_pages"`
	CurrentPage int           `json:"current_page"`
}

func GetAllItems(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		// Ambil query parameters
		pageStr := r.URL.Query().Get("page")
		limitStr := r.URL.Query().Get("limit")
		sort := r.URL.Query().Get("sort")
		order := r.URL.Query().Get("order")
		search := r.URL.Query().Get("search")

		// Default values
		page := 1
		limit := 10
		if pageStr != "" {
			var err error
			page, err = strconv.Atoi(pageStr)
			if err != nil || page < 1 {
				page = 1
			}
		}
		if limitStr != "" {
			var err error
			limit, err = strconv.Atoi(limitStr)
			if err != nil || limit < 1 {
				limit = 10
			}
		}

		// Set sorting defaults
		if sort == "" {
			sort = "id" // Default sort by id
		}
		if order == "" {
			order = "ASC" // Default order
		}

		// Prepare the SQL query for counting total items
		countQuery := `SELECT COUNT(*) FROM tr_items`
		var totalItems int
		if search != "" {
			countQuery += " WHERE name LIKE ? OR description LIKE ?"
			err := db.QueryRowContext(ctx, countQuery, "%"+search+"%", "%"+search+"%").Scan(&totalItems)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			err := db.QueryRowContext(ctx, countQuery).Scan(&totalItems)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Calculate total pages
		totalPages := (totalItems + limit - 1) / limit

		// Prepare the SQL query for fetching items
		query := `
			SELECT 
				tr_items.id, 
				tr_items.category_id, 
				tr_items.name, 
				tr_items.description, 
				tr_items.price, 
				tr_items.created_at,
				cm_category.name AS category_name
			FROM 
				tr_items
			JOIN 
				cm_category ON tr_items.category_id = cm_category.id
		`

		var args []interface{}
		if search != "" {
			query += " WHERE tr_items.name LIKE ? OR tr_items.description LIKE ?"
			args = append(args, "%"+search+"%", "%"+search+"%")
		}

		query += fmt.Sprintf(" ORDER BY %s %s LIMIT ?, ?", sort, order)
		args = append(args, (page-1)*limit, limit)

		// Execute the query
		rows, err := db.QueryContext(ctx, query, args...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var items []models.Item
		for rows.Next() {
			var item models.Item
			if err := rows.Scan(&item.ID, &item.CategoryID, &item.Name, &item.Description, &item.Price, &item.CreatedAt, &item.CategoryName); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			items = append(items, item)
		}

		// Create the paginated response
		response := PaginatedResponse{
			Items:       items,
			TotalItems:  totalItems,
			TotalPages:  totalPages,
			CurrentPage: page,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
func CreateItem(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.ItemCreateRequest
		if err := utils.DecodeAndValidate(r, &req, validate); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		stmt, err := tx.Prepare("INSERT INTO " + table_item + " (category_id, name, description, price) VALUES (?, ?, ?, ?)")
		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		result, err := stmt.Exec(req.CategoryID, req.Name, req.Description, req.Price)
		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id, _ := result.LastInsertId()
		response := map[string]interface{}{
			"ID":          int(id),
			"CategoryID":  req.CategoryID,
			"Name":        req.Name,
			"Description": req.Description,
			"Price":       req.Price,
		}

		if err = tx.Commit(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

func GetItem(db *sql.DB, id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		var item models.Item
		stmt, err := db.PrepareContext(ctx, "SELECT id, category_id, name, description, price, created_at FROM "+table_item+" WHERE id = ?")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		err = stmt.QueryRowContext(ctx, id).Scan(&item.ID, &item.CategoryID, &item.Name, &item.Description, &item.Price, &item.CreatedAt)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "item not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(item)
	}
}

func UpdateItem(db *sql.DB, id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.ItemUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		stmt, err := tx.Prepare("UPDATE " + table_item + " SET category_id = ?, name = ?, description = ?, price = ? WHERE id = ?")
		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(req.CategoryID, req.Name, req.Description, req.Price, id)
		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = tx.Commit(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func DeleteItem(db *sql.DB, id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		stmt, err := tx.Prepare("DELETE FROM " + table_item + " WHERE id = ?")
		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(id)
		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = tx.Commit(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
