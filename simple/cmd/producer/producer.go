package main

import (
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/dungtc/kafka-playground/simple"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	brokerList               = kingpin.Flag("brokerList", "List of brokers").Default("0.0.0.0:9092").Strings()
	topic                    = kingpin.Flag("topic", "Topic name").Default("youtube").String()
	version                  = kingpin.Flag("version", "Kafka version").Default("2.5.0").String()
	maxRetry                 = kingpin.Flag("maxRetry", "Retry limit").Default("5").Int()
	partitions        *int32 = kingpin.Flag("partitions", "Partitions").Default("1").Int32()
	replicationFactor *int16 = kingpin.Flag("replicas", "Replication factor broker").Default("1").Int16()
)

func main() {
	kingpin.Parse()

	// init administration
	clusterAdmin, err := simple.NewClusterAdmin(*version, *brokerList...)
	if err != nil {
		panic(err)
	}

	if err := clusterAdmin.NewTopic(*topic, *partitions, *replicationFactor); err != nil {
		panic(err)
	}

	producer, err := simple.NewProducer(*maxRetry, *brokerList...)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Panic(err)
		}
	}()
	for i := 0; i < 10; i++ {
		partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
			Topic: *topic,
			Key:   sarama.StringEncoder(fmt.Sprintf("id_%v", i)),
			Value: sarama.StringEncoder(fmt.Sprintf("hello %v", i)),
		})
		if err != nil {
			log.Panic(err)
		}
		log.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", *topic, partition, offset)
	}
}
