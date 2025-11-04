package fabricClient

import (
	"fmt"
	"lib-CG/client/node"
	"lib-CG/client/node/alchemy"
	"lib-CG/common/logging"
	"lib-CG/models/configs"
)

const (
	alchemyType = "alchemy"
)

func NewProvider(cfg configs.Provider, logger *logging.Logger) (node.Provider, error) {
	switch cfg.ProviderType {
	case alchemyType:
		return alchemy.NewAlchemyClient(cfg, logger)
	default:
		return nil, fmt.Errorf("not found provider for client type: %s", cfg.ProviderType)
	}
}
