package handlers

import (
	"github.com/go-chi/chi/v5"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func (u handlerConfig) Routes() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Heartbeat("/ping"))
	router.Post("/items", u.pushItemToQueue)
	return router
}
