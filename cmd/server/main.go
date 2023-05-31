package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"reby/api"
	"reby/api/handlers"
	"reby/app/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	conf := config.Get()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(api.JSONResponseMiddleware)
	r.Use(api.RecovererMiddleware)
	h := handlers.InitHandlers(conf)

	handlers.AddRideEndpoints(r, h.Ride)

	server := &http.Server{
		ReadTimeout:       3 * time.Second,
		WriteTimeout:      6 * time.Second,
		Addr:              fmt.Sprintf("%s:%s", conf.APIURL, conf.APIPort),
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           r,
	}

	if err := server.ListenAndServe(); err != http.ErrServerClosed && err != nil {
		log.Fatalf("Error starting http server <%s>", err)
	}
}
