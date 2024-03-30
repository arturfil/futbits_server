package handlers

import (
	"chi_soccer/helpers"
	"chi_soccer/services"
	// "encoding/json"
	"net/http"
	// "github.com/go-chi/chi/v5"
)

// var field services.Field
var db services.DB

// GET/fields
func getAllFields(w http.ResponseWriter, r *http.Request) {
	// var fields services.Field
	// var db services.DB

	all, err := db.Models.Field.GetAllFields()
	// all, err := db.Models.Field.GetAllFields()
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"fields": all})
}

// GET/fields/field
// func getFieldById(w http.ResponseWriter, r *http.Request) {
// 	id := chi.URLParam(r, "id")
//
//     singleField, err := field.GetFieldById(id)
// 	if err != nil {
// 		helpers.MessageLogs.ErrorLog.Println("\nerror -> ->", err)
// 		return
// 	}
//
// 	helpers.WriteJSON(w, http.StatusOK, singleField)
// }
//
// // POST/fields/field
// func createField(w http.ResponseWriter, r *http.Request) {
// 	err := json.NewDecoder(r.Body).Decode(&field)
//
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
//
// 	helpers.WriteJSON(w, http.StatusOK, field)
// 	err = field.CreateField(field)
//
// 	if err != nil {
// 		helpers.MessageLogs.ErrorLog.Println(err)
//         return
// 	}
//
//
// 	helpers.WriteJSON(w, http.StatusOK, "Created Field sucessfully")
// }
//
// // PUT/field
// func updateField(w http.ResponseWriter, r *http.Request) {
// 	var f services.Field
//
// 	err := json.NewDecoder(r.Body).Decode(&f)
//
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
//
// 	helpers.WriteJSON(w, http.StatusOK, f)
//
// 	err = f.UpdateField()
//
// 	if err != nil {
// 		helpers.MessageLogs.ErrorLog.Println(err)
// 		helpers.WriteJSON(w, http.StatusOK, "updated field")
// 	}
// }
