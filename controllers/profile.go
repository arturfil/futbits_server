package controllers

import (
	"chi_soccer/helpers"
	"chi_soccer/services"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

var profile services.Profile

func CreateProfile(w http.ResponseWriter, r *http.Request) {
	var p services.Profile

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := profile.CreateProfile(p)

    log.Println("ID", id)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
	}

	newProfile, _ := profile.GetProfileById(id)
	helpers.WriteJSON(w, http.StatusOK, newProfile)
}

func GetProfileById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	profile, err := profile.GetProfileById(id)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, profile)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")

    var p services.Profile

    err := json.NewDecoder(r.Body).Decode(&p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    helpers.WriteJSON(w, http.StatusOK, p)

    err = p.UpdateProfile(id, p)
    if err != nil {
        helpers.MessageLogs.ErrorLog.Println(err)
        helpers.WriteJSON(w, http.StatusOK, "updated profile")
    }
}
