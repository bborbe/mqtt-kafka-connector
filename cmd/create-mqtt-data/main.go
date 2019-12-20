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

	"github.com/bborbe/argument"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/glog"
	"github.com/pkg/errors"
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
	MqttBroker  string `required:"true" arg:"mqtt-broker" env:"MQTT_BROKER" usage:"broker address to connect"`
	MqttTopic   string `required:"true" arg:"mqtt-topic" env:"MQTT_TOPIC" usage:"topic name dummy data are written to"`
	MqttPayload string `required:"true" arg:"mqtt-payload" env:"MQTT_PAYLOAD" usage:"content written to topic"`
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
	client := mqtt.NewClient(
		mqtt.NewClientOptions().AddBroker(a.MqttBroker),
	)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return errors.Wrap(token.Error(), "connect failed")
	}
	if token := client.Publish(
		a.MqttTopic,
		0,
		false,
		a.MqttPayload,
	); token.Wait() && token.Error() != nil {
		return errors.Wrap(token.Error(), "publish failed")
	}
	glog.V(2).Infof("message published successful")
	return nil
}
