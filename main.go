package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"ticket-system/config"
	"ticket-system/database"
	"ticket-system/handlers"
	authmiddleware "ticket-system/middleware"
)

func SetupRouter(cfg *config.Config, db *database.DB) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	// Enable CORS for web clients and automated testing suites
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	authHandler := handlers.NewAuthHandler(db, cfg.JWTSecret)
	ticketHandler := handlers.NewTicketHandler(db)

	r.Get("/health", handlers.HealthCheck)

	r.Post("/auth/register", authHandler.Register)
	r.Post("/auth/login", authHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(authmiddleware.JWTAuth(cfg.JWTSecret))

		r.Post("/tickets", ticketHandler.CreateTicket)
		r.Get("/tickets", ticketHandler.ListTickets)
		r.Get("/tickets/{id}", ticketHandler.GetTicketByID)
		r.Patch("/tickets/{id}/status", ticketHandler.UpdateTicketStatus)
	})

	return r
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.ConnectDB(cfg.MongoURI, cfg.DBName)
	if err != nil {
		log.Printf("warning: mongo connection failed: %v", err)
	}

	router := SetupRouter(cfg, db)

	log.Printf("Starting Ticket System API server on port %s...", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
