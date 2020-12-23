#!/bin/bash

# Standalone mode for develop purpose
### define a property file and use connect process in command line 
docker run --rm -it -v $(pwd)/source/standalone:/connectors --net=host landoop/fast-data-dev:2.5.0 bash
cd connectors
connect-standalone worker.properties source-connector.properties

# using rest API
curl -X POST http://localhost:8083/connectors -H "Content-Type: application/json" -d @source/standalone/standalone.json