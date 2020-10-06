package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/dungtc/kafka-playground/youtube-stream/domain"
	"github.com/dungtc/kafka-playground/youtube-stream/repository"
	"github.com/dungtc/kafka-playground/youtube-stream/setting"
	"github.com/dungtc/kafka-playground/youtube-stream/storage"

	"strings"
	"sync"

	"github.com/Shopify/sarama"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	brokerList        = kingpin.Flag("brokerList", "List of brokers").Default("localhost:9092").Strings()
	topics            = kingpin.Flag("topic", "Topic name").Default("youtube").String()
	partition         = kingpin.Flag("partition", "Partition").Default("0").String()
	groupId           = kingpin.Flag("groupId", "Group id").Default("my-application").String()
	offsetType        = kingpin.Flag("offsetType", "Offset Type (OffsetNewest | OffsetOldest)").Default("-1").Int()
	messageCountStart = kingpin.Flag("messageCountStart", "Message counter start from:").Int()
	maxRetry          = kingpin.Flag("maxRetry", "Retry limit").Default("5").Int()
	version           = kingpin.Flag("version", "Kafka version").Default("2.5.0").String()
)

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready             chan bool
	youtubeRepository repository.VideoRepository
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

		video := repository.Video{}
		if err := json.Unmarshal(message.Value, &video); err != nil {
			log.Print("err", err)
		}

		consumer.youtubeRepository.Save(session.Context(), &domain.Video{VideoId: video.Id, CategoryID: video.Snippet.CategoryId, ETag: video.Etag})
		session.MarkMessage(message, "")
	}

	return nil
}

func main() {
	log.Println("Start consumer")

	kingpin.Parse()

	conf := setting.Configuration{}

	// Load configuration
	err := setting.LoadConfiguration("env", "config", "toml", &conf)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	dbstring := fmt.Sprintf(
		"user=%s host=%s port=%s dbname=%s sslmode=disable password=%s",
		conf.Database.User,
		conf.Database.Host,
		conf.Database.Port,
		conf.Database.Name,
		conf.Database.Passsword,
	)

	db, err := storage.NewStorage(dbstring)
	if err != nil {
		panic(err)
	}
	youtubeRepository := repository.NewVideoRepository(*db)

	config := sarama.NewConfig()
	version, err := sarama.ParseKafkaVersion(*version)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", err)
	}
	config.Version = version
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer := Consumer{
		ready:             make(chan bool),
		youtubeRepository: youtubeRepository,
	}

	// Start a new consumer group
	client, err := sarama.NewConsumerGroup(*brokerList, *groupId, config)
	if err != nil {
		log.Println(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			log.Println("Consume")
			if err := client.Consume(ctx, strings.Split(*topics, ","), &consumer); err != nil {
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
	if err = client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}

	log.Println("Processed", *messageCountStart, "messages")
}
