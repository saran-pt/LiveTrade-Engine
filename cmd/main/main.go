package main

import (
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/saran-pt/livetrade-engine/pkg/config"
	"github.com/saran-pt/livetrade-engine/pkg/routes"
)

func main() {
	config.LoadEnv()
	port := config.GetEnv("PORT", "8080")

	fmt.Printf(
		"Server running on Port: %s\n",
		port,
	)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	routes.RegisterRoutes(v1Router)
}
