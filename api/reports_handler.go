package main

import (
	"chi_soccer/internal/data"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// GET/reports/report
func (app *application) GetAllReports(w http.ResponseWriter, r *http.Request) {
	var reports data.Report
	all, err := reports.GetAllReports()
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	app.writeJSON(w, http.StatusOK, envelope{"reports": all})
}

// GET/reports/report/id
func (app *application) GetReportById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	report, err := app.models.Report.GetReportById(id)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	app.writeJSON(w, http.StatusOK, report)
}

// GET/reports/user_id
func (app *application) GetReportsOfUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "user_id")
	reports, err := app.models.Report.GetAllReporstById(id)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	app.writeJSON(w, http.StatusOK, reports)
}

// POST/reports/report
func (app *application) CreateReport(w http.ResponseWriter, r *http.Request) {
	var rp data.Report
	err := json.NewDecoder(r.Body).Decode(&rp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := app.models.Report.CreateReport(rp)
	if err != nil {
		app.errorLog.Println(err)
	}
	newGroup, _ := app.models.Report.GetReportById(id)
	app.writeJSON(w, http.StatusOK, newGroup)
}

// func (app *application) GetGroupById(w http.ResponseWriter, r *http.ResponseWriter) {
// 	id := chi.URLParam(r, "id")
// 	group, err := app.models.Group.GetRepo
// }
