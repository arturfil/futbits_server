package main

import (
	"chi_soccer/internal/data"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GET/groups
func (app *application) GetAllGroups(w http.ResponseWriter, r *http.Request) {
	var groups data.Group
	all, err := groups.GetAllGroups()
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	app.writeJSON(w, http.StatusOK, envelope{"groups": all})
}

// GET/groups/group/:id
func (app *application) GetGroupById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	group, err := app.models.Group.GetGroupById(id)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	app.writeJSON(w, http.StatusOK, group)

}

// POST/groups/create
func (app *application) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var g data.Group // 0xffddc234 \\ {name: "Manchester Utd", id: 1234}
	err := json.NewDecoder(r.Body).Decode(&g)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := app.models.Group.CreateGroup(g)
	if err != nil {
		app.errorLog.Println(err)
	}
	newGroup, _ := app.models.Group.GetGroupById(id)
	app.writeJSON(w, http.StatusOK, newGroup)
}
