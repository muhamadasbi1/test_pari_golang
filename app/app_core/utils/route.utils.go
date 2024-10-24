package utils

import (
	"database/sql"
	"net/http"
	routes "test_kerja/app_core/midleware"
)

type RouteStruct struct {
	Get    http.HandlerFunc
	Post   http.HandlerFunc
	Update http.HandlerFunc
	Delete http.HandlerFunc
}

func Route(w http.ResponseWriter, r *http.Request, route RouteStruct, db *sql.DB) {
	switch r.Method {
	case http.MethodGet:
		if route.Get != nil {
			route.Get(w, r)
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	case http.MethodPost:
		if route.Post != nil {
			routes.BasicAuthMiddleware(db, route.Post)(w, r)
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	case http.MethodPut:
		if route.Update != nil {
			routes.BasicAuthMiddleware(db, route.Post)(w, r)
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	case http.MethodDelete:
		if route.Delete != nil {
			routes.BasicAuthMiddleware(db, route.Post)(w, r)
		} else {
			http.Error(w, "Not found", http.StatusNotFound)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
