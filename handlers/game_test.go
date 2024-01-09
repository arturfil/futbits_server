package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type testGame struct {
	name           string
	method         string
	json           string
	paramID        string
	handler        http.HandlerFunc
	expectedStatus int
}


func Test_app_game(t *testing.T) {

	//var tests = []testGame{
	//	{"/api/v1/games/35069b00-a556-4b4b-acb2-d842007b8ffa", "GET", "", "", GetAllGames, http.StatusOK},
	//}

    req, err := http.NewRequest("GET", "/api/v1/games/35069b00-a556-4b4b-acb2-d842007b8ffa", nil)
    if err != nil {
        t.Errorf("Error creating a new request: %v",err) 
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(GetAllGames)
    handler.ServeHTTP(rr, req)

    if rr.Code != http.StatusOK {
        t.Errorf("%d: wrong sttaus code returned", rr.Code)
    }

}
