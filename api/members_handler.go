package main

import (
	"chi_soccer/internal/data"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// GET/members/group_id
func (app *application) GetAllMembersFromGroup(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "group_id"))
	app.infoLog.Println("id", id)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	members, err := app.models.Member.GetAllMemberOfAGroup(id)
	app.infoLog.Println(members)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	app.writeJSON(w, http.StatusOK, members)
}

// POST/members/create
func (app *application) CreateMember(w http.ResponseWriter, r *http.Request) {
	var m data.Member
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := app.models.Member.CreateMember(m)
	if err != nil {
		app.errorLog.Println(err)
	}
	app.writeJSON(w, http.StatusOK, id)
	// id, err := app.models.Member.CreateMember()
}
