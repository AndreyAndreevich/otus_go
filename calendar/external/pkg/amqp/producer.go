package amqp

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/domain"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

// Producer implements events.Producer
type Producer struct {
	logger     *zap.Logger
	connection *amqp.Connection
	channel    *amqp.Channel
	exchange   string
}

// NewProducer - create rabbit producer
func NewProducer(logger *zap.Logger,
	errorChan chan<- error,
	dsn string,
	exchange string,
	waitGroup *sync.WaitGroup,
) (result *Producer, err error) {

	producer := &Producer{
		exchange: exchange,
		logger:   logger,
	}

	if producer.connection, err = amqp.Dial(dsn); err != nil {
		logger.Error("parse amqp dsn error", zap.Error(err))
		return nil, err
	}

	if producer.channel, err = producer.connection.Channel(); err != nil {
		producer.connection.Close()
		logger.Error("create channel error", zap.Error(err))
		return nil, err
	}

	err = producer.channel.ExchangeDeclare(
		exchange,
		amqp.ExchangeFanout,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		producer.channel.Close()
		producer.connection.Close()

		logger.Error("declare exchange error", zap.Error(err))
		return nil, err
	}

	errChan := make(chan *amqp.Error)
	producer.channel.NotifyClose(errChan)

	waitGroup.Add(1)
	go func(waitGroup *sync.WaitGroup) {
		defer waitGroup.Done()
		if err, ok := <-errChan; ok && err != nil {
			errorChan <- errors.New(err.Error())
		}
	}(waitGroup)

	return producer, nil
}

// Publish event to exchange
func (p *Producer) Publish(event domain.Event) error {
	jsonEvent := eventToJSON(event)

	body, err := json.Marshal(jsonEvent)
	if err != nil {
		p.logger.Error("cannot marshal event", zap.Error(err))
		return err
	}

	err = p.channel.Publish(
		p.exchange,
		"",
		false,
		false,
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            body,
			DeliveryMode:    amqp.Persistent,
			Priority:        0,
		},
	)
	if err != nil {
		p.logger.Error("cannot publish event", zap.Error(err), zap.Any("event", event))
		return err
	}

	return nil
}

// Close producer
func (p *Producer) Close() error {
	if err := p.channel.Close(); err != nil {
		return err
	}

	if err := p.connection.Close(); err != nil {
		return err
	}

	return nil
}
