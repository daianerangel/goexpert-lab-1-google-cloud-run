package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

	if tempC, ok := response["temp_C"]; !ok || tempC == 0 {
		t.Errorf("temp_C is missing or zero: got %v", tempC)
	}
	if tempF, ok := response["temp_F"]; !ok || tempF == 0 {
		t.Errorf("temp_F is missing or zero: got %v", tempF)
	}
	if tempK, ok := response["temp_K"]; !ok || tempK == 0 {
		t.Errorf("temp_K is missing or zero: got %v", tempK)
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
