package events

import (
	"broker-service/internal/logger"
	"broker-service/pkg/models"
	"context"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	connection *amqp.Connection
	exchange   string
	topic      string
}

func (e *Emitter) setup() error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	return declareExchange(channel, e.exchange)
}

func NewEventEmitter(conn *amqp.Connection, exchange, topic string) (Emitter, error) {
	emitter := Emitter{
		connection: conn,
		exchange:   exchange,
		topic:      topic,
	}

	err := emitter.setup()
	if err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}

func (e *Emitter) PushToQueue(item models.TodoItem) error {
	itemBytes, err := json.Marshal(&item)
	if err != nil {
		return err
	}
	err = e.publish(string(itemBytes))
	if err != nil {
		return err
	}

	return nil
}

func (e *Emitter) publish(event string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	logger.Info("Pushing to channel")
	err = channel.PublishWithContext(
		context.TODO(),
		e.exchange,
		e.topic,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)
	if err != nil {
		logger.Error("Could not push to channel", err)
	}
	return nil
}
