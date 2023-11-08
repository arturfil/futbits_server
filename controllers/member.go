package controllers

import (
	"chi_soccer/helpers"
	"chi_soccer/services"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

var member services.Member

// GET/members/group_id
func GetAllMembersFromGroup(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "group_id")
	members, err := member.GetAllMemberOfAGroup(id)
	helpers.MessageLogs.InfoLog.Println(members)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, members)
}

// POST/members/create
func CreateMember(w http.ResponseWriter, r *http.Request) {
	var m services.Member
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := member.CreateMember(m)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
	}
	helpers.WriteJSON(w, http.StatusOK, id)
	// id, err := app.services.Member.CreateMember()
}
