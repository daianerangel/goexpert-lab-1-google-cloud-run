package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type LocationInfo struct {
	Localidade string `json:"localidade"`
}

type WeatherInfo struct {
    Current struct {
        Temperature            float64 `json:"temp_c"`
    } `json:"current"`
}

func getLocation(zipCode string) (string, error) {
	tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{Transport: tr}
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", zipCode)
	resp, err := client.Get(url)
	fmt.Println("StatusCode:", resp)
	fmt.Println("err:", err)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var location LocationInfo
	if err := json.NewDecoder(resp.Body).Decode(&location); err != nil {
		return "", err
	}

	return location.Localidade, nil
}

func getWeather(city string) (WeatherInfo, error) {
	tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{Transport: tr}
	encodedCity := url.QueryEscape(city)
	completeUrl := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=6c0e6aefacc44ed0a69130616242705&q=%s", encodedCity)
	resp, err := client.Get(completeUrl)
	fmt.Println("StatusCode:", resp)
	fmt.Println("err:", err)
	if err != nil {
		return WeatherInfo{}, err
	}
	defer resp.Body.Close()

	var weather WeatherInfo
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return WeatherInfo{}, err
	}

	return weather, nil
}

func temperatureHandler(w http.ResponseWriter, r *http.Request) {
	zipCode := r.URL.Query().Get("zipcode")
	if len(zipCode) != 8 {
		http.Error(w, "invalid zipcode", 422)
		return
	}

	city, err := getLocation(zipCode)
	if err != nil || city == "" {
		http.Error(w, "can not find zipcode", 404)
		return
	}

	weather, err := getWeather(city)
	if err != nil {
		http.Error(w, "failed to get weather info", 500)
		return
	}

	tempC := weather.Current.Temperature
	tempF := tempC*1.8 + 32
	tempK := tempC + 273

	response := map[string]float64{
		"temp_C": tempC,
		"temp_F": tempF,
		"temp_K": tempK,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}



func main() {
	http.HandleFunc("/temperature", temperatureHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
