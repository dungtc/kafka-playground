#!/bin/bash

# # Distributed mode for production
# avro schema
curl -X POST http://localhost:8083/connectors -H "Content-Type: application/json" -d @source/distributed/connector-avro.json

# json without schema
curl -X POST http://localhost:8083/connectors -H "Content-Type: application/json" -d @source/distributed/connector-json.json

# protobuf
curl -X POST http://localhost:8083/connectors -H "Content-Type: application/json" -d @source/distributed/connector-protobuf.json