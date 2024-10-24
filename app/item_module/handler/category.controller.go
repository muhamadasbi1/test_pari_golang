package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"test_kerja/app_core/utils"
	"test_kerja/item_module/dto"
	models "test_kerja/item_module/model"
	"time"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()
var table = "cm_category"

func GetAll(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		var categories []models.Category
		rows, err := db.QueryContext(ctx, "SELECT id, name, created_at, updated_at FROM "+table)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var category models.Category
			if err := rows.Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			categories = append(categories, category)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(categories)
	}
}

func Create(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.CategoryCreateRequest
		if err := utils.DecodeAndValidate(r, &req, validate); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()
		tx, err := db.BeginTx(ctx, nil)
		if utils.CekError(w, err) {
			return
		}
		stmt, err := tx.Prepare("INSERT INTO " + table + " (name) VALUES (?)")
		if utils.CekError(w, err) {
			tx.Rollback()
			return
		}
		defer stmt.Close()
		result, err := stmt.Exec(req.Name)
		if utils.CekError(w, err) {
			tx.Rollback()
			return
		}
		id, _ := result.LastInsertId()
		response := map[string]interface{}{
			"ID":   int(id),
			"Name": req.Name,
		}
		if err = tx.Commit(); utils.CekError(w, err) {
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}

func Get(db *sql.DB, id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		var category models.Category
		stmt, err := db.PrepareContext(ctx, "SELECT id, name FROM "+table+" WHERE id = ?")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer stmt.Close()
		err = stmt.QueryRowContext(ctx, id).Scan(&category.ID, &category.Name)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "category not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(category)
	}
}

func Update(db *sql.DB, id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.CategoryUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		tx, err := db.BeginTx(ctx, nil)
		if utils.CekError(w, err) {
			return
		}

		stmt, err := tx.Prepare("UPDATE " + table + " SET name = ? WHERE id = ?")
		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(req.Name, id)
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

func Delete(db *sql.DB, id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		stmt, err := tx.Prepare("DELETE FROM " + table + " WHERE id = ?")
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
