package main

import (
	"context"
	"example/buddyseller-api/api"
	"example/buddyseller-api/api/handlers"
	"example/buddyseller-api/db"
	"example/buddyseller-api/db/datastore"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := loadEnv()

	if err != nil {
		log.Fatalf("Error loading environment: \n%v", err)
	}

	ctx := context.Background()
	conn, err := db.InitDB(ctx)

	if err != nil {
		log.Fatalf("Error initializing database: \n%v", err)
	}

	defer conn.Close(ctx)

	ds := datastore.New(conn)

	userHandler := handlers.PostgresUserHandler{DS: ds}
	orderHandler := handlers.PostgresOrderHandler{DS: ds}
	productHandler := handlers.PostgresProductHandler{DS: ds}
	sessionHandler := handlers.PostgresSessionHandler{DS: ds}

	handlers := api.RouterHandlers{
		UserHandler:    &userHandler,
		OrderHandler:   &orderHandler,
		ProductHandler: &productHandler,
		SessionHandler: &sessionHandler,
	}

	r := api.RouterSetup(&handlers)

	// Start HTTP server
	srv := &http.Server{
		Addr:    ":9000",
		Handler: r.Handler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe failed: \n%v", err)
		}
	}()

	// Wait for termination signal
	waitForShutdown(srv)
}

func loadEnv() error {
	var mode string
	flag.StringVar(&mode, "mode", "", "Provide environment to be executed. 'release' | 'debug' | 'test'")

	flag.Parse()

	var err error
	switch mode {
	case "release":
		gin.SetMode(gin.ReleaseMode)
		err = godotenv.Load(".env")
	case "test":
	case "debug":
		gin.SetMode(gin.DebugMode)
		err = godotenv.Load(".env.dev")
	default:
		panic("Must provide environment to be executed. 'release' | 'debug' | 'test'")
	}

	return err
}

func waitForShutdown(server *http.Server) {
	// Create a channel to receive OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	log.Println("Shutting down server...")
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: \n%v", err)
	}
	log.Println("Server gracefully stopped")
}
