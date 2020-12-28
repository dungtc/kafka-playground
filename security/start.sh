#!/bin/bash

# Producer
docker-compose exec kafka-cluster kafka-console-producer --broker-list localhost:19092 --topic kafka-security --producer.config /etc/kafka/secrets/host.producer.ssl.config


# Consumer
docker-compose exec kafka-cluster kafka-console-consumer --broker-list localhost:19092 --topic kafka-security --consumer.config /etc/kafka/secrets/host.consumer.ssl.config
