[![Go Report Card](https://goreportcard.com/badge/github.com/dungtc/kafka-playground)](https://goreportcard.com/report/github.com/dungtc/kafka-playground)
[![License: MIT](https://img.shields.io/github/license/dungtc/kafka-playground)](https://opensource.org/licenses/Apache-2.0)

### Author
- [Linkedin profile](https://www.linkedin.com/in/tran-dung-076221179/)
Contact me if you have any questions or joining contributors

### Contributing
- Change repository code and make a commit. 
- Create a pull request on your feature branch

# 1. Messaging with Simple Publish-Subscribe (pub/sub)
A publisher or multiple publishers write messages to many partitions and single or multiple consumer groups consume those messages.

**Consumer group** is a group of consumers, it can parallelize the data consumption.

When a new consumer joins to consumer group or removes a consumer, **rebalancing** happen

**Consumer group coordinator** implements the rebalancing strategy and manages the state of the group. It auto rebalance message to the appropriate consumer in each consumer group

![alt text](https://images.squarespace-cdn.com/content/v1/56894e581c1210fead06f878/1512813188413-EEI12VI1FMQLJ4XMRQTJ/ke17ZwdGBToddI8pDm48kJMyHperWZxre3bsQoFNoPhZw-zPPgdn4jUwVcJE1ZvWEtT5uBSRWt4vQZAgTJucoTqqXjS3CfNDSuuf31e0tVGDclntk9GVn4cF1XFdv7wlNvED_LyEM5kIdmOo2jMRZpu3E9Ef3XsXP1C_826c-iU/KafkaPubSub.png?format=500w)

**Producer cmd**

```console
$ go run ./cmd/producer/producer.go --topic topic-name --partitions 1 --replicas 1
```

**You have to config Kafka cluster when change number partitions and replicas > 1**


In case, you want to modify topic configuration ***manually*** after setting up Kafka cluster, you can use Kafka-CLI to execute:

Add new topic
```console
$ docker-compose exec broker kafka-topics --create --topic youtube --partitions 3 -replication-factor 3 --bootstrap-server broker:9092
```

Alter exist topic
```console
$ docker-compose exec broker kafka-topics --alter --topic youtube --partitions 3 -replication-factor 3 --bootstrap-server broker:9092
```

Kafka doesn't support reducing the number of partitions for a topic 

**Consumer cmd**

```console
$ go run ./cmd/consumer/consumer.go --topic topic-name
```

# 2. Youtube streaming
A simple streaming model to collect videos data and sync data between **POSTGRESQL** and **ELASTICSEARCH**

![alt text](https://github.com/dungtc/kafka-playground/blob/develop/youtube-stream/diagram/youtube.png?raw=true)

**Prepare connector image**
The Kafka connect base image only contain a few connectors. To add a new connector, you need to build a new docker image that has new connectors installed.

***Docker file***
```dockerfile
FROM confluentinc/cp-kafka-connect-base:6.0.0
RUN confluent-hub install --no-prompt confluentinc/kafka-connect-elasticsearch:10.0.1
```

***Build***

```console
$ docker build . my-connect-image:1.0.0
```

***Prepare connector config***
```json
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
```console
$ ./run.sh
```
Detail at https://docs.confluent.io/current/connect/managing/extending.html

***Create KSQL stream***

```sql
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
```console
$ curl -X GET "localhost:9200/_search?pretty&size=100" -H 'Content-Type: application/json' -d'
{
    "query": {
        "match_all": {}
    }
}
'
```
# 3. Kafka cluster

When a Kafka client connects to one of the brokers, it will auto-connect to entire Kafka cluster brokers. Kafka broker call **bootstrap server**.

**Zookeeper** responsible for managing Kafka brokers, it performs partition leader election and in-sync-replicas.
Furtherly it informs **Kafka** if any configurations change happen (such as topics, partitions...).

***In production environment, you need to config **Zookeeper Cluster** too***

![alt text](https://github.com/dungtc/kafka-playground/blob/develop/kafka-cluster/Kafka-cluster.png?raw=true)

**Environment**
```dockerfile
KAFKA_ZOOKEEPER_CONNECT: 'zookeeper:2181'
KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://broker2:39092,PLAINTEXT_HOST://localhost:9093
```

You should notice to these env variables, it defines internal network and external network publish to client

# 4. Schema Registry

Producers send bytes to the Kafka cluster, the consumers don't know any data that the producer sends, it has no verification and will break the consumer. In a distributed system it becomes very hard for evolving microservices and breaks consistency between other teams. That's why the schema registry was created.

![alt text](https://github.com/dungtc/kafka-playground/blob/develop/schema/kafka-schema-registry.png?raw=true)


**Schema registry** provides a restful interface for storing and retrieving Avro, JSON, Protobuf schema between producer and consumer


### [Schema registry API usage](https://docs.confluent.io/current/schema-registry/develop/using.html)

List all subjects
```console
$ curl -X GET http://localhost:8081/subjects
```
Example result:
```
["Kafka-value","Kafka-key"]
```

Get schema registered version of subject
```console
$ curl -X GET http://localhost:8081/subjects/${subject-name}/versions/1
```
Delete schema version of the subject

```console
$ curl -X DELETE http://localhost:8081/subjects/${subject-name}/versions/1
```

### [Avro schema](https://avro.apache.org/docs/current/spec.html)

Avro schema form includes 4 fields: type, namespace, name, fields
```json
{
     "type": "record",
     "namespace": "com.example",
     "name": "FullName",
     "fields": [
       { "name": "first", "type": "string" },
     ]
} 
```
**Primitive Data Types** int, long, float, double, string, null, bytes, boolean

**Complex Data Types** record, enum, array, map, unions, fixed (fixed-sized field)

**"record"** represents an object 
```json
{
     "type": "record",
     "namespace": "com.example",
     "name": "FullName",
     "fields": [
       { "name": "first", "type": "string" },
     ]
} 
```

**"enum"** defines enumerated values

```json
{"type": "enum", "symbols": ["A", "B","C"]}
```

**"array"** defines an array fields with a type
```json
{"type": "array", "items": "string"}
```

**"map"** defines key, value pairs. The key must be a string
```json
{"type": "map", "values": "string"}
```

**"unions"** define a list of variable types
```json
{"type": ["int", "string", "null"]}
```

**Logical Data Types** supported in Avro version 1.7.7
Provides more meaning with existing primitive types
- decimal (bytes)
- date (int)
- time-millis (long) number of milliseconds after midnight
- timestamp-millis (long) number of milliseconds from Unix epoch

```json
{"type": "int", "logicalType": "decimal"}
```

# 5. [Connect](https://docs.confluent.io/platform/current/connect/index.html)

Kafka Connect is used for streaming data between Kafka and other data systems. It's simple, low latency, reliable, and scalable. Kafka Connect can push large data sets in and out of Kafka. It becomes an integral component of an ETL pipeline when combined with Kafka and streaming processing framework.

![Connector architecture](./connector/kafka-connector.png)

List of available connectors: https://www.confluent.io/hub
Kafka Connect with schema registry: https://docs.confluent.io/platform/current/schema-registry/connect.html

### [Connect API Usage](https://docs.confluent.io/platform/current/connect/references/restapi.html)

Create a new connector

```console
$ curl -X POST http://localhost:8083/connectors -H "Content-Type: application/json" -d @source/distributed/connector-json.json
```

Get connector information
```console
$ curl -X GET http://localhost:8083/connectors/${connector-name}
```

There are 2 types of worker configuration: **standalone mode** and **distributed mode**

**Standalone Mode** is a single process that runs connectors and tasks. For development purposes only.

**Distributed Mode** is multiple workers that run connectors and tasks, it's scalable and fault-tolerant. If a worker dies, rebalance happens, tasks will be delivered to other workers. Use for **production**.

### Source connector

The source connector gives you the way to import data from any data sets and write to Kafka topic

Create a new source connector
```console
$ curl -X POST http://localhost:8083/connectors -H "Content-Type: application/json" -d @source/distributed/connector-json.json
```

**Example JDBC source connector**
```yaml
connector.class=io.confluent.connect.jdbc.JdbcSourceConnector
timestamp.column.name=[column_name1]
incrementing.column.name=[auto_increment_column]
table.whitelist=orders
mode=timestamp
topic.prefix=postgresql-source-connector-
connection.url=jdbc:postgresql://postgres:5432/postgres
connection.user=postgres
connection.password=password
key.converter.schemas.enable=false
value.converter.schemas.enable=false
value.converter=org.apache.kafka.connect.json.JsonConverter
key.converter=org.apache.kafka.connect.json.JsonConverter
```

### Sink connector

The Kafka sink connector allows you to export data from Kafka to any data sets

Create a new sink connector
```console
$ curl -X POST http://localhost:8083/connectors -H "Content-Type: application/json" -d @sink/distributed/elastic-json.json
```

**Example Elastic sink connector with json schema less**
```yaml
connector.class=io.confluent.connect.elasticsearch.ElasticsearchSinkConnector
type.name=_doc
errors.retry.timeout=-1
topics=postgresql-source-connector-orders
transforms=ExtractTimestamp
key.ignore=true
schema.ignore=true
transforms.ExtractTimestamp.type=org.apache.kafka.connect.transforms.InsertField$Value
value.converter.schemas.enable=false
value.converter=org.apache.kafka.connect.json.JsonConverter
connection.url=http://elasticsearch:9200
key.converter=org.apache.kafka.connect.storage.StringConverter
transforms.ExtractTimestamp.timestamp.field=timestamp
```

**Query Elastic data**
```console
$ curl -X GET "localhost:9200/_search?pretty&size=100" -H 'Content-Type: application/json' -d'{"query":{"match_all":{}}}'
```