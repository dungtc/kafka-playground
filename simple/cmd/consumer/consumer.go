package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/dungtc/kafka-playground/simple"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	brokerList        = kingpin.Flag("brokerList", "List of brokers").Default("localhost:9092").Strings()
	topics            = kingpin.Flag("topic", "Topic name").Default("test").String()
	partition         = kingpin.Flag("partition", "Partition").Default("0").String()
	groupId           = kingpin.Flag("groupId", "Group id").Default("my-application").String()
	offsetType        = kingpin.Flag("offsetType", "Offset Type (OffsetNewest | OffsetOldest)").Default("-1").Int()
	messageCountStart = kingpin.Flag("messageCountStart", "Message counter start from:").Int()
	maxRetry          = kingpin.Flag("maxRetry", "Retry limit").Default("5").Int()
	version           = kingpin.Flag("version", "Kafka version").Default("2.5.0").String()
)

func main() {
	consumerGroup, err := simple.NewConsumerGroup()
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := consumerGroup.Consume(ctx, strings.Split(*topics, ","), &simple.Consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				log.Println("cancel yes")
				return
			}
			// consumer.ready = make(chan bool)
		}
	}()

	// <-consumer.ready // Await till the consumer has been set up
	select {
	case <-ctx.Done():
		log.Println("terminating: context cancelled")
	case <-signals:
		log.Println("terminating: via signal")
	}

	cancel()
	wg.Wait()
	if err = consumerGroup.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}

	log.Println("Processed", *messageCountStart, "messages")
}
