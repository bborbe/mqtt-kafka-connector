// Copyright (c) 2019 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/bborbe/argument/v2"
	"github.com/bborbe/errors"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/glog"
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
	MqttBroker   string `required:"true"  arg:"mqtt-broker"   env:"MQTT_BROKER"   usage:"broker address to connect"`
	MqttUsername string `required:"false" arg:"mqtt-user"     env:"MQTT_USER"     usage:"mqtt user"`
	MqttPassword string `required:"false" arg:"mqtt-password" env:"MQTT_PASSWORD" usage:"mqtt password"                        display:"length"`
	MqttTopic    string `required:"true"  arg:"mqtt-topic"    env:"MQTT_TOPIC"    usage:"topic name dummy data are written to"`
	MqttPayload  string `required:"true"  arg:"mqtt-payload"  env:"MQTT_PAYLOAD"  usage:"content written to topic"`
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

func (a *application) run(ctx context.Context) error {
	mqttClient := mqtt.NewClient(
		mqtt.NewClientOptions().
			AddBroker(a.MqttBroker).
			SetUsername(a.MqttUsername).
			SetPassword(a.MqttPassword),
	)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return errors.Wrap(ctx, token.Error(), "connect failed")
	}
	if token := mqttClient.Publish(
		a.MqttTopic,
		0,
		false,
		a.MqttPayload,
	); token.Wait() && token.Error() != nil {
		return errors.Wrap(ctx, token.Error(), "publish failed")
	}
	glog.V(2).Infof("message published successful")
	return nil
}
