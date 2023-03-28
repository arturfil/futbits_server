package controllers

import (
	"chi_soccer/helpers"
	"chi_soccer/services"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GET/fields
func GetAllFields(w http.ResponseWriter, r *http.Request) {
	var fields services.Field
	all, err := fields.GetAllFields()
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"fields": all})
}

// GET/fields/field
func GetFieldById(w http.ResponseWriter, r *http.Request) {
	// var field services.Field
	id := chi.URLParam(r, "id")
	singleField, err := mod.Field.GetFieldById(id)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, singleField)
}

// POST/fields/field
func CreateField(w http.ResponseWriter, r *http.Request) {
	var f services.Field
	err := json.NewDecoder(r.Body).Decode(&f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, f)
	id, err := mod.Field.CreateField(f)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		newField, _ := mod.Field.GetFieldById(id)
		helpers.WriteJSON(w, http.StatusOK, newField)
	}
}

// PUT/field
func UpdateField(w http.ResponseWriter, r *http.Request) {
	var f services.Field
	err := json.NewDecoder(r.Body).Decode(&f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, f)
	err = f.UpdateField()
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		helpers.WriteJSON(w, http.StatusOK, "updated field")
	}
}
