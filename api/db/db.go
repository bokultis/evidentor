package db

import (
	"errors"
	"evidentor/api/logger"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//database global
var DB *gorm.DB

type GormLogger struct{}

func (*GormLogger) Print(v ...interface{}) {

	log := logger.NewLogger()
	log.SetReportCaller(false)
	if v[0] == "sql" {
		log.Printf("%v %v  %v", v[1], v[2], v[3])
	}
	if v[0] == "log" {
		log.Print(v[2])
	}
}

func SetupDB() *gorm.DB {

	//db config vars
	var dbHost string = os.Getenv("DB_HOST")
	var dbPort string = os.Getenv("DB_PORT")
	var dbName string = os.Getenv("DB_NAME")
	var dbUser string = os.Getenv("DB_USERNAME")
	var dbPassword string = os.Getenv("DB_PASSWORD")

	////fmt.Printf("dbName: %s, dbHost: %s", dbName, dbHost)
	//connect to db
	db, dbError := gorm.Open("mysql", dbUser+":"+dbPassword+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?charset=utf8&parseTime=True&loc=Local")
	//db, dbError := gorm.Open("mysql", dbUser+":"+dbPassword+"@/"+dbName+"?charset=utf8&parseTime=True&loc=Local")
	if dbError != nil {
		panic(dbError)
	}

	db.SetLogger(&GormLogger{})
	db.DB().SetMaxIdleConns(0)

	return db
}

// get database name from env
func GetDbName() (string, error) {

	database := os.Getenv("DB_NAME")
	if database == "" {
		return "", errors.New("Database not set")
	}

	return database, nil
}
