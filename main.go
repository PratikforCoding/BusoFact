package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	controller "github.com/PratikforCoding/BusoFact.git/controllers"
	"github.com/PratikforCoding/BusoFact.git/database"
	"github.com/PratikforCoding/BusoFact.git/router"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongouri := os.Getenv("connectlink")

	busCol, usrCol, err := database.CreateDB(mongouri)
	if err != nil {
		log.Fatal("Didn't create connection to mongodb")
	}
	defer database.CloseDB()

	apicfg := controller.NewAPIConfig(busCol, usrCol)

	fmt.Println("MongoDB API")
	r := router.Router(apicfg)

	corsMux := middlewareCors(r)
	server := &http.Server {
		Addr: ":8080",
		Handler: corsMux,
	}

	log.Println("Server is getting started at port: 8080 ....")
	log.Fatal(server.ListenAndServe())
	log.Println("Server is runnig at port: 8080 ....")
}