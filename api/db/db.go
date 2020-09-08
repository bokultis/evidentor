package db

import (
	"errors"
	"log"
	"time"

	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//DB database global
var DB *gorm.DB

func SetupDB() *gorm.DB {

	//db config vars
	var dbHost string = os.Getenv("DB_HOST")
	var dbPort string = os.Getenv("DB_PORT")
	var dbName string = os.Getenv("DB_NAME")
	var dbUser string = os.Getenv("DB_USERNAME")
	var dbPassword string = os.Getenv("DB_PASSWORD")

	//setup loger
	logger.Default.LogMode(logger.Error)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Millisecond, // Slow SQL threshold
			LogLevel:      logger.Error,     // Log level
			Colorful:      true,             // Enable color
		},
	)

	//connect to db
	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8&parseTime=True&loc=Local"

	db, dbError := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if dbError != nil {
		panic(dbError)
	}

	return db
}

// GetDbName get database name from env
func GetDbName() (string, error) {

	database := os.Getenv("DB_NAME")
	if database == "" {
		return "", errors.New("Database not set")
	}

	return database, nil
}
