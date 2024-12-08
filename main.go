package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type WeatherData struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temperature"`
	Condition   string  `json:"condition"`
}

type OpenWeatherResponse struct {
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

func fetchWeatherData(city string) WeatherData {
	apiKey := "" 
		apiURL := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=%s", city, apiKey)
	
		client := http.Client{Timeout: 10 * time.Second}
		resp, err := client.Get(apiURL)
		if err != nil {
			fmt.Printf("Error fetching data for %s: %v\n", city, err)
			return WeatherData{City: city, Temperature: 0, Condition: "Unknown"}
		}
		defer resp.Body.Close()
	
	
		fmt.Printf("Response Status for %s: %s\n", city, resp.Status)
	
		var openWeatherResponse OpenWeatherResponse
		if err := json.NewDecoder(resp.Body).Decode(&openWeatherResponse); err != nil {
			fmt.Printf("Error decoding response for %s: %v\n", city, err)
			return WeatherData{City: city, Temperature: 0, Condition: "Unknown"}
		}

		if len(openWeatherResponse.Weather) == 0 {
			fmt.Printf("No weather data for %s\n", city)
			return WeatherData{City: city, Temperature: 0, Condition: "Unknown"}
		}
	
		fmt.Printf("Weather Data for %s: %+v\n", city, openWeatherResponse)
	
		return WeatherData{
			City:       city,
			Temperature: openWeatherResponse.Main.Temp,
			Condition:   openWeatherResponse.Weather[0].Description,
		}
	}
	

	func weatherHandler(w http.ResponseWriter, r *http.Request) {
		cities := strings.Split(r.URL.Query().Get("cities"), ",")
		for i, city := range cities {
			cities[i] = strings.TrimSpace(city)
		}
	
		workerCount := 4
		workerPool := NewWorkerPool(workerCount)
		
		workerPool.Start()
	
		for _, city := range cities {
			task := Task{City: city}
			workerPool.AddTask(task)
		}
	
		go workerPool.Stop()
	
		resultMap := make(map[string]WeatherData)
	
		for result := range workerPool.Results() {
			resultMap[result.Weather.City] = result.Weather
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resultMap); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}	

func main() {
	http.HandleFunc("/weather", weatherHandler)
	port := ":8080"
	fmt.Println("Server is running on http://localhost" + port)
	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}




