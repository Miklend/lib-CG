package fabricClient

import (
	"github.com/Miklend/lib-CG/common/logging"

	"github.com/Miklend/lib-CG/client/broker"
	"github.com/Miklend/lib-CG/client/broker/kafka"
	"github.com/Miklend/lib-CG/models/configs"
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
