package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/PratikforCoding/BusoFact.git/auth"
	reply "github.com/PratikforCoding/BusoFact.git/json"
	model "github.com/PratikforCoding/BusoFact.git/models"
)

func (apiCfg *APIConfig)HandlerGetBuses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	source := r.URL.Query().Get("source")
	destination := r.URL.Query().Get("destination")

	buses := apiCfg.getBuses(source, destination)
	reply.RespondWithJson(w, http.StatusFound, buses)
}

func (apiCfg *APIConfig)HandlerAddBuses(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
		StopageName string `json:"stopageName"`
	}

	token, err := auth.GetTokenFromCookie(r, apiCfg.jwtSecret)
	if err != nil {	
		reply.RespondWtihError(w, http.StatusUnauthorized, "Couldn't get token from request")
		return
	}

	claims, err := auth.ValidateJWT(token, apiCfg.jwtSecret)
	if err != nil {
		reply.RespondWtihError(w, http.StatusUnauthorized, "Couldn't validate token")
		return
	}
	fmt.Println(claims)
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
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
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	busName := r.URL.Query().Get("busname")

	foundBus, err := apiCfg.getBusByName(busName)
	if err != nil {
		reply.RespondWtihError(w, http.StatusNotFound, "Bus not found")
		return 
	}
	reply.RespondWithJson(w, http.StatusFound, foundBus)
}

func (apiCfg *APIConfig) HandlerAddStopage(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
		StopageName string `json:"stopageName"`
		BeforeStopage string `json:"beforeStopage"`
	}

	decoder := json.NewDecoder(r.Body);
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		reply.RespondWtihError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	upadtedBus, err := apiCfg.addBusStopage(params.Name, params.StopageName, params.BeforeStopage)
	if err != nil {
		reply.RespondWtihError(w, http.StatusInternalServerError, "couldn't update the bus")
		return
	}
	reply.RespondWithJson(w, http.StatusOK, upadtedBus)

} 

func (apiCfg *APIConfig)HandlerCreateAccount(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		FristName string `json:"firstName"`
		LastName string `json:"lastName"`
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

	user, err := apiCfg.createUser(params.FristName, params.LastName, params.Email, params.Password)

	if err != nil {
		reply.RespondWtihError(w, http.StatusInternalServerError, "error creating user")
		return
	}
	retUser := model.User{
		ID: user.ID,
		FristName: user.FristName,
		LastName: user.LastName,
		Email: user.Email,
		Role: user.Role,
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

	accessTokenTime := 60 * 60
	refreshTokenTime := 60 *24 * 60 * 60
	accessToken, err := auth.MakeAccessToken(user.ID.Hex(), apiCfg.jwtSecret, time.Duration(accessTokenTime) * time.Second)
	if err != nil {
		reply.RespondWtihError(w, http.StatusInternalServerError, "Couldn't not create Access Token")
		return
	}

	refreshToken, err := auth.MakeRefreshToken(user.ID.Hex(), apiCfg.jwtSecret, time.Duration(refreshTokenTime) * time.Second)
	if err != nil {
		reply.RespondWtihError(w, http.StatusInternalServerError, "Couldn't not create Access Token")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HttpOnly: true,
		Secure:   true, 
		SameSite: http.SameSiteLaxMode,
	})
	
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true, 
		SameSite: http.SameSiteLaxMode,
	})
	

	retUser := model.User{
		ID: user.ID,
		FristName: user.FristName,
		LastName: user.LastName,
		Email: user.Email,
	}
	reply.RespondWithJson(w, http.StatusOK, retUser)
}

func (apiCfg *APIConfig) HandlerGetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := apiCfg.getAllUsers()
	if err != nil {
		reply.RespondWtihError(w, http.StatusNotFound, "couldn't get users")
		return 
	}
	reply.RespondWithJson(w, http.StatusFound, users)
}

func (apiCfg *APIConfig) HandlerMakeAdmin(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	user, err := apiCfg.makeAdmin(email)
	if err != nil {
		reply.RespondWtihError(w, http.StatusBadRequest, "couldn't promote to admin")
		return 
	}
	retUser := model.User{
		ID: user.ID,
		FristName: user.FristName,
		LastName: user.LastName,
		Email: user.Email,
		Role: user.Role,
	}
	reply.RespondWithJson(w, http.StatusOK, retUser)
}



