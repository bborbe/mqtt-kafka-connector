# MQTT Kafka Connector

Subscribes a MQTT topic and sends each message to a Kafka topic.

## Start MQTT + Kafka

```bash
make run
```       

```bash
kafka-console-consumer -topic=mqtt -brokers=localhost:9092
```

## Run

```bash
go run main.go \
-initial-delay=1s \
-mqtt-broker=tcp://localhost:1883 \
-mqtt-topic=mytopic/test \
-kafka-brokers=localhost:9092 \
-kafka-topic=mqtt \
-v=2
```
