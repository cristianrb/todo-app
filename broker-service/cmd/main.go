package main

import (
	"broker-service/internal/events"
	"broker-service/internal/handlers"
	"broker-service/internal/logger"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"math"
	"net/http"
	"os"
	"time"
)

const (
	PORT              = "8080"
	RABBITMQ_USERNAME = "RABBITMQ_USERNAME"
	RABBITMQ_PASSWORD = "RABBITMQ_PASSWORD"
	RABBITMQ_HOST     = "RABBITMQ_HOST"
	RABBIT_EXCHANGE   = "RABBIT_EXCHANGE"
	RABBIT_TOPIC      = "RABBIT_TOPIC"
)

func main() {
	logger.Info(fmt.Sprintf("Starting broker service on port %s", PORT))
	rabbitmqHost := os.Getenv(RABBITMQ_HOST)
	rabbitmqUsername := os.Getenv(RABBITMQ_USERNAME)
	rabbitmqPassword := os.Getenv(RABBITMQ_PASSWORD)
	exchange := os.Getenv(RABBIT_EXCHANGE)
	topic := os.Getenv(RABBIT_TOPIC)
	rabbitConn, err := connectToRabbitMQ(rabbitmqHost, rabbitmqUsername, rabbitmqPassword)
	if err != nil {
		panic(err)
	}
	defer rabbitConn.Close()

	emitter, err := events.NewEventEmitter(rabbitConn, exchange, topic)
	if err != nil {
		panic(err)
	}
	handlerConfig := handlers.NewHandlerConfig(&http.Client{}, emitter)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: handlerConfig.Routes(),
	}
	err = srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func connectToRabbitMQ(host, username, password string) (*amqp.Connection, error) {
	var attempts int64
	var backOff = 1 * time.Second

	for {
		c, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s", username, password, host))
		if err != nil {
			logger.Info("RabbitMQ not yet ready...")
			attempts++
		} else {
			logger.Info("Connected to RabbitMQ")
			return c, nil
		}

		if attempts > 5 {
			logger.Error("Unable to connect to RabbitMQ", err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(attempts), 2)) * time.Second
		logger.Info("Backing off...")
		time.Sleep(backOff)
	}

}
