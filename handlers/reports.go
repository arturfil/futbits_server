package handlers

import (
	"chi_soccer/helpers"
	"chi_soccer/services"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

var report services.Report

// GET/reports/report
func GetAllReports(w http.ResponseWriter, r *http.Request) {
	var reports services.Report
	all, err := reports.GetAllReports()
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"reports": all})
}

// GET/reports/report/id
func GetReportById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	report, err := report.GetReportById(id)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, report)
}

// GET/reports/game_id
func GetReportsOfGame(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "game_id")

	reports, err := report.GetAllReportsByGameId(id)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, reports)
}

// GET/reports/group_id
func GetReportsOfGroup(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "group_id")
	reports, err := report.GetAllReportsByGroupId(id)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, reports)
}

// GET/reports/user_id
func GetReportsOfUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "user_id")
	reports, err := report.GetAllReporstByUserId(id)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, reports)
}

// POST/reports/upload
func UploadReportCSV(w http.ResponseWriter, r *http.Request) {
	var rp services.Report
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("reports")
	if err != nil {
		fmt.Println("Error Retrieving File")
		fmt.Println(err)
		return
	}
	defer file.Close()

	fmt.Printf("Uploaded file: %+v\n", handler.Filename)

	reader := csv.NewReader(file)
	rp.UploadReport(reader)

	msg := "CSV file processed successfully"
	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"msg": msg})
}

// POST/reports/report
func CreateReport(w http.ResponseWriter, r *http.Request) {
	var rp services.Report
	err := json.NewDecoder(r.Body).Decode(&rp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newReport, err := report.CreateReport(rp)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
	}
	helpers.WriteJSON(w, http.StatusOK, newReport)
}

// func GetGroupById(w http.ResponseWriter, r *http.ResponseWriter) {
// 	id := chi.URLParam(r, "id")
// 	group, err := mod.Group.GetRepo
// }
