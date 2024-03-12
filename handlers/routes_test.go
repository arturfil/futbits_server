package handlers

import (
	"net/http"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

type registeredRoutes struct {
    route string
    method string
}

func Test_app_routes(t *testing.T) {
    var registered = []registeredRoutes {
        {"/api/v1/games/{user_id}", "GET"},
        {"/api/v1/groups/", "GET"},
        {"/api/v1/reports/group/{group_id}", "GET"},
    }

    mux := Routes() 

    chiRoutes := mux.(chi.Routes)

    for _, route := range registered {
        if !routeExists(route.route, route.method, chiRoutes) {
            t.Errorf("route %s is not registered", route.route)
        }
    }
}

func routeExists(testRoute, testMethod string, chiRoutes chi.Routes) bool {
    found := false
    _ = chi.Walk(
        chiRoutes, 
        func(method string, route string, handler http.Handler, middleware ...func(http.Handler) http.Handler) error {
            if strings.EqualFold(method, testMethod) && strings.EqualFold(route, testRoute,) {
                found = true
            }
            return nil 
        })
        return found
}


