package main

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/mdalboni/goexpert_3/internals/handlers"
)

func main() {
	slog.Info("Starting server on port 8080")
	runServer()
}

func runServer() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	weatherHandler := handlers.NewWeatherHandler()
	r.With(addContext).Get("/weather/{zipCode}", weatherHandler.GetWeather)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		slog.Error("Error starting server: ", err)
	}
}

type contextKey string

func addContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), contextKey("request_id"), uuid.New().String())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
