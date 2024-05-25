package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mdalboni/goexpert_3/internals/services"
)

type GetWeatherResponse struct {
	TempC float64 `json:"temp_c"`
	TempF float64 `json:"temp_f"`
	TempK float64 `json:"temp_k"`
}

type WeatherHandler struct {
	CEPService     services.CEPService
	WeatherService services.WeatherService
}

func NewWeatherHandler() *WeatherHandler {
	return &WeatherHandler{
		CEPService:     services.NewViaCEPService(),
		WeatherService: services.NewWeatherAPIService(),
	}
}

// GetWeather returns the weather
func (wh *WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	zipCode := chi.URLParam(r, "zipCode")
	if len(zipCode) != 8 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("invalid zipcode"))
		return
	}
	responseCEP, error := wh.CEPService.GetAddressByCEP(zipCode)
	if error != nil {
		switch error {
		case services.ErrCEPNotFound:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("can not find zipcode"))
			return
		case services.ErrInvalidCEP:
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("invalid zipcode"))
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
			return
		}
	}
	responseWeather, error := wh.WeatherService.GetWeatherByCity(responseCEP.Localidade)
	if error != nil {
		switch error {
		case services.ErrCEPNotFound:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("can not find zipcode"))
			return
		case services.ErrInvalidCEP:
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("invalid zipcode"))
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal server error"))
			return
		}
	}
	// Truncating values to ensure only 1 decimal place
	output := GetWeatherResponse{
		TempC: float64(int(responseWeather.Current.TempC*10)) / 10,
		TempF: float64(int(responseWeather.Current.TempF*10)) / 10,
		TempK: float64(int((responseWeather.Current.TempC+273.15)*10)) / 10,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}
