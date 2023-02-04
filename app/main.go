package main

import (
	"fmt"

	_ "go.uber.org/cadence/.gen/go/cadence"
	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/cadence/worker"

	"github.com/uber-go/tally"
	"go.uber.org/yarpc"
	_ "go.uber.org/yarpc/api/transport"
	"go.uber.org/yarpc/transport/tchannel"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// cadenceService should always be cadence-frontend
	CadenceService = "cadence-frontend"
	ClientName     = "greetings-worker"
	Domain         = "tavern"
	Host           = "127.0.0.1:7933"
	TaskList       = "greetings"
)

func main() {
	worker, logger, err := newWorkerServiceClient()

	if err != nil {
		panic(err)
	}

	if err := worker.Start(); err != nil {
		panic(fmt.Errorf("failed to start the worker: %v", err))
	}

	logger.Info("Started Worker.", zap.String("worker", TaskList))

	select {}
}

func newWorkerServiceClient() (worker.Worker, *zap.Logger, error) {
	logger, err := newLogger()
	if err != nil {
		return nil, nil, err
	}

	worketOptions := worker.Options{
		Logger:       logger,
		MetricsScope: tally.NewTestScope(TaskList, map[string]string{}),
	}

	connection, err := newCadenceConnection(ClientName)
	if err != nil {
		return nil, nil, err
	}

	return worker.New(connection, Domain, TaskList, worketOptions), logger, nil

}

func newCadenceConnection(clientName string) (workflowserviceclient.Interface, error) {
	ch, err := tchannel.NewChannelTransport(tchannel.ServiceName(ClientName))
	if err != nil {
		return nil, fmt.Errorf("failed to set up Transport channel: %v", err)
	}

	dispatcher := yarpc.NewDispatcher(yarpc.Config{
		Name: ClientName,
		Outbounds: yarpc.Outbounds{
			CadenceService: {Unary: ch.NewSingleOutbound(Host)},
		},
	})

	if err := dispatcher.Start(); err != nil {
		return nil, fmt.Errorf("failed to start dispatcher: %v", err)
	}

	return workflowserviceclient.New(dispatcher.ClientConfig(CadenceService)), nil

}

func newLogger() (*zap.Logger, error) {
	config := zap.NewDevelopmentConfig()

	config.Level.SetLevel(zapcore.InfoLevel)

	var err error
	logger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %v", err)
	}

	return logger, nil
}
