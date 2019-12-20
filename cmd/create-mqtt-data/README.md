# Mqtt Data Creator

Insert data into a topic.

## Run

```bash
go run main.go \
-mqtt-broker=tcp://mqtt:1883 \
-mqtt-topic=mytopic/test \
-mqtt-payload=banana \
-v=2
```
