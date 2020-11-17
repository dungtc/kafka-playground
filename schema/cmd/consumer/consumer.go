package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/dungtc/kafka-playground/schema"
	"github.com/dungtc/kafka-playground/simple"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	brokerList        = kingpin.Flag("brokerList", "List of brokers").Default("localhost:9092").Strings()
	topics            = kingpin.Flag("topic", "Topic name").Default("schema-1").String()
	partition         = kingpin.Flag("partition", "Partition").Default("0").String()
	groupId           = kingpin.Flag("groupId", "Group id").Default("schema-application").String()
	offsetType        = kingpin.Flag("offsetType", "Offset Type (OffsetNewest | OffsetOldest)").Default("-1").Int()
	messageCountStart = kingpin.Flag("messageCountStart", "Message counter start from:").Int()
	maxRetry          = kingpin.Flag("maxRetry", "Retry limit").Default("5").Int()
	version           = kingpin.Flag("version", "Kafka version").Default("2.5.0").String()
	assignor          = kingpin.Flag("assignor", "Kafka partition assignor").Default("sticky").String()
)

func main() {
	kingpin.Parse()

	// init schema client
	schemaRegistryClient := schema.NewSchema("http://localhost:8081")

	// create kafka consumer group by id
	logrus.Println("Init consumer group")
	consumerGroup, err := simple.NewConsumerGroup(*version, *groupId, *assignor, *brokerList...)
	if err != nil {
		panic(err)
	}

	// init context signal cancellation
	ctx, cancel := context.WithCancel(context.Background())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// create consumer
	consumer := schema.NewConsumer(schemaRegistryClient)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := consumerGroup.Consume(ctx, strings.Split(*topics, ","), &consumer); err != nil {
				logrus.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				logrus.Println("cancel yes")
				return
			}
			// consumer.ready = make(chan bool)
		}
	}()

	// <-consumer.ready // Await till the consumer has been set up
	select {
	case <-ctx.Done():
		logrus.Println("terminating: context cancelled")
	case <-signals:
		logrus.Println("terminating: via signal")
	}

	cancel()
	wg.Wait()
	if err = consumerGroup.Close(); err != nil {
		logrus.Panicf("Error closing client: %v", err)
	}

	logrus.Println("Processed", *messageCountStart, "messages")
}
