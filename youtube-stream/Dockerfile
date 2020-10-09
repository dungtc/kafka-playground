FROM confluentinc/cp-kafka-connect:5.3.1
# Create plugin directory
RUN mkdir -p /usr/share/java/plugins

COPY connector-elastic-sink.json /tmp

# datagen config
RUN confluent-hub install --no-prompt confluentinc/kafka-connect-elasticsearch:10.0.1