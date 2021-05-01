package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/justhackit/go-webapp-template/datastore/animal"
	handlerAnimal "github.com/justhackit/go-webapp-template/delivery/animal"
	"github.com/justhackit/go-webapp-template/driver"
)

func main() {
	os.Setenv("SQL_HOST", "localhost")
	os.Setenv("SQL_USER", "root")
	os.Setenv("SQL_PASSWORD", "Secret_123")
	os.Setenv("SQL_PORT", "3306")
	os.Setenv("SQL_DB", "animal")
	// get the mysql configs from env:
	conf := driver.MySQLConfig{
		Host:     os.Getenv("SQL_HOST"),
		User:     os.Getenv("SQL_USER"),
		Password: os.Getenv("SQL_PASSWORD"),
		Port:     os.Getenv("SQL_PORT"),
		Db:       os.Getenv("SQL_DB"),
	}
	var err error

	db, err := driver.ConnectToMySQL(conf) //Create Database connection
	if err != nil {
		log.Println("could not connect to sql, err:", err)
		return
	}

	datastore := animal.New(db)             //Instantiate datastore's Animal interface's implementation
	handler := handlerAnimal.New(datastore) //DB Injeceting

	http.HandleFunc("/animal", handler.Handler)
	fmt.Println(http.ListenAndServe(":9000", nil))
}
