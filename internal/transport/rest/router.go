package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	router   *chi.Mux
	handlers *Handlers
}

func NewRouter(h *Handlers) *Router {

	r := chi.NewRouter()
	return &Router{
		router:   r,
		handlers: h,
	}
}

func (r *Router) Run() {

	r.router.Use(middleware.URLFormat)
	r.router.Get("/orders/{order_id}", r.handlers.GetOrder)

	http.ListenAndServe(":8081", r.router)
}
