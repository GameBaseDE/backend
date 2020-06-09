package main

import (
	"github.com/gin-contrib/cors"
	"gitlab.tandashi.de/GameBase/gamebase-backend/openapi"
	"os"
)

func main() {
	router := openapi.NewRouter()
	router.Use(cors.Default())

	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	if router.Run(":"+port) != nil {
		println("Could not start the server")
	}
}
