package fabricClient

import (
	"lib-CG/client/broker"
	"lib-CG/client/broker/kafka"
	"lib-CG/common/logging"
	"lib-CG/models/configs"
)

const (
	kafkaBrokerType = "kafka"
	mock            = "mock"
)

func NewBroker(cfg configs.Broker, logger *logging.Logger) broker.BrokerClient {
	switch cfg.BrockerType {
	case kafkaBrokerType:
		return kafka.NewKafkaBroker(cfg, logger)
	case mock:
		return broker.NewMockBrokerClient()
	default:
		return nil
	}
}
