package fabricClient

import (
	"fmt"

	"github.com/Miklend/lib-CG/client/node"
	"github.com/Miklend/lib-CG/client/node/alchemy"
	"github.com/Miklend/lib-CG/common/logging"
	"github.com/Miklend/lib-CG/models/configs"
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
