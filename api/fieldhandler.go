package main

import (
	"chi_soccer/internal/data"
	"encoding/json"
	"net/http"
)

func (app *application) GetAllFields(w http.ResponseWriter, r *http.Request) {
	var fields data.Field
	all, err := fields.GetAllFields()
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	app.writeJSON(w, http.StatusOK, envelope{"fields": all})
}

func (app *application) CreateField(w http.ResponseWriter, r *http.Request) {
	var f data.Field
	err := json.NewDecoder(r.Body).Decode(&f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	app.writeJSON(w, http.StatusOK, f)
	id, err := app.models.Field.CreateField(f)
	if err != nil {
		app.errorLog.Println(err)
		newField, _ := app.models.Field.GetFieldById(id)
		app.writeJSON(w, http.StatusOK, newField)
	}
}

func (app *application) UpdateField(w http.ResponseWriter, r *http.Request) {
	var f data.Field
	err := json.NewDecoder(r.Body).Decode(&f)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	app.writeJSON(w, http.StatusOK, f)
	err = f.UpdateField()
	if err != nil {
		app.errorLog.Println(err)
		app.writeJSON(w, http.StatusOK, "updated field")
	}

}
