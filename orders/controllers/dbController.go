package controllers

import (
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/orders/microservices/db"
)

var DB *pgxpool.Pool

func init() {
	db, err := db.ConnectToDB()
	if err != nil {
		log.Fatal("Error connecting to the database:", err)

	}

	DB = db
}
