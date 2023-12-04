package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Nebezdar/tech_lvl0.git/internal/app"
	"github.com/Nebezdar/tech_lvl0.git/internal/cache"
	"github.com/Nebezdar/tech_lvl0.git/internal/db"
	"github.com/Nebezdar/tech_lvl0.git/internal/http"
	"github.com/Nebezdar/tech_lvl0.git/internal/messaging"
)

func main() {
	// Connect to PostgreSQL
	db, err := db.NewDatabase("postgres://your_username:your_password@localhost/your_database?sslmode=disable")
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}
	defer db.Close()

	// Connect to Nats-streaming
	nats, err := messaging.NewNats("nats://localhost:4222")
	if err != nil {
		log.Fatal("Error connecting to Nats:", err)
	}
	defer nats.Close()

	// Create Cache
	cache := cache.NewCache(3600)

	// Create OrderService
	orderService := app.NewOrderService(db, cache, nats)
	orderService.Start()
	defer orderService.Stop()

	// Create HTTP Server
	httpServer := http.NewServer(orderService)
	go httpServer.Start("8080")

	// Handle shutdown gracefully
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Wait for a signal to shutdown
	<-stop

	log.Println("Shutting down gracefully...")
}
