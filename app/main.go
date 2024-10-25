package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	routes "test_kerja/app_core/route"
	"test_kerja/app_core/utils"
	routes_item "test_kerja/item_module/route"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// CORS middleware
func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	dbUser := getEnv("DB_USER", "test_kerja")
	dbPassword := getEnv("DB_PASSWORD", "TKLqNHoNmjWTnWD@test_kerja")
	dbName := getEnv("DB_NAME", "test_kerja")
	dbHost := getEnv("DB_HOST", "mysql")
	dbPort := getEnv("DB_PORT", "3306")
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	// dataSource := "test_kerja:TKLqNHoNmjWTnWD@test_kerja@tcp(mysql:3306)/test_kerja" // Perbaiki format datasource

	log.Println("Starting the server...")
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	log.Println("Received request for:", r.URL.Path)
	// 	w.Write([]byte("Hello, world!"))
	// })
	database, err := utils.InitDB(dataSource)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	migrationFiles := []string{
		"app_core/migration/create_user_table.sql",
		"app_core/migration/create_item_category_table.sql",
		"app_core/migration/create_item_table.sql",
		"app_core/migration/create_user.sql",
	}

	for _, file := range migrationFiles {
		err = utils.Migrate(database, file)
		if err != nil {
			log.Fatal("Migration failed:", err)
		}
	}

	routes.SetupRoutes(database)
	routes_item.SetupRoutes(database)

	log.Println("Server is running on :8080")
	// http.Handle("/", cors(http.DefaultServeMux)) // Wrap the default mux with CORS
	log.Fatal(http.ListenAndServe(":8080", nil))
}
