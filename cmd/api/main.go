package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Nel-sokchhunly/go-rest-api/cmd/api/router"
	"github.com/Nel-sokchhunly/go-rest-api/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

const fmtDBString = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable"

// @title Go REST API
// @version 1.0
// @description This is a sample server Go REST API server with Folder Structure.

// @license.name MIT License
// @license.url http://opensource.org/licenses/MIT

// @host localhost:8080
// @router /api/v1 [get]
func main() {
	c := config.New()

	var logLevel gormlogger.LogLevel
	if c.DB.Debug {
		logLevel = gormlogger.Info
	} else {
		logLevel = gormlogger.Error
	}

	dbString := fmt.Sprintf(fmtDBString, c.DB.Host, c.DB.User, c.DB.Password, c.DB.DBName, c.DB.Port)
	db, err := gorm.Open(postgres.Open(dbString), &gorm.Config{Logger: gormlogger.Default.LogMode(logLevel)})
	if err != nil {
		log.Fatal("DB connection start failure")
		return
	}

	router := router.New(db)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", c.Server.Port),
		Handler:      router,
		ReadTimeout:  c.Server.TimeoutRead,
		WriteTimeout: c.Server.TimeoutWrite,
		IdleTimeout:  c.Server.TimeoutIdle,
	}

	log.Printf("Server started on port %d", c.Server.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server failed to start: ", err)
	}
}
