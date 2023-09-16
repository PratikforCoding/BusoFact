package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	controller "github.com/PratikforCoding/BusoFact.git/controllers"
	reply "github.com/PratikforCoding/BusoFact.git/json"
)

func HandlerGetBuses(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Source string `json:"source"`
		Destination string `json:"destination"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		reply.RespondWtihError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	source := strings.ToLower(params.Source)
	destination := strings.ToLower(params.Destination)

	buses := controller.GetBuses(source, destination)
	reply.RespondWithJson(w, http.StatusFound, buses)
}

func HandlerAddBuses(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
		StopageName string `json:"stopageName"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		reply.RespondWtihError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}
	
	newBus, err := controller.AddBuses(params.Name, params.StopageName)
	if err != nil {
		reply.RespondWtihError(w, http.StatusInternalServerError, "Couldn't add the bus")
		return
	}
	reply.RespondWithJson(w, http.StatusOK, newBus)
}

func HandlerGetBusByName(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		reply.RespondWtihError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	foundBus, err := controller.GetBusByName(params.Name)
	if err != nil {
		reply.RespondWtihError(w, http.StatusNotFound, "Bus not found")
		return 
	}
	reply.RespondWithJson(w, http.StatusFound, foundBus)
}