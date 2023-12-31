package server

import (
	"context"
	"errors"
	"log"
	"net/http"

	repository "github.com/F-Dupraz/go-rest-api.git/repository"
	database "github.com/F-Dupraz/go-rest-api.git/database"
	"github.com/F-Dupraz/go-rest-api.git/websocket"

	"github.com/rs/cors"
	"github.com/gorilla/mux"
)

type Config struct {
	Port        string
	JWTSecret   string
	DatabaseURL string
}

type Server interface {
	Config() *Config
	Hub() *websocket.Hub
}

type Broker struct {
	config *Config
	router *mux.Router
	hub *websocket.Hub
}

func (b *Broker) Config() *Config {
	return b.config
}

func (b *Broker) Hub() *websocket.Hub {
	return b.hub
}

func NewServer(ctx context.Context, config *Config) (*Broker, error) {
	if config.Port == "" {
		return nil, errors.New("Port is required!")
	}
	if config.JWTSecret == "" {
		return nil, errors.New("JWT secret is required!")
	}
	if config.DatabaseURL == "" {
		return nil, errors.New("Database url is required!")
	}
	broker := &Broker{
		config: config,
		router: mux.NewRouter(),
		hub:    websocket.NewHub(),
	}
	return broker, nil
}

func (b *Broker) Start(binder func(s Server, r *mux.Router)) {
	b.router = mux.NewRouter()
	binder(b, b.router)
	handler := cors.Default().Handler(b.router)
	repo, err := database.NewPostgresRepository(b.config.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	go b.hub.Run()
	repository.SetRepository(repo)
	log.Println("Starting server on port", b.config.Port)
	if err := http.ListenAndServe(b.config.Port, handler); err != nil {
		log.Println("Error starting server:", err)
	} else {
		log.Fatalf("Server stopped!")
	}
}