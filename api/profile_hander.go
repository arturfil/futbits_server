package main

import (
	"chi_soccer/internal/data"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) CreateProfile(w http.ResponseWriter, r *http.Request) {
	var p data.Profile
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := app.models.Profile.CreateProfile(p)
	if err != nil {
		app.errorLog.Println(err)
	}
	newProfile, _ := app.models.Profile.GetProfileById(id)
	app.writeJSON(w, http.StatusOK, newProfile)
}

func (app *application) GetProfileById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	profile, err := app.models.Profile.GetProfileById(id)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	app.writeJSON(w, http.StatusOK, profile)
}
