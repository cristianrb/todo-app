package main

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"net/http"
	"os"
	"record-service/internal/events"
	"record-service/internal/handlers"
	"record-service/internal/logger"
	"record-service/internal/services"
	"time"
)

const (
	PORT              = "8081"
	MONGO_USERNAME    = "MONGO_USERNAME"
	MONGO_PASSWORD    = "MONGO_PASSWORD"
	MONGO_URL         = "MONGO_URL"
	DB_NAME           = "DB_NAME"
	COLLECTION_NAME   = "COLLECTION_NAME"
	RABBITMQ_USERNAME = "RABBITMQ_USERNAME"
	RABBITMQ_PASSWORD = "RABBITMQ_PASSWORD"
	RABBITMQ_HOST     = "RABBITMQ_HOST"
	RABBIT_EXCHANGE   = "RABBIT_EXCHANGE"
	RABBIT_TOPIC      = "RABBIT_TOPIC"
)

func main() {
	logger.Info(fmt.Sprintf("Starting broker service on port %s", PORT))
	// RabbitMQ
	rabbitmqHost := os.Getenv(RABBITMQ_HOST)
	rabbitmqUsername := os.Getenv(RABBITMQ_USERNAME)
	rabbitmqPassword := os.Getenv(RABBITMQ_PASSWORD)
	rabbitConn, err := connectToRabbitMQ(rabbitmqHost, rabbitmqUsername, rabbitmqPassword)
	if err != nil {
		panic(err)
	}
	defer rabbitConn.Close()

	// MongoDB
	mongoUrl := os.Getenv(MONGO_URL)
	mongoUsername := os.Getenv(MONGO_USERNAME)
	mongoPassword := os.Getenv(MONGO_PASSWORD)
	dbName := os.Getenv(DB_NAME)
	collectionName := os.Getenv(COLLECTION_NAME)
	mongoClient, err := connectToMongo(mongoUrl, mongoUsername, mongoPassword)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	itemsService := services.NewItemsService(mongoClient, dbName, collectionName)
	handlerConfig := handlers.NewHandlerConfig(&http.Client{}, itemsService)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: handlerConfig.Routes(),
	}
	go srv.ListenAndServe()

	exchange := os.Getenv(RABBIT_EXCHANGE)
	topic := os.Getenv(RABBIT_TOPIC)
	consumer, err := events.NewConsumer(rabbitConn, exchange, itemsService)
	if err != nil {
		panic(err)
	}

	err = consumer.Listen([]string{topic})
	if err != nil {
		logger.Error("Error while listening to rabbit mq", err)
	}
}

func connectToMongo(mongoUrl, username, password string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoUrl)
	clientOptions.SetAuth(options.Credential{
		Username: username,
		Password: password,
	})

	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logger.Error("Unable to connect to mongo", err)
		return nil, err
	}

	logger.Info("Connected to mongo")
	return c, nil
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
