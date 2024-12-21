package main

import (
	"log"

	"github.com/babaYaga451/social/internal/db"
	"github.com/babaYaga451/social/internal/env"
	"github.com/babaYaga451/social/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://postgres:password@localhost/postgres?sslmode=disable")

	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	store := store.NewStorage(conn)

	db.Seed(store, conn)
}
