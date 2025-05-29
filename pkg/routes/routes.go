package routes

import (
	"github.com/go-chi/chi/v5"
	handler "github.com/saran-pt/livetrade-engine/pkg/handlers"
	// db "github.com/saran-pt/livetrade-engine/pkg/sql/dbal"
)

// var RegisterRoutes = func(router *chi.Mux, conn *sql.DB) {
// 	router.Get("/depth", handler.GetDepth)
// }

func RegisterRoutes(router *chi.Mux, cfg *handler.ApiConfig) {
	router.Post("/users", cfg.CreateUser)
	router.Post("/order", cfg.PlaceOrder)
	router.Get("/balance", cfg.GetBalance)
	router.Get("/depth", cfg.GetDepth)
}
