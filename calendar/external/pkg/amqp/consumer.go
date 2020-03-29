package amqp

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/domain"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

var (
	ErrConsumerStop           = errors.New("consumer stops")
	ErrConsumerWithoutHandler = errors.New("handler not defined")
)

// Consumer implements events.Consumer
type Consumer struct {
	logger       *zap.Logger
	errorChan    chan<- error
	connection   *amqp.Connection
	channel      *amqp.Channel
	queue        *amqp.Queue
	exchangeName string
}

func NewConsumer(
	logger *zap.Logger,
	errorChan chan<- error,
	dsn string,
	exchangeName,
	queueName string,
	waitGroup *sync.WaitGroup,
) (*Consumer, error) {

	consumer := &Consumer{
		exchangeName: exchangeName,
		logger:       logger,
		errorChan:    errorChan,
	}

	var err error
	consumer.connection, err = amqp.Dial(dsn)
	if err != nil {
		return nil, err
	}

	consumer.channel, err = consumer.connection.Channel()
	if err != nil {
		consumer.connection.Close()
		return nil, err
	}

	err = consumer.channel.ExchangeDeclare(
		exchangeName,
		amqp.ExchangeFanout,
		false,
		false,
		false,
		false,
		amqp.Table{},
	)
	if err != nil {
		consumer.Close()
		return nil, err
	}

	queue, err := consumer.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		true,
		amqp.Table{},
	)
	if err != nil {
		consumer.Close()
		return nil, err
	}

	consumer.queue = &queue

	errChan := make(chan *amqp.Error)
	consumer.channel.NotifyClose(errChan)

	waitGroup.Add(1)
	go func(waitGroup *sync.WaitGroup) {
		defer waitGroup.Done()
		if err, ok := <-errChan; ok && err != nil {
			errorChan <- errors.New(err.Error())
		}
	}(waitGroup)

	return consumer, nil
}

func (c *Consumer) Consume(ctx context.Context, handler func(event domain.Event) error) error {

	err := c.channel.QueueBind(
		c.queue.Name,
		"",
		c.exchangeName,
		false,
		nil,
	)
	if err != nil {
		c.logger.Error("bind queue to exchange error", zap.Error(err))
		return err
	}

	if handler == nil {
		return ErrConsumerWithoutHandler
	}

	deliveryChan, err := c.channel.Consume(
		c.queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		c.logger.Error("consume error", zap.Error(err))
		return err
	}

	waitGroup := &sync.WaitGroup{}

	waitGroup.Add(1)
	go func(waitGroup *sync.WaitGroup) {
		defer waitGroup.Done()
		for {
			select {
			case <-ctx.Done():
				c.logger.Debug("consumer close")
				return
			case msg, ok := <-deliveryChan:
				if ok {
					c.logger.Debug("income message")
					var jsonEvent Event
					if err := json.Unmarshal(msg.Body, &jsonEvent); err != nil {
						c.logger.Error("unmarshal error",
							zap.ByteString("message", msg.Body),
							zap.Error(err),
						)
						c.errorChan <- err
						return
					}

					event, err := eventFromJson(jsonEvent)
					if err != nil {
						c.logger.Error("parse event error",
							zap.Reflect("json event", jsonEvent),
							zap.Error(err),
						)
						c.errorChan <- err
						return
					}

					if err := handler(event); err != nil {
						c.logger.Warn("handle event error", zap.Error(err))
					}
				} else {
					c.errorChan <- ErrConsumerStop
					return
				}
			}
		}
	}(waitGroup)

	waitGroup.Wait()

	return nil
}

func (c *Consumer) Close() error {
	if err := c.channel.Close(); err != nil {
		return err
	}

	if err := c.connection.Close(); err != nil {
		return err
	}

	return nil
}
