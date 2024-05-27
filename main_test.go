package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)


func TestTemperatureHandlerSuccess(t *testing.T) {

	req, err := http.NewRequest("GET", "/temperature?zipcode=10001", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(temperatureHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"temp_C":28.5,"temp_F":83.3,"temp_K":301.65}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
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
