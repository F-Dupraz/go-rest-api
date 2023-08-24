package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	
	"github.com/F-Dupraz/go-rest-api.git/handlers"
	"github.com/F-Dupraz/go-rest-api.git/middlewares"
	"github.com/F-Dupraz/go-rest-api.git/server"
	"github.com/F-Dupraz/go-rest-api.git/websocket"
)

func BindRoutes(s server.Server, r *mux.Router) {
	hub := websocket.NewHub()

	r.Use(middlewares.CheckAuthMiddleware(s))
	
	r.HandleFunc("/", handlers.HomeHandlers(s)).Methods(http.MethodGet)
	r.HandleFunc("/signup", handlers.SignUpHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/login", handlers.LoginHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/me", handlers.MeHandler(s)).Methods(http.MethodGet)

	r.HandleFunc("/posts", handlers.InsertPostHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/posts", handlers.InsertPostHandler(s)).Methods(http.MethodPost)
	r.HandleFunc("/posts/{id}", handlers.GetPostById(s)).Methods(http.MethodGet)
	r.HandleFunc("/posts/{id}", handlers.UpdatePost(s)).Methods(http.MethodPut)
	r.HandleFunc("/posts", handlers.ListPost(s)).Methods(http.MethodGet)

	go hub.Run()
	r.HandleFunc("/ws", hub.HandleWebSocket)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file %v.\n", err)
	}

	PORT := os.Getenv("PORT")
	JWT_SECRET := os.Getenv("JWT_SECRET")
	DATABASE_URL := os.Getenv("DATABASE_URL")

	s, err := server.NewServer(context.Background(), &server.Config {
		Port: PORT,
		JWTSecret: JWT_SECRET,
		DatabaseURL: DATABASE_URL,
	})

	if err != nil {
		log.Fatalf("Error creating server %v.\n", err)
	}

	s.Start(BindRoutes)
}
