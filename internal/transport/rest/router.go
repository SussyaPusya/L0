package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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

	r.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // или []string{"*"} для всех
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Максимальное время кэширования preflight-запроса (в секундах)
	}))

	r.router.Use(middleware.URLFormat)
	r.router.Get("/orders/{order_id}", r.handlers.GetOrder)

	r.router.Get("/", r.handlers.Main)

	http.ListenAndServe(":8081", r.router)
}
