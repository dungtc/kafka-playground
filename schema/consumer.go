package schema

import (
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

type Consumer struct {
	ready          chan bool
	schemaRegistry SchemaRegistry
}

func NewConsumer(registry SchemaRegistry) Consumer {
	return Consumer{
		ready:          make(chan bool),
		schemaRegistry: registry,
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	log.Println("Consumer listening")
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

		schema, err := consumer.schemaRegistry.GetSchemaFromMsg(message.Value)
		if err != nil {
			return fmt.Errorf("Error getting the schema with id '%d' %s", schema.ID(), err)
		}

		// deserialize data
		payload := schema.DeserializeData(message.Value)

		logrus.Printf("Message value = %s, schema_id = %d, timestamp = %v, topic = %s, partition = %v", string(payload), schema.ID(), message.Timestamp, message.Topic, message.Partition)
		session.MarkMessage(message, "")
	}

	return nil
}
