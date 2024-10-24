package routes_item

import (
	"database/sql"
	"net/http"
	"strings"
	"test_kerja/app_core/utils"
	handlers "test_kerja/item_module/handler"
)

func SetupRoutes(db *sql.DB) {
	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		utils.Route(w, r, utils.RouteStruct{
			Get:  handlers.GetAll(db),
			Post: handlers.Create(db), // POST handler yang akan diautentikasi
		}, db)
	})

	http.HandleFunc("/categories/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/categories/")
		idParts := strings.Split(id, "/")
		if len(idParts) != 1 {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		utils.Route(w, r, utils.RouteStruct{
			Get:    handlers.Get(db, idParts[0]),    // Pass id to Update
			Update: handlers.Update(db, idParts[0]), // Pass id to Update
			Delete: handlers.Delete(db, idParts[0]), // Pass id to Update
		}, db)
	})

	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		utils.Route(w, r, utils.RouteStruct{
			Get:  handlers.GetAllItems(db), // Handle GET request for /user
			Post: handlers.CreateItem(db),
		}, db)
	})

	http.HandleFunc("/items/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/items/")
		idParts := strings.Split(id, "/")
		if len(idParts) != 1 {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		utils.Route(w, r, utils.RouteStruct{
			Get:    handlers.GetItem(db, idParts[0]),    // Pass id to Update
			Update: handlers.UpdateItem(db, idParts[0]), // Pass id to Update
			Delete: handlers.DeleteItem(db, idParts[0]), // Pass id to Update
		}, db)
	})

}
