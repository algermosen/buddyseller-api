package main

import (
	"example/buddyseller-api/database"
	"example/buddyseller-api/routes"
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	database.InitDB()
	server := gin.Default()
	routes.RegisterRoutes(server)
	server.Run(":9000")
}

func init() {
	var mode string
	flag.StringVar(&mode, "mode", "", "Provide environment to be executed. 'release' | 'debug' | 'test'")

	flag.Parse()

	switch mode {
	case "release":
		gin.SetMode(gin.ReleaseMode)
		godotenv.Load(".env")
	case "test":
		gin.SetMode(gin.TestMode)
		godotenv.Load(".env.test")
	case "debug":
	default:
		gin.SetMode(gin.DebugMode)
		godotenv.Load(".env.dev")
	}
}
