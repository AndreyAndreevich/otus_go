package amqp

import (
	"context"
	"errors"
	"sync"

	"github.com/streadway/amqp"
	"gitlab.octafx.com/develop/ms-ctrader-parser/events"
	libevents "gitlab.octafx.com/go-libs/events"
	"go.uber.org/zap"
)

var (
	ErrConsumerStop           = errors.New("consumer stops")
	ErrConsumerWithoutHandler = errors.New("handler not defined")
)

// Consumer implements events.Consumer
type Consumer struct {
	logger       *zap.Logger
	waitGroup    *sync.WaitGroup
	ctx          context.Context
	errorChan    chan<- error
	connection   *amqp.Connection
	channel      *amqp.Channel
	queue        *amqp.Queue
	exchangeName string
	handler      libevents.Handler
}

func (c *Consumer) Consume(events []string, handler libevents.Handler) error {
	c.handler = handler

	// Bind Queue for each event from eventsToListen
	for _, eventName := range events {
		err := c.channel.QueueBind(
			c.queue.Name,
			eventName, // routingKey equal with eventName
			c.exchangeName,
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}

	if c.handler == nil {
		return ErrConsumerWithoutHandler
	}

	msgs, err := c.channel.Consume(
		c.queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	c.waitGroup.Add(1)
	go func() {
		defer c.waitGroup.Done()
		for {
			select {
			case <-c.ctx.Done():
				c.logger.Debug("consumer worker close")
				return
			case msg, ok := <-msgs:
				if ok {
					c.logger.Debug("income message")
					c.handleMessage(&msg)
				} else {
					c.errorChan <- ErrConsumerStop
					return
				}
			}
		}
	}()

	return nil
}

func NewConsumer(
	logger *zap.Logger,
	waitGroup *sync.WaitGroup,
	ctx context.Context,
	errorChan chan<- error,
	address Address,
	exchangeName,
	queueName string,
) (res libevents.Consumer, err error) {

	logger.Debug("RabbitMQ events consumer connecting", zap.String("address", address.SafeCreds()))
	consumer := &Consumer{
		exchangeName: exchangeName,
		logger:       logger,
		waitGroup:    waitGroup,
		ctx:          ctx,
		errorChan:    errorChan,
	}

	consumer.connection, err = amqp.Dial(address.String())
	if err != nil {
		return nil, err
	}

	consumer.channel, err = consumer.connection.Channel()
	if err != nil {
		return nil, err
	}

	err = consumer.channel.ExchangeDeclare(
		exchangeName,
		amqp.ExchangeFanout,
		true,
		false,
		false,
		false,
		amqp.Table{},
	)
	if err != nil {
		return nil, err
	}

	queue, err := consumer.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		amqp.Table{},
	)
	if err != nil {
		return nil, err
	}
	consumer.queue = &queue

	// OnClose
	cerr := make(chan *amqp.Error)
	consumer.channel.NotifyClose(cerr)

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		select {
		case <-ctx.Done():
			return
		case err, ok := <-cerr:
			if ok && err != nil {
				errorChan <- errors.New(err.Error())
			}
		}
	}()

	return consumer, nil
}

func (c *Consumer) handleMessage(msg *amqp.Delivery) {
	// Ack all messages after handling
	defer func() {
		if err := msg.Ack(false); err != nil {
			c.logger.Warn("error on ACK message",
				zap.Error(err),
				zap.Uint64("delivery_tag", msg.DeliveryTag),
				zap.ByteString("message", msg.Body))
		}
	}()

	if len(msg.Body) == 0 {
		c.logger.Warn("empty message of amqp message", zap.Any("msg", msg))
		return
	}

	// parse event
	event := &libevents.Event{
		Event:   msg.Headers["messageType"].(string),
		Payload: msg.Body,
	}

	// handle event
	c.logger.Debug("handle event", zap.String("eventType", event.Event))

	err := c.handler(event)
	switch err {
	case events.ErrUnexpectedEvent:
	case nil:
	default:
		c.logger.Warn("eventHandle returns error",
			zap.Error(err),
			zap.Any("event", event),
			zap.Uint64("delivery_tag", msg.DeliveryTag))
	}
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
