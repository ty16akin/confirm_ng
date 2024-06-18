package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
	"github.com/ty16akin/ConfirmNG/internal/database"
	"github.com/ty16akin/ConfirmNG/internal/handler"
)

func main() {

	enverr := godotenv.Load()
	if enverr != nil {
		log.Fatal("Error loading .env file")
	}
	var dbconnection = os.Getenv("MONGO_URI")
	err := database.Init(dbconnection, "confirm_ng")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Connected to MongoDB!")

	defer func() {
		err := database.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	r := chi.NewRouter()
	server := &http.Server{
		Addr:    ":3001",
		Handler: r,
	}
	r.Use(middleware.Logger)

	r.Route("/users", loadUserRoutes)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run the server
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
		return
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}

func loadUserRoutes(router chi.Router) {
	userHandler := &handler.User{}
	fuelStationHandler := &handler.FS{}

	router.Post("/", userHandler.CreateUser)
	router.Get("/", userHandler.GetUsers)
	router.Get("/{id}", userHandler.GetUserById)
	router.Patch("/{id}", userHandler.UpdateUserById)
	router.Delete("/{id}", userHandler.DeleteUserById)
	router.Get("/fuel-stations", fuelStationHandler.SearchFuelStations)
}
