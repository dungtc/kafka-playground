package simple

import (
	"log"

	"github.com/Shopify/sarama"
)

type Admin struct {
	sarama.ClusterAdmin
}

func NewClusterAdmin(ver string, brokerList ...string) (*Admin, error) {
	config := sarama.NewConfig()
	version, err := sarama.ParseKafkaVersion(ver)
	if err != nil {
		log.Panicf("Error parsing Kafka version: %v", version)
	}
	config.Version = version
	admin, err := sarama.NewClusterAdmin(brokerList, config)
	if err != nil {
		log.Fatal("Error while creating cluster admin: ", err.Error())
	}
	// defer func() { _ = admin.Close() }()
	return &Admin{
		admin,
	}, err
}

func (a *Admin) NewTopic(topic string, numPartitions int32, replicationFactor int16) error {
	isExist, err := a.FindTopic(topic)
	if err != nil {
		return err
	}
	if isExist {
		return nil
	}

	err = a.ClusterAdmin.CreateTopic(topic, &sarama.TopicDetail{
		NumPartitions:     numPartitions,
		ReplicationFactor: replicationFactor,
	}, false)
	if err != nil {
		log.Fatal("Error while creating topic: ", err.Error())
	}
	return err
}

func (a *Admin) FindTopic(topic string) (bool, error) {
	topics, err := a.ClusterAdmin.ListTopics()
	if err != nil {
		log.Fatal("Error while creating topic: ", err.Error())
	}
	for k := range topics {
		if k == topic {
			return true, nil
		}
		// fmt.Println(topics[k])
	}

	return false, nil
}
