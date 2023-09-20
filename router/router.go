package router

import (
	controller "github.com/PratikforCoding/BusoFact.git/controllers"
	"github.com/go-chi/chi/v5"
)
func Router(apiCfg *controller.APIConfig) *chi.Mux {
	router := chi.NewRouter()
	userRouter := chi.NewRouter()

	router.Get("/getbuses", apiCfg.HandlerGetBuses)
	router.Get("/getbusbyname", apiCfg.HandlerGetBusByName)
	router.Post("/addbus", apiCfg.HandlerAddBuses)

	userRouter.Post("/createaccount", apiCfg.HandlerCreateAccount)
	userRouter.Post("/login", apiCfg.HandlerLogin)

	router.Mount("/api", router)
	router.Mount("/usr", userRouter)

	return router
}