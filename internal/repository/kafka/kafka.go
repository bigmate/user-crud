package kafka

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"user-crud/internal/config"

	"github.com/bigmate/closer"

	"github.com/Shopify/sarama"
)

type Client struct {
	producer sarama.AsyncProducer
	cli      sarama.Client
}

func (c *Client) Ping(_ context.Context) error {
	for _, broker := range c.cli.Brokers() {
		ok, err := broker.Connected()
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("no connection with broker: %v", broker.ID())
		}
	}

	return nil
}

func NewClient(config *config.Config) (*Client, error) {
	conf := sarama.NewConfig()
	conf.Producer.Return.Successes = true

	cli, err := sarama.NewClient(config.Kafka.Brokers, conf)
	if err != nil {
		return nil, err
	}

	producer, err := sarama.NewAsyncProducerFromClient(cli)
	if err != nil {
		return nil, err
	}

	closer.Add(producer.Close)
	closer.Add(cli.Close)

	go func() {
		for sendErr := range producer.Errors() {
			log.Printf("kafka: failed to send message %v", sendErr)
		}
	}()

	return &Client{
		producer: producer,
		cli:      cli,
	}, nil
}

func (c *Client) Publish(topic string, r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		log.Printf("kafka: failed to read message %v", err)
		return err
	}

	c.producer.Input() <- &sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.StringEncoder(strconv.FormatInt(1, 10)),
		Value:     sarama.ByteEncoder(data),
		Timestamp: time.Now(),
	}

	return err
}
