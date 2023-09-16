package main

import (
	"log"
	"net/http"
	"github.com/PratikforCoding/BusoFact.git/router"
)

func main() {
	r := router.Router()

	corsMux := middlewareCors(r)
	server := &http.Server {
		Addr: ":8080",
		Handler: corsMux,
	}

	log.Println("Server is running at port: 8080 ....")
	log.Fatal(server.ListenAndServe())
}