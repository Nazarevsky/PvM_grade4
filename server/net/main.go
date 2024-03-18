package net

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"net/http"
	"server/net/handlers"
)

func Run(port uint) {
	r := chi.NewRouter()

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           300,
	})

	r.Use(corsMiddleware.Handler)

	// Роутер ендпоїнтів
	r.Route("/", func(r chi.Router) {
		// Ендпоїнт на заміну круглих дужок на квадратні
		r.Post("/replace_round_brackets", handlers.ReplaceRoundBrackets)
	})

	// Слухач підключень
	startHttpListener(port, r)
}

func startHttpListener(port uint, r *chi.Mux) {
	fmt.Printf("Server is listening on port %d...\n", port)

	// Метод, що запускає сервер на вказаному порту
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}
