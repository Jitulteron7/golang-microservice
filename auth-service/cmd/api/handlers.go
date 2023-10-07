package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (app *Config) Auth(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"Email"`
		Password string `json:"Password"`
	}

	err := app.readJSON(w, r, &requestPayload)

	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	user, err := app.Models.User.GetByEmail(requestPayload.Email)

	if err != nil {
		app.errorJSON(w, errors.New("invalid creadentials"), http.StatusBadGateway)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil {
		app.errorJSON(w, errors.New("invalid creadentials"), http.StatusBadGateway)
		return
	}
	if !valid {
		app.errorJSON(w, errors.New("invalid creadentials"), http.StatusBadGateway)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	app.writeJSON(w, http.StatusAccepted, payload)

}

func (app *Config) Hi(w http.ResponseWriter, r *http.Request) {

	payload := jsonResponse{
		Error:   false,
		Message: "Hi from auth",
		Data:    nil,
	}

	err := app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		log.Panic(err)
	}

}
