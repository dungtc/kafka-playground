package main

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/Shopify/sarama"
	"github.com/dungtc/kafka-playground/schema"
	"github.com/dungtc/kafka-playground/schema/domain"
	"github.com/dungtc/kafka-playground/simple"
	"github.com/google/uuid"
	"github.com/riferrei/srclient"
	"github.com/sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	brokerList               = kingpin.Flag("brokerList", "List of brokers").Default("localhost:9092").Strings()
	topic                    = kingpin.Flag("topic", "Topic name").Default("schema-1").String()
	version                  = kingpin.Flag("version", "Kafka version").Default("2.5.0").String()
	maxRetry                 = kingpin.Flag("maxRetry", "Retry limit").Default("5").Int()
	partitions        *int32 = kingpin.Flag("partitions", "Partitions").Default("1").Int32()
	replicationFactor *int16 = kingpin.Flag("replicas", "Replication factor broker").Default("1").Int16()
)

type ComplexType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	kingpin.Parse()

	// init schema client
	schemaRegistryClient := schema.NewSchema("http://localhost:8081")

	// get latest schema
	logrus.Info("Action get lastest schema version")
	schema, err := schemaRegistryClient.GetLastestSchemaVersion(*topic)
	if schema != nil || err != nil {
		logrus.Info("Action create schema from file")

		// create schema if not exist
		schema, err = schemaRegistryClient.CreateSchemaFromFile(*topic, "avro-schema/schema-example.avsc", srclient.Avro.String(), false)
		if err != nil {
			panic(err)
		}
	}

	// init administration
	clusterAdmin, err := simple.NewClusterAdmin(*version, *brokerList...)
	if err != nil {
		panic(err)
	}

	// create topic if not exist
	if err := clusterAdmin.NewTopic(*topic, *partitions, *replicationFactor); err != nil {
		panic(err)
	}

	// create new producer
	producer, err := simple.NewProducer(*maxRetry, *brokerList...)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Panic(err)
		}
	}()

	// send message
	for i := 0; i < 10; i++ {
		data := domain.Video{VideoId: strconv.Itoa(i), CategoryID: uuid.New().String(), ETag: ""}

		payload, _ := json.Marshal(data)
		partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
			Topic: *topic,
			Key:   sarama.ByteEncoder([]byte(uuid.New().String())),
			Value: sarama.ByteEncoder(schema.SerializeData(payload)),
		})
		if err != nil {
			log.Panic(err)
		}
		logrus.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d), schema %v\n", *topic, partition, offset, schema.ID())
	}
}
