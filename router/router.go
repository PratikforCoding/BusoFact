package router

import (
	"github.com/PratikforCoding/BusoFact.git/controllers"
	"github.com/go-chi/chi/v5"
)
func Router(apiCfg *controller.APIConfig) *chi.Mux {
	router := chi.NewRouter()
	apiRouter := chi.NewRouter()
	userRouter := chi.NewRouter()
	adminRouter := chi.NewRouter()
	

	apiRouter.Get("/getbuses", apiCfg.HandlerGetBuses)
	apiRouter.Get("/getbusbyname", apiCfg.HandlerGetBusByName)

	userRouter.Put("/addbus", apiCfg.HandlerAddBuses)
	userRouter.Put("/addstopage", apiCfg.HandlerAddStopage)

	userRouter.Post("/createaccount", apiCfg.HandlerCreateAccount)
	userRouter.Post("/login", apiCfg.HandlerLogin)

	adminRouter.Put("/makeadmin", apiCfg.HandlerMakeAdmin)
	adminRouter.Get("/getusers", apiCfg.HandlerGetAllUsers)

	router.Mount("/api", apiRouter)
	router.Mount("/usr", userRouter)
	router.Mount("/admin", adminRouter)

	return router
}