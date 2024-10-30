package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"okx-bot/frontend-service/controllers"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("Hello test")
	r := mux.NewRouter()
	r.Headers("user", "1")
	r.HandleFunc("/api/signal/okx/all", controllers.GetAllOkxSignals).Methods("GET")
	ts := httptest.NewServer(r)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/api/signal/okx/all")
	if err != nil {
		fmt.Errorf("Expected nil, received %s", err.Error())
	}
	if res.StatusCode != http.StatusOK {
		fmt.Errorf("Expected %d, received %d", http.StatusOK, res.StatusCode)
	}
}
