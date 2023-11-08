// Copyright (c) 2019 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/bborbe/argument"
	"github.com/bborbe/run"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	runtime.GOMAXPROCS(runtime.NumCPU())
	_ = flag.Set("logtostderr", "true")

	app := &application{}
	if err := argument.Parse(app); err != nil {
		glog.Exitf("parse app failed: %v", err)
	}

	glog.V(0).Infof("application started")
	if err := app.run(
		contextWithSig(context.Background()),
	); err != nil {
		glog.Exitf("application failed: %+v", err)
	}
	glog.V(0).Infof("application finished")
	os.Exit(0)
}

type application struct {
	InitialDelay time.Duration `required:"false" arg:"initial-delay" env:"INITIAL_DELAY" usage:"initial time before processing starts" default:"1m"`
	MqttBroker   string        `required:"true" arg:"mqtt-broker" env:"MQTT_BROKER" usage:"broker address to connect"`
	MqttUsername string        `required:"false" arg:"mqtt-user" env:"MQTT_USER" usage:"mqtt user"`
	MqttPassword string        `required:"false" arg:"mqtt-password" env:"MQTT_PASSWORD" usage:"mqtt password" display:"length"`
	MqttTopic    string        `required:"true" arg:"mqtt-topic" env:"MQTT_TOPIC" usage:"topic name dummy data are written to"`
	KafkaBrokers string        `required:"true" arg:"kafka-brokers" env:"KAFKA_BROKERS" usage:"kafka brokers"`
	KafkaTopic   string        `required:"true" arg:"kafka-topic" env:"KAFKA_TOPIC" usage:"kafka topic"`
	Port         int           `required:"false" arg:"port" env:"PORT" usage:"port to listen" default:"9022"`
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
		a.runHttpServer,
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
			return errors.Wrap(err, "create client failed")
		}
		defer client.Close()

		producer, err := sarama.NewSyncProducerFromClient(client)
		if err != nil {
			return errors.Wrap(err, "create sync producer failed")
		}
		defer producer.Close()

		mqttClient := mqtt.NewClient(
			mqtt.NewClientOptions().
				AddBroker(a.MqttBroker).
				SetUsername(a.MqttUsername).
				SetPassword(a.MqttPassword),
		)
		if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
			return errors.Wrap(token.Error(), "connect failed")
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
		}); token.Wait() && token.Error() != nil {
			return errors.Wrap(token.Error(), "subscribe failed")
		}
		return <-errs
	}
}

func (a *application) runHttpServer(ctx context.Context) error {
	router := mux.NewRouter()
	router.HandleFunc("/healthz", a.check)
	router.HandleFunc("/readiness", a.check)
	router.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.Port),
		Handler: router,
	}

	go func() {
		select {
		case <-ctx.Done():
			if err := server.Shutdown(ctx); err != nil {
				glog.Warningf("shutdown failed: %v", err)
			}
		}
	}()
	err := server.ListenAndServe()
	if err == http.ErrServerClosed {
		glog.V(0).Info(err)
		return nil
	}
	return errors.Wrap(err, "httpServer failed")
}

func (a *application) check(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprint(resp, "OK")
}
