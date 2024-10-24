package routes

import (
	"database/sql"
	"net/http"
	"strings"
	handlers "test_kerja/app_core/handler"
	routes "test_kerja/app_core/midleware"
	"test_kerja/app_core/utils"
)

// SetupRoutes sets up the HTTP routes
func SetupRoutes(db *sql.DB) {
	http.HandleFunc("/user", routes.BasicAuthMiddleware(db, func(w http.ResponseWriter, r *http.Request) {
		utils.Route(w, r, utils.RouteStruct{
			Get:  handlers.GetAll(db), // Handle GET request for /user
			Post: handlers.Create(db),
		}, db)
	}))

	http.HandleFunc("/user/", routes.BasicAuthMiddleware(db, func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/user/")
		idParts := strings.Split(id, "/")
		if len(idParts) != 1 {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		utils.Route(w, r, utils.RouteStruct{
			Get:    handlers.Get(db, idParts[0]),    // Handle GET request for /user/{id}
			Update: handlers.Update(db, idParts[0]), // Handle PUT request for /user/{id}
			Delete: handlers.Delete(db, idParts[0]), // Handle DELETE request for /user/{id}
		}, db)
	}))
}
