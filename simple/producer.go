package simple

import (
	"log"
	"time"

	"github.com/Shopify/sarama"
)

func NewProducer(maxRetry int, brokerList ...string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Version = sarama.V2_5_0_0

	// safe producer
	config.Producer.Idempotent = true
	config.Net.MaxOpenRequests = 1 // max.in.flight.request.per.connection
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = maxRetry
	config.Producer.Retry.Backoff = 20 * time.Second
	config.Producer.Timeout = 10 * time.Second
	config.Producer.Flush.MaxMessages = 5 // max message in a request

	// hight throughtput
	config.Producer.Compression = sarama.CompressionSnappy // compression algorithm
	// config.Producer.MaxMessageBytes = 15000                // max size each message
	config.Producer.Flush.Bytes = 64 * 1024                // batch size ,default 16 KB
	config.Producer.Flush.Frequency = 5 * time.Millisecond // linger.ms

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Panic(err)
		return nil, err
	}
	return producer, nil
}
