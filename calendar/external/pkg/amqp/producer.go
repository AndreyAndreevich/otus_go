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
	logger       *zap.Logger
	errorChan    chan *amqp.Error
	connection   *amqp.Connection
	channel      *amqp.Channel
	exchangeName string
	dsn          string
}

// NewProducer - create rabbit producer
func NewProducer(logger *zap.Logger,
	errorChan chan<- error,
	dsn string,
	exchange string,
	waitGroup *sync.WaitGroup,
) (result *Producer, err error) {

	producer := &Producer{
		logger:       logger,
		errorChan:    make(chan *amqp.Error),
		exchangeName: exchange,
		dsn:          dsn,
	}

	if err := producer.reconnect(); err != nil {
		return nil, err
	}

	waitGroup.Add(1)
	go func(waitGroup *sync.WaitGroup) {
		defer waitGroup.Done()
		for {
			if err, ok := <-producer.errorChan; ok && err != nil {
				logger.Warn("channel error", zap.Error(err))
				if err := producer.reconnect(); err != nil {
					errorChan <- errors.New(err.Error())
					return
				}
			}
		}
	}(waitGroup)

	return producer, nil
}

func (p *Producer) reconnect() (err error) {
	if p.connection, err = amqp.Dial(p.dsn); err != nil {
		p.logger.Error("parse amqp dsn error", zap.Error(err))
		return err
	}

	if p.channel, err = p.connection.Channel(); err != nil {
		p.connection.Close()
		p.logger.Error("create channel error", zap.Error(err))
		return err
	}

	err = p.channel.ExchangeDeclare(
		p.exchangeName,
		amqp.ExchangeFanout,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		p.Close()

		p.logger.Error("declare exchange error", zap.Error(err))
		return err
	}

	p.channel.NotifyClose(p.errorChan)

	return nil
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
		p.exchangeName,
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
