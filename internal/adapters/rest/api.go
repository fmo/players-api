package rest

import (
	"context"
	"github.com/fmo/players-api/internal/api"
	"github.com/fmo/players-api/internal/ports"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Adapter struct {
	api  ports.APIPorts
	port string
}

func NewAdapter(api ports.APIPorts, port string) *Adapter {
	return &Adapter{
		api:  api,
		port: port,
	}
}

func (a Adapter) Run(ctx context.Context) {
	r := chi.NewMux()

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	r.Use(corsHandler.Handler)

	h := api.HandlerFromMux(a, r)

	srv := &http.Server{
		Addr:    a.port,
		Handler: h,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
