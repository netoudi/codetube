package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"transcoder/internal/converter"
	"transcoder/internal/rabbitmq"

	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
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

	rabbitmqUrl := getEnvOrDefault("RABBITMQ_URL", "amqp://guest:guest@host.docker.internal:5672")
	rabbitmqClient, err := rabbitmq.NewRabbitClient(rabbitmqUrl)
	if err != nil {
		panic(err)
	}
	defer rabbitmqClient.Close()

	conversionExchange := getEnvOrDefault("CONVERSION_EXCHANGE", "conversion_exchange")
	conversionQueue := getEnvOrDefault("CONVERSION_QUEUE", "conversion_queue")
	conversionKey := getEnvOrDefault("CONVERSION_KEY", "conversion_key")
	confirmationKey := getEnvOrDefault("CONFIRMATION_KEY", "confirmation_key")
	confirmationQueue := getEnvOrDefault("CONFIRMATION_QUEUE", "confirmation_queue")

	converter := converter.NewVideoConvert(db, rabbitmqClient)

	msgs, err := rabbitmqClient.ConsumeMessages(conversionExchange, conversionKey, conversionQueue)
	if err != nil {
		slog.Error("Error consuming messages from RabbitMQ", slog.String("error", err.Error()))
	}

	for d := range msgs {
		go func(delivery amqp.Delivery) {
			converter.Handle(delivery, conversionExchange, confirmationKey, confirmationQueue)
		}(d)
	}
}
