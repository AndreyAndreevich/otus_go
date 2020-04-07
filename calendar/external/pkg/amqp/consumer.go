package amqp

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/AndreyAndreevich/otus_go/calendar/internal/domain"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

var (
	// ErrConsumerStop - consumer stop error
	ErrConsumerStop = errors.New("consumer stops")
	// ErrConsumerWithoutHandler - empty handler
	ErrConsumerWithoutHandler = errors.New("handler not defined")
)

// Consumer implements events.Consumer
type Consumer struct {
	logger       *zap.Logger
	errorChan    chan *amqp.Error
	connection   *amqp.Connection
	channel      *amqp.Channel
	queue        *amqp.Queue
	exchangeName string
	queueName    string
	dsn          string
}

// NewConsumer - create new rabbit consumer
func NewConsumer(
	logger *zap.Logger,
	dsn string,
	exchangeName,
	queueName string,
) (*Consumer, error) {

	consumer := &Consumer{
		logger:       logger,
		errorChan:    make(chan *amqp.Error),
		exchangeName: exchangeName,
		queueName:    queueName,
		dsn:          dsn,
	}

	if err := consumer.reconnect(); err != nil {
		return nil, err
	}

	return consumer, nil
}

func (c *Consumer) reconnect() (err error) {
	c.connection, err = amqp.Dial(c.dsn)
	if err != nil {
		return err
	}

	c.channel, err = c.connection.Channel()
	if err != nil {
		c.connection.Close()
		return err
	}

	err = c.channel.ExchangeDeclare(
		c.exchangeName,
		amqp.ExchangeFanout,
		false,
		false,
		false,
		false,
		amqp.Table{},
	)
	if err != nil {
		c.Close()
		return err
	}

	queue, err := c.channel.QueueDeclare(
		c.queueName,
		true,
		false,
		false,
		true,
		amqp.Table{},
	)
	if err != nil {
		c.Close()
		return err
	}

	c.channel.NotifyClose(c.errorChan)
	c.queue = &queue

	return nil
}

// Consume - consume messages and handle events
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
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		c.logger.Error("consume error", zap.Error(err))
		return err
	}

	for {
		select {
		case <-ctx.Done():
			c.logger.Debug("consumer close")
			return nil
		case err := <-c.errorChan:
			c.logger.Warn("channel error", zap.Error(err))
			if err := c.reconnect(); err != nil {
				return err
			}
		case msg, ok := <-deliveryChan:
			if !ok {
				return nil
			}

			var jsonEvent Event
			if err := json.Unmarshal(msg.Body, &jsonEvent); err != nil {
				c.logger.Error("unmarshal error",
					zap.ByteString("message", msg.Body),
					zap.Error(err),
				)
				return err
			}

			event, err := eventFromJSON(jsonEvent)
			if err != nil {
				c.logger.Error("parse event error",
					zap.Reflect("json event", jsonEvent),
					zap.Error(err),
				)
				return err
			}

			if err := handler(event); err != nil {
				c.logger.Warn("handle event error", zap.Error(err))
			}
		}
	}
}

// Close consumer
func (c *Consumer) Close() error {
	if err := c.channel.Close(); err != nil {
		return err
	}

	if err := c.connection.Close(); err != nil {
		return err
	}

	return nil
}
