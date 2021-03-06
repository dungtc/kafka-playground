#!/bin/sh
  
docker-compose up -d --build

# Creating the user kafka
# kafka is configured as a super user, no need for additional ACL
docker-compose exec kafka-1 kafka-configs --zookeeper zookeeper:2182 --alter --add-config 'SCRAM-SHA-256=[password=kafka-pass],SCRAM-SHA-512=[password=kafka-pass]' --entity-type users --entity-name kafka
docker-compose exec kafka-1 kafka-configs --zookeeper zookeeper:2182 --alter --add-config 'SCRAM-SHA-256=[password=admin-pass],SCRAM-SHA-512=[password=admin-pass]' --entity-type users --entity-name admin
docker-compose exec kafka-1 kafka-configs --zookeeper zookeeper:2182 --alter --add-config 'SCRAM-SHA-256=[password=producer-pass],SCRAM-SHA-512=[password=producer-pass]' --entity-type users --entity-name producer
docker-compose exec kafka-1 kafka-configs --zookeeper zookeeper:2182 --alter --add-config 'SCRAM-SHA-256=[password=consumer-pass],SCRAM-SHA-512=[password=consumer-pass]' --entity-type users --entity-name consumer

# ACLs
docker-compose exec kafka-1 kafka-acls  --authorizer-properties zookeeper.connect=zookeeper:2182 --add --allow-principal User:producer --producer --topic=*
docker-compose exec kafka-1 kafka-acls  --authorizer-properties zookeeper.connect=zookeeper:2182 --add --allow-principal User:consumer --consumer --topic=* --group=*

echo "Example configuration:"
echo "-> kafka-console-producer --broker-list localhost:9093 --producer.config /tmp/producer.conf --topic test"
echo "-> kafka-console-consumer --bootstrap-server localhost:9094 --consumer.config /tmp/consumer.conf --topic test --from-beginning"
echo "ZooKeeper shell with authorization from host:"
echo "-> KAFKA_OPTS=\"-Djava.security.auth.login.config=zookeeper.sasl.jaas.conf\" zookeeper-shell localhost:2182"
echo "ZooKeeper shell with authorization within container (KAFKA_OPTS already set):"
echo "-> docker-compose exec kafka-1 zookeeper-shell zookeeper:2182"
echo "Kafkacat with authorization from host:"
echo "-> kafkacat -L -b localhost:9094 -F kafka/kafkacat.conf"