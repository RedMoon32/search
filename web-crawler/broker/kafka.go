package broker

import (
	"context"
	"log"

	"github.com/IBM/sarama"
)

type Broker interface {
	Consume() (<-chan string, error)
	Produce(message string) error
}

type KafkaBroker struct {
	producer       sarama.SyncProducer
	consumerGroup  sarama.ConsumerGroup
	topic          string
	ready          chan bool
	messageChannel chan string
}

func NewKafkaBroker(brokers, topic string) Broker {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	producer, err := sarama.NewSyncProducer([]string{brokers}, config)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}

	consumerGroup, err := sarama.NewConsumerGroup([]string{brokers}, "web-crawler-group", config)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer group: %v", err)
	}

	return &KafkaBroker{
		producer:       producer,
		consumerGroup:  consumerGroup,
		topic:          topic,
		ready:          make(chan bool),
		messageChannel: make(chan string, 100),
	}
}

func (k *KafkaBroker) Consume() (<-chan string, error) {
	ctx := context.Background()
	go func() {
		for {
			if err := k.consumerGroup.Consume(ctx, []string{k.topic}, k); err != nil {
				log.Fatalf("Error consuming messages: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
			k.ready = make(chan bool)
		}
	}()
	<-k.ready
	return k.messageChannel, nil
}

func (k *KafkaBroker) Produce(message string) error {
	_, _, err := k.producer.SendMessage(&sarama.ProducerMessage{
		Topic: k.topic,
		Value: sarama.StringEncoder(message),
	})
	return err
}

func (k *KafkaBroker) Setup(sarama.ConsumerGroupSession) error {
	close(k.ready)
	return nil
}

func (k *KafkaBroker) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (k *KafkaBroker) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		k.messageChannel <- string(msg.Value)
		session.MarkMessage(msg, "")
	}
	return nil
}
