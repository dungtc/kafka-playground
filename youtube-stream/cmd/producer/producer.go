package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/dungtc/kafka-playground/youtube-stream/repository"
	"github.com/dungtc/kafka-playground/youtube-stream/setting"

	"github.com/Shopify/sarama"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	brokerList = kingpin.Flag("brokerList", "List of brokers").Default("localhost:9092").Strings()
	topic      = kingpin.Flag("topic", "Topic name").Default("youtube").String()
	maxRetry   = kingpin.Flag("maxRetry", "Retry limit").Default("5").Int()
	Conf       setting.Configuration
)

func init() {
	// Load configuration
	err := setting.LoadConfiguration("env", "config", "toml", &Conf)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func main() {
	kingpin.Parse()
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Version = sarama.V2_5_0_0

	// safe producer
	config.Producer.Idempotent = true
	config.Net.MaxOpenRequests = 1 // max.in.flight.request.per.connection
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = *maxRetry
	config.Producer.Retry.Backoff = 20 * time.Second
	config.Producer.Timeout = 10 * time.Second
	config.Producer.Flush.MaxMessages = 5 // max message in a request

	// hight throughtput
	config.Producer.Compression = sarama.CompressionSnappy // compression algorithm
	// config.Producer.MaxMessageBytes = 15000                // max size each message
	config.Producer.Flush.Bytes = 64 * 1024                // batch size ,default 16 KB
	config.Producer.Flush.Frequency = 5 * time.Millisecond // linger.ms

	producer, err := sarama.NewSyncProducer(*brokerList, config)
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Panic(err)
		}
	}()

	musicCategory := "10"

	youtubeRepository := repository.NewYoutubeRepository(Conf.GoogleConnection)

	videos, err := youtubeRepository.GetVideos(musicCategory)
	if err != nil {
		log.Panic(err)
	}

	for i := 0; i < len(videos.Items); i++ {
		fmt.Printf("video %+v\n", videos.Items[i])

		val, err := json.Marshal(videos.Items[i])
		if err != nil {
			log.Panic(err)
		}

		partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
			Topic: *topic,
			Key:   sarama.StringEncoder(videos.Items[i].Id),
			Value: sarama.StringEncoder(string(val)),
		})
		if err != nil {
			log.Panic(err)
		}
		log.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", *topic, partition, offset)
	}
}
