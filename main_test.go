package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)


func TestTemperatureHandlerSuccess(t *testing.T) {

	req, err := http.NewRequest("GET", "/temperature?zipcode=23093010", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(temperatureHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]float64
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	expectedResponse := map[string]float64{
		"temp_C": 22,
		"temp_F": 71.6,
		"temp_K": 295,
	}

	if !reflect.DeepEqual(response, expectedResponse){
		t.Errorf("handler returned unexpected body: got %v want %v", response, expectedResponse)
	}
}

func TestTemperatureHandlerInvalidZipCode(t *testing.T) {
	req, err := http.NewRequest("GET", "/temperature?zipcode=123", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(temperatureHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnprocessableEntity)
	}
}

func TestZipCodeNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "/temperature?zipcode=00000000", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(temperatureHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnprocessableEntity)
	}
}
