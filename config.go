package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var accessLogFile *os.File
var db *gorm.DB

func InitDB() {
	config := map[string]string{
		"DB_Username": os.Getenv("DB_USERNAME"),
		"DB_Password": os.Getenv("DB_PASSWORD"),
		"DB_Port":     os.Getenv("DB_PORT"),
		"DB_Host":     os.Getenv("DB_HOST"),
		"DB_Name":     os.Getenv("DB_NAME"),
	}
	connectionString :=
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
			config["DB_Username"],
			config["DB_Password"],
			config["DB_Host"],
			config["DB_Port"],
			config["DB_Name"])
	var e error
	// Create a custom GORM logger
	var customLogger = logger.New(
		log.New(io.MultiWriter(os.Stdout, accessLogFile), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             250 * time.Millisecond,
			IgnoreRecordNotFoundError: false,
			LogLevel:                  logger.Info,
			Colorful:                  true,
		},
	)
	// Create the database connection with the custom logger
	db, e = gorm.Open(mysql.Open(connectionString), &gorm.Config{
		PrepareStmt: true,
		Logger:      customLogger,
	})
	if e != nil {
		// Handle database connection error
		log.Fatalf("Error initializing the database!")
	}
	// InitMigrate()
}
