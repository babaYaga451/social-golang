package main

import (
	"log"

	"github.com/babaYaga451/social/internal/db"
	"github.com/babaYaga451/social/internal/env"
	"github.com/babaYaga451/social/internal/store"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

//	@title			Go-Social
//	@description	API for social platform to follow users and post content

//	@BasePath					/v1
//	@securityDefinitions.apiKey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config{
		addr:   env.GetString("ADDR", ":8080"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgresql://postgres:password@localhost/postgres?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Database
	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)

	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("Database connection pool established")

	store := store.NewStorage(db)

	app := &application{
		conf:   cfg,
		store:  store,
		logger: logger,
	}

	mux := app.mount()
	logger.Fatal(app.run(mux))
}
