package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"transcoder/internal/converter"

	_ "github.com/lib/pq"
)

func connectPostgres() (*sql.DB, error) {
	host := getEnvOrDefault("POSTGRES_HOST", "host.docker.internal")
	port := getEnvOrDefault("POSTGRES_PORT", "5431")
	user := getEnvOrDefault("POSTGRES_USER", "root")
	password := getEnvOrDefault("POSTGRES_PASSWORD", "root")
	dbname := getEnvOrDefault("POSTGRES_DBNAME", "codetube_transcoder")
	sslmode := getEnvOrDefault("POSTGRES_SSLMODE", "disable")

	connSrt := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)

	db, err := sql.Open("postgres", connSrt)
	if err != nil {
		slog.Error("Error connecting to database", slog.String("error", err.Error()))
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		slog.Error("Error pinging database", slog.String("error", err.Error()))
		return nil, err
	}

	slog.Info("Connected to Postgres successfully")
	return db, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func main() {
	db, err := connectPostgres()
	if err != nil {
		panic(err)
	}
	converter := converter.NewVideoConvert(db)
	converter.Handle([]byte(`{"video_id": 6, "path": "/media/uploads/6"}`))
}
