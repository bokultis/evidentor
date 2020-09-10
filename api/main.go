package main

import (
	"net/http"
	"os"

	"github.com/bokultis/evidentor/api/db"
	"github.com/bokultis/evidentor/api/logger"
	"github.com/bokultis/evidentor/api/redis"
	"github.com/bokultis/evidentor/api/router"

	"github.com/gorilla/handlers"
	"github.com/sirupsen/logrus"
)

//Logger ...
var Logger *logrus.Logger

func main() {
	logger.Logger = logger.NewLogger()

	//append user routes
	port := os.Getenv("PORT")
	router := router.NewRouter()

	headers := handlers.AllowedHeaders([]string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"})
	methods := handlers.AllowedMethods([]string{"POST", "GET", "OPTIONS", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})

	//Setup database
	db.DB = db.SetupDB()

	//Setup redis
	redis.RedisClient = redis.SetupRedis()

	logger.Logger.Printf("Server starts at localhost: %s", port)
	//create http server
	logger.Logger.Fatal(http.ListenAndServe(":"+port, handlers.CORS(headers, methods, origins)(router)))
}
