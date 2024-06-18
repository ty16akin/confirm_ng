package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

const googlePlacesAPIURL = "https://maps.googleapis.com/maps/api/place/nearbysearch/json"

type FS struct{}

type Place struct {
	Name    string `json:"name"`
	Address string `json:"vicinity"`
}

type PlacesResponse struct {
	Results []Place `json:"results"`
}

func (f *FS) SearchFuelStations(w http.ResponseWriter, r *http.Request) {

	location := r.URL.Query().Get("location")
	if location == "" {
		http.Error(w, "Location is required", http.StatusBadRequest)
		return
	}

	enverr := godotenv.Load()
	if enverr != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		http.Error(w, "API key is missing", http.StatusInternalServerError)
		return
	}

	params := url.Values{}
	params.Add("location", location)
	params.Add("radius", "5000") // 5 km radius
	params.Add("type", "gas_station")
	params.Add("key", apiKey)

	apiURL := fmt.Sprintf("%s?%s", googlePlacesAPIURL, params.Encode())

	resp, err := http.Get(apiURL)
	if err != nil {
		log.Printf("Failed to make API request: %v", err)
		http.Error(w, "Failed to retrieve data", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Unexpected status code: %v", resp.StatusCode)
		http.Error(w, "Failed to retrieve data", http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		http.Error(w, "Failed to retrieve data", http.StatusInternalServerError)
		return
	}

	var placesResp PlacesResponse
	err = json.Unmarshal(body, &placesResp)
	if err != nil {
		fmt.Println("Failed to marshal:", err)
		http.Error(w, "Failed to parse data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(placesResp.Results)
}
