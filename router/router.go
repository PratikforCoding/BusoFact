package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/PratikforCoding/BusoFact.git/handlers"
)
func Router() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/getbuses", handler.HandlerGetBuses)
	router.Get("/getbusbyname", handler.HandlerAddBuses)
	router.Post("/addbus", handler.HandlerAddBuses)

	router.Mount("/api", router)

	return router
}