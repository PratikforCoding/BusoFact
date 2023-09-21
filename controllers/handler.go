package controller

import (
	"encoding/json"
	"net/http"
	"strings"

	reply "github.com/PratikforCoding/BusoFact.git/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (apiCfg *APIConfig)HandlerGetBuses(w http.ResponseWriter, r *http.Request) {
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

	buses := apiCfg.getBuses(source, destination)
	reply.RespondWithJson(w, http.StatusFound, buses)
}

func (apiCfg *APIConfig)HandlerAddBuses(w http.ResponseWriter, r *http.Request) {
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
	
	newBus, err := apiCfg.addBuses(params.Name, params.StopageName)
	if err != nil {
		reply.RespondWtihError(w, http.StatusInternalServerError, "Couldn't add the bus")
		return
	}
	reply.RespondWithJson(w, http.StatusOK, newBus)
}

func (apiCfg *APIConfig)HandlerGetBusByName(w http.ResponseWriter, r *http.Request) {
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

	foundBus, err := apiCfg.getBusByName(params.Name)
	if err != nil {
		reply.RespondWtihError(w, http.StatusNotFound, "Bus not found")
		return 
	}
	reply.RespondWithJson(w, http.StatusFound, foundBus)
}

func (apiCfg *APIConfig)HandlerCreateAccount(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Email string `json:"email"`
		Password string `json:"password"`
		ConfirmPassword string `json:"confpassword"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		reply.RespondWtihError(w, http.StatusInternalServerError, "Couldn,t decode parameters")
		return
	}
	
	if params.Password != params.ConfirmPassword {
		reply.RespondWtihError(w, http.StatusNotAcceptable, "confirm password correctly")
		return
	}

	user, err := apiCfg.createUser(params.Username, params.Email, params.Password)
	if err != nil {
		reply.RespondWtihError(w, http.StatusConflict, "User already exists")
		return
	}
	idStr := user["_id"].(primitive.ObjectID).Hex()
	retUser := bson.M {
		"username": user["username"].(string),
		"email": user["email"].(string),
		"id": idStr,
	}
	reply.RespondWithJson(w, http.StatusCreated, retUser)
}

func (apiCfg *APIConfig)HandlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		reply.RespondWtihError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := apiCfg.userLogin(params.Email, params.Password)
	if err != nil {
		errorMsg := "User authentication failed"
		if err.Error() == "user doesn't exist" {
			errorMsg = "User doesn't exist"
		} else if err.Error() == "wrong password" {
			errorMsg = "Wrong password"
		}
		reply.RespondWtihError(w, http.StatusNotFound, errorMsg)
		return
	}
	idStr := user["_id"].(primitive.ObjectID).Hex()
	retUser := bson.M {
		"username": user["username"].(string),
		"email": user["email"].(string),
		"id": idStr,
	}
	reply.RespondWithJson(w, http.StatusOK, retUser)
}

