#!/bin/sh

# Producer avro schema
docker run -it --rm --net=host confluentinc/cp-schema-registry:5.5.1 bash
kafka-avro-console-producer \
         --broker-list localhost:9092 --topic schema-1 --property schema.registry.url=http://localhost:8081 --property value.schema='{"type":"record","name":"myrecord","fields":[{"name":"f1","type":"string"}]}'

# Data example
# {"f1":"value1"}

# Consumer avro schema
kafka-avro-console-consumer --bootstrap-server localhost:9092 --topic schema-1 \
 --property schema.registry.url=http://localhost:8081 --from-beginning