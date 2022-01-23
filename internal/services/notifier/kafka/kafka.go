package kafka

import (
	"io"
	"log"
	"strconv"
	"time"

	"user-crud/internal/config"
	"user-crud/pkg/closer"

	"github.com/Shopify/sarama"
)

type Client struct {
	producer sarama.SyncProducer
}

func NewClient(config *config.Config) (*Client, error) {
	conf := sarama.NewConfig()
	conf.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(config.Kafka.Brokers, conf)

	if err != nil {
		return nil, err
	}

	closer.Add(producer.Close)

	return &Client{
		producer: producer,
	}, nil
}

func (s *Client) Publish(topic string, r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		log.Printf("kafka: failed to read message %v", err)
		return err
	}

	_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
		Topic:     topic,
		Key:       sarama.StringEncoder(strconv.FormatInt(1, 10)),
		Value:     sarama.ByteEncoder(data),
		Timestamp: time.Now(),
	})
	if err != nil {
		log.Printf("kafka: failed to send message %v", err)
	}

	return err
}
