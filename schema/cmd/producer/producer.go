package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/Shopify/sarama"
	"github.com/dungtc/kafka-playground/schema"
	"github.com/dungtc/kafka-playground/schema/domain"
	"github.com/dungtc/kafka-playground/simple"
	"github.com/google/uuid"
	"github.com/riferrei/srclient"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	brokerList               = kingpin.Flag("brokerList", "List of brokers").Default("localhost:9092").Strings()
	topic                    = kingpin.Flag("topic", "Topic name").Default("youtube").String()
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

	// schema
	schemaRegistryClient := schema.NewSchema("http://localhost:8081")
	// schema, err := schemaRegistryClient.GetLastestSchemaVersion(*topic)
	// if schema != nil || err != nil {
	// 	schema = schemaRegistryClient.CreateSchemaFromFile(*topic, "schema-example.avsc", srclient.Avro, false)
	// }
	fmt.Println("voo")
	schema, err := schemaRegistryClient.CreateSchemaFromFile(*topic, "schema-example.avsc", srclient.Avro.String(), false)
	if err != nil {
		panic(err)
	}

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
		data := domain.Video{VideoId: strconv.Itoa(i), CategoryID: uuid.New().String(), ETag: ""}
		// data2 := ComplexType{ID: 1, Name: "Gopher"}

		payload, _ := json.Marshal(data)

		partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
			Topic: *topic,
			Key:   sarama.ByteEncoder([]byte(uuid.New().String())),
			Value: sarama.ByteEncoder(schema.SerializeData(payload)),
		})
		if err != nil {
			log.Panic(err)
		}
		log.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", *topic, partition, offset)
	}
}
