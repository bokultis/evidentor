package main

import (
	"evidentor/api/db"
	"evidentor/api/logger"
	"evidentor/api/router"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

//Logger ...
var Logger *logrus.Logger

func main() {

	logger.Logger = logger.NewLogger()

	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		logger.Logger.Print("No .env file found")
	}

	//append user routes

	port := os.Getenv("PORT")
	router := router.NewRouter()

	headers := handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"})
	methods := handlers.AllowedMethods([]string{"POST", "GET", "OPTIONS", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})

	//Setup database
	db.DB = db.SetupDB()
	defer db.DB.Close()

	logger.Logger.Printf("Server starts at localhost: %s", port)
	//create http server
	logger.Logger.Fatal(http.ListenAndServe(":"+port, handlers.CORS(headers, methods, origins)(router)))
}
