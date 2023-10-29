package router

import (
	"github.com/PratikforCoding/BusoFact.git/controllers"
	"github.com/go-chi/chi/v5"
)
func Router(apiCfg *controller.APIConfig) *chi.Mux {
	router := chi.NewRouter()
	apiRouter := chi.NewRouter()
	userRouter := chi.NewRouter()
	
	

	apiRouter.Get("/getbuses", apiCfg.HandlerGetBuses)
	apiRouter.Get("/getbusbyname", apiCfg.HandlerGetBusByName)

	userRouter.Post("/addbus", apiCfg.HandlerAddBuses)
	userRouter.Put("/addstopage", apiCfg.HandlerAddStopage)

	userRouter.Post("/createaccount", apiCfg.HandlerCreateAccount)
	userRouter.Post("/login", apiCfg.HandlerLogin)

    userRouter.Get("/getusers", apiCfg.HandlerGetAllUsers)
	userRouter.Put("/makeadmin", apiCfg.HandlerMakeAdmin)
	userRouter.Delete("/deletebus", apiCfg.HandlerDeleteBus)
	userRouter.Delete("/deleteuser", apiCfg.HandlerDeleteUser)
	
	router.Mount("/api", apiRouter)
	router.Mount("/usr", userRouter)

	return router
}