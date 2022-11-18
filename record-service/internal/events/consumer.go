package events

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"record-service/internal/logger"
	"record-service/internal/services"
	"record-service/pkg/models"
)

type Consumer struct {
	conn         *amqp.Connection
	exchange     string
	itemsService services.ItemsService
}

func NewConsumer(conn *amqp.Connection, exchange string, itemsService services.ItemsService) (Consumer, error) {
	consumer := Consumer{
		conn:         conn,
		exchange:     exchange,
		itemsService: itemsService,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (c *Consumer) setup() error {
	channel, err := c.conn.Channel()
	if err != nil {
		logger.Error("Error while setup consumer", err)
		return err
	}

	return declareExchange(channel, c.exchange)
}

func (c *Consumer) Listen(topics []string) error {
	ch, err := c.conn.Channel()
	if err != nil {
		logger.Error("Error while listening", err)
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		logger.Error("Error while declaring random queue", err)
		return err
	}

	for _, topic := range topics {
		err = ch.QueueBind(q.Name, topic, c.exchange, false, nil)
		if err != nil {
			logger.Error("Error while queue bind", err)
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		logger.Error("Error while consuming channel", err)
		return err
	}

	forever := make(chan bool)
	go func() {
		for message := range messages {
			logger.Info("Received message on queue")
			var item models.TodoItem
			err = json.Unmarshal(message.Body, &item)
			if err != nil {
				logger.Error("Error while unmarshalling message", err)
			}

			go c.handleTodoItem(item)
		}
	}()

	<-forever
	return nil
}

func (c *Consumer) handleTodoItem(item models.TodoItem) {
	err := c.itemsService.Insert(item)
	if err != nil {
		logger.Error("Error while inserting item", err)
	}
}
