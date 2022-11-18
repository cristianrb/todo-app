package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (u *handlerConfig) Routes() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Heartbeat("/ping"))
	router.Get("/items", u.getTodoItems)
	return router
}
