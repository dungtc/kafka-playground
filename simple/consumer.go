package simple

import (
	"log"

	"github.com/Shopify/sarama"
)

func NewConsumerGroup(ver string, groupId string, brokerList ...string) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	version, err := sarama.ParseKafkaVersion(ver)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}
	config.Version = version
	config.Consumer.Return.Errors = true
	client, err := sarama.NewClient(brokerList, config)
	if err != nil {
		log.Panic(err)
	}

	// Start a new consumer group
	group, err := sarama.NewConsumerGroupFromClient(groupId, client)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return group, nil
}

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready chan bool
}

func NewConsumer() Consumer {
	return Consumer{
		ready: make(chan bool),
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	// close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s, partition = %v", string(message.Value), message.Timestamp, message.Topic, message.Partition)
		session.MarkMessage(message, "")
	}

	return nil
}
