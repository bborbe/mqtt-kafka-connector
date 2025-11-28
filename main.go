// Copyright (c) 2019 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/bborbe/argument/v2"
	"github.com/bborbe/errors"
	libhttp "github.com/bborbe/http"
	"github.com/bborbe/run"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	runtime.GOMAXPROCS(runtime.NumCPU())
	_ = flag.Set("logtostderr", "true")

	ctx := context.Background()
	app := &application{}
	if err := argument.Parse(ctx, app); err != nil {
		glog.Exitf("parse app failed: %v", err)
	}

	glog.V(0).Infof("application started")
	if err := app.run(
		contextWithSig(ctx),
	); err != nil {
		glog.Exitf("application failed: %+v", err)
	}
	glog.V(0).Infof("application finished")
	os.Exit(0)
}

type application struct {
	InitialDelay time.Duration `required:"false" arg:"initial-delay" env:"INITIAL_DELAY" usage:"initial time before processing starts" default:"1m"`
	MqttBroker   string        `required:"true"  arg:"mqtt-broker"   env:"MQTT_BROKER"   usage:"broker address to connect"`
	MqttUsername string        `required:"false" arg:"mqtt-user"     env:"MQTT_USER"     usage:"mqtt user"`
	MqttPassword string        `required:"false" arg:"mqtt-password" env:"MQTT_PASSWORD" usage:"mqtt password"                                        display:"length"`
	MqttTopic    string        `required:"true"  arg:"mqtt-topic"    env:"MQTT_TOPIC"    usage:"topic name dummy data are written to"`
	KafkaBrokers string        `required:"true"  arg:"kafka-brokers" env:"KAFKA_BROKERS" usage:"kafka brokers"`
	KafkaTopic   string        `required:"true"  arg:"kafka-topic"   env:"KAFKA_TOPIC"   usage:"kafka topic"`
	Port         int           `required:"false" arg:"port"          env:"PORT"          usage:"port to listen"                        default:"9022"`
}

func contextWithSig(ctx context.Context) context.Context {
	ctxWithCancel, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()

		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		select {
		case <-signalCh:
		case <-ctx.Done():
		}
	}()

	return ctxWithCancel
}

func (a *application) run(
	ctx context.Context,
) error {
	return run.CancelOnFirstFinish(
		ctx,
		run.Delayed(a.createFetcherCron(), a.InitialDelay),
		a.runHTTPServer,
	)
}

func (a *application) createFetcherCron() func(ctx context.Context) error {
	return func(ctx context.Context) error {
		config := sarama.NewConfig()
		config.Version = sarama.V2_0_0_0
		config.Producer.RequiredAcks = sarama.WaitForAll
		config.Producer.Retry.Max = 10
		config.Producer.Return.Successes = true

		client, err := sarama.NewClient(strings.Split(a.KafkaBrokers, ","), config)
		if err != nil {
			return errors.Wrap(ctx, err, "create client failed")
		}
		defer client.Close()

		producer, err := sarama.NewSyncProducerFromClient(client)
		if err != nil {
			return errors.Wrap(ctx, err, "create sync producer failed")
		}
		defer producer.Close()

		mqttClient := mqtt.NewClient(
			mqtt.NewClientOptions().
				AddBroker(a.MqttBroker).
				SetUsername(a.MqttUsername).
				SetPassword(a.MqttPassword),
		)
		if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
			return errors.Wrap(ctx, token.Error(), "connect failed")
		}

		errs := make(chan error, runtime.NumCPU())
		defer close(errs)

		if token := mqttClient.Subscribe(a.MqttTopic, 0, func(mqttClient mqtt.Client, mqttMessage mqtt.Message) {
			partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
				Topic: a.KafkaTopic,
				Value: sarama.ByteEncoder(mqttMessage.Payload()),
			})
			if err != nil {
				select {
				case <-ctx.Done():
					return
				case errs <- err:
					return
				}
			}
			glog.V(2).Infof("send message successful to %s with partition %d offset %d", a.KafkaTopic, partition, offset)
		}); token.Wait() &&
			token.Error() != nil {
			return errors.Wrap(ctx, token.Error(), "subscribe failed")
		}
		return <-errs
	}
}

func (a *application) runHTTPServer(ctx context.Context) error {
	router := mux.NewRouter()
	router.Path("/healthz").Handler(libhttp.NewPrintHandler("OK"))
	router.Path("/readiness").Handler(libhttp.NewPrintHandler("OK"))
	router.Path("/metrics").Handler(promhttp.Handler())

	glog.V(2).Infof("starting http server listen on :%d", a.Port)
	return libhttp.NewServer(
		fmt.Sprintf(":%d", a.Port),
		router,
	).Run(ctx)
}
