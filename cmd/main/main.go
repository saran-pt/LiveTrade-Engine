package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	_ "github.com/lib/pq"
	"github.com/saran-pt/livetrade-engine/pkg/config"
	handler "github.com/saran-pt/livetrade-engine/pkg/handlers"
	"github.com/saran-pt/livetrade-engine/pkg/routes"
	"github.com/saran-pt/livetrade-engine/pkg/sql/dbal"
)

func main() {
	config.LoadEnv()
	port := config.GetEnv("PORT", "8080")

	// DataBase Configuration
	dbURL, ok := os.LookupEnv("DB_URL")
	// dbURL, ok := "postgres://postgres:password@localhost:5432/livetradeengine", true
	if !ok {
		log.Fatal("DB-URL not found!")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to Database: ", err)
	}

	apiConf := &handler.ApiConfig{DB: dbal.New(conn)}

	defer conn.Close()

	// Router Configuration
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	// V1 Router Configuration
	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	routes.RegisterRoutes(router, apiConf)
	log.Fatal(srv.ListenAndServe())
}
