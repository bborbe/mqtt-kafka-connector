# MQTT Kafka Connector

Copy all data from MQTT to Kafka

## Start MQTT + Kafka

echo "127.0.0.1 kafka" >> /etc/hosts
echo "127.0.0.1 mqtt" >> /etc/hosts
echo "127.0.0.1 ksql-server" >> /etc/hosts
echo "127.0.0.1 schema-registry" >> /etc/hosts
echo "127.0.0.1 zookeeper" >> /etc/hosts

```bash
make run
```
## Run

```bash
go run main.go \
-initial-delay=1s \
-mqtt-broker=tcp://mqtt:1883 \
-mqtt-topic=mytopic/test \
-kafka-brokers=kafka:9092 \
-kafka-topic=mqtt \
-v=2
```
