package router

import (
	controller "github.com/PratikforCoding/BusoFact.git/controllers"
	"github.com/go-chi/chi/v5"
)
func Router(apiCfg *controller.APIConfig) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/getbuses", apiCfg.HandlerGetBuses)
	router.Get("/getbusbyname", apiCfg.HandlerGetBusByName)
	router.Post("/addbus", apiCfg.HandlerAddBuses)

	router.Mount("/api", router)

	return router
}