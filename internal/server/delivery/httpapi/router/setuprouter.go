package router

import "github.com/go-chi/chi/v5"

func (rout *Router) SetupRouter() {
	rout.Mux = chi.NewRouter()

	// Группа без проверки JWT
	rout.Mux.Group(func(public chi.Router) {
		public.Post("/api/create", rout.handler.Auth.PostAuthToken) // Получение токенов
	})

	// все остальные запросы
	rout.Mux.Group(func(protected chi.Router) {
		protected.Use(rout.middleware.JwtCheck)
		protected.Put("/api/refresh", rout.handler.Auth.PutRefresh) // Refresh токенов
	})
}
