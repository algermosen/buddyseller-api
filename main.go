package main

import (
	"context"
	"example/buddyseller-api/db"
	"example/buddyseller-api/db/data_store"
	"example/buddyseller-api/routes"
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
		log.Fatalf("Error loading environment: \n\t%v", err)
	}

	ctx := context.Background()
	DB, err := db.InitDB(ctx)

	if err != nil {
		log.Fatalf("Error initializing database: \n\t%v", err)
	}

	defer DB.Close(ctx)

	ds := data_store.New(DB)

	r := gin.Default()
	routes.RegisterRoutes(r)

	// Start HTTP server
	srv := &http.Server{
		Addr:    ":9000",
		Handler: r.Handler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe failed: \n\t%v", err)
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
	default:
		gin.SetMode(gin.DebugMode)
		err = godotenv.Load(".env.dev")
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
		log.Fatalf("Server shutdown failed: \n\t%v", err)
	}
	log.Println("Server gracefully stopped")
}
