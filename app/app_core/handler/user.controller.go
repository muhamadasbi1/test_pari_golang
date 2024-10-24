package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"test_kerja/app_core/dto"
	models "test_kerja/app_core/model"
	"test_kerja/app_core/utils"
	"time"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()
var table = "auth_users"

func GetAll(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		var user models.User
		stmt, err := db.PrepareContext(ctx, "SELECT id, name, email FROM "+table+" ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer stmt.Close()
		err = stmt.QueryRowContext(ctx).Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "user not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func Create(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.UserCreateRequest
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
		stmt, err := tx.Prepare("INSERT INTO " + table + " (name, email) VALUES (?, ?)")
		if utils.CekError(w, err) {
			tx.Rollback()
			return
		}
		defer stmt.Close()
		result, err := stmt.Exec(req.Name, req.Email)
		if utils.CekError(w, err) {
			tx.Rollback()
			return
		}
		id, _ := result.LastInsertId()
		response := map[string]interface{}{
			"ID":    int(id),
			"Name":  req.Name,
			"Email": req.Email,
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

		var user models.User
		stmt, err := db.PrepareContext(ctx, "SELECT id, name, email FROM "+table+" WHERE id = ?")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer stmt.Close()
		err = stmt.QueryRowContext(ctx, id).Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "user not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func Update(db *sql.DB, id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		tx, err := db.BeginTx(ctx, nil)
		if utils.CekError(w, err) {
			return
		}

		stmt, err := tx.Prepare("UPDATE " + table + " SET name = ?, email = ? WHERE id = ?")
		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(user.Name, user.Email, id)
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
