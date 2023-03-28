package controllers

import (
	"chi_soccer/helpers"
	"chi_soccer/services"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func CreateProfile(w http.ResponseWriter, r *http.Request) {
	var p services.Profile
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := mod.Profile.CreateProfile(p)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
	}
	newProfile, _ := mod.Profile.GetProfileById(id)
	helpers.WriteJSON(w, http.StatusOK, newProfile)
}

func GetProfileById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	profile, err := mod.Profile.GetProfileById(id)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, profile)
}
