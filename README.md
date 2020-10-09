# 1. Messaging with Simple Publish Subscribe (pub/sub)
A publisher or multiple publisher write messages to many partitions and single or multiple consumer groups consume those messages.

**Consumer group** is a group of consumers, it can parallelise the data consumption.

When a new consumer joins to consumer group or remove a consumer, **rebalancing** happen

**Consumer group coordinator** implement rebalacing strategy and manage the state of the group. It auto rebalace message to appropriate consumer in each consumer group

![alt text](https://images.squarespace-cdn.com/content/v1/56894e581c1210fead06f878/1512813188413-EEI12VI1FMQLJ4XMRQTJ/ke17ZwdGBToddI8pDm48kJMyHperWZxre3bsQoFNoPhZw-zPPgdn4jUwVcJE1ZvWEtT5uBSRWt4vQZAgTJucoTqqXjS3CfNDSuuf31e0tVGDclntk9GVn4cF1XFdv7wlNvED_LyEM5kIdmOo2jMRZpu3E9Ef3XsXP1C_826c-iU/KafkaPubSub.png?format=500w)

# 2. Youtube streaming
A simple streaming model to collect videos data and sync data between **POSTGRESQL** and **ELASTICSEARCH**

![alt text](https://github.com/dungtc/kafka-playground/blob/develop/youtube-stream/diagram/youtube.png?raw=true)

**Prepare connector image**
The kafka connect base image only contain a few connectors. To add new connector, you need to build a new docker image that have new connectors installed.

***Docker file***
```
FROM confluentinc/cp-kafka-connect-base:6.0.0
RUN confluent-hub install --no-prompt confluentinc/kafka-connect-elasticsearch:10.0.1
```

***Build***

```
docker build . my-connect-image:1.0.0
```

***Prepare connector config***
```
{
  "name": "elastic-sink4",
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
}
```

***Insert connector configuration***
```
./run.sh
```
Detail at https://docs.confluent.io/current/connect/managing/extending.html

***Create ksql stream***

```
CREATE STREAM youtube (
    kind VARCHAR,
    etag VARCHAR,
    id VARCHAR,
	snippet MAP<VARCHAR,VARCHAR>
  ) WITH (
    KAFKA_TOPIC='youtube2',
    VALUE_FORMAT='JSON'
  );
```

***Query elasticsearch***
```
curl -X GET "localhost:9200/_search?pretty&size=100" -H 'Content-Type: application/json' -d'
{
    "query": {
        "match_all": {}
    }
}
'
```
