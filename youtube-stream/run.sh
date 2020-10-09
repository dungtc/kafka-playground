#!/bin/bash
echo -e "\nAdding ElastichSearch Kafka Sink Connector for the 'youtube' topic into the 'youtube' collection"

curl -X POST -H "Content-Type:application/json" --data '{
  "name": "elastic-sink5",
  "config": {
    "connector.class": "io.confluent.connect.elasticsearch.ElasticsearchSinkConnector",
    "key.converter": "org.apache.kafka.connect.storage.StringConverter",
    "value.converter": "org.apache.kafka.connect.json.JsonConverter",
    "value.converter.schemas.enable": "false",
    "errors.retry.timeout": "-1",
    "topics": "youtube2",
    "connection.url": "http://elasticsearch:9200",
    "type.name": "_doc",
    "key.ignore": "false",
    "schema.ignore": "false",
    "transforms": "ExtractTimestamp",
    "transforms.ExtractTimestamp.type": "org.apache.kafka.connect.transforms.InsertField$Value",
    "transforms.ExtractTimestamp.timestamp.field": "timestamp"
  }
}' http://localhost:8083/connectors