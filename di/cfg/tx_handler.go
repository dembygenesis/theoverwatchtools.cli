package cfg

import (
	"github.com/dembygenesis/local.tools/internal/config"
	"github.com/dembygenesis/local.tools/internal/persistence"
	"github.com/sarulabs/dingo/v4"
)

const (
	txHandlerLayer = "tx_handler_layer"
)

func getTxHandlerLayer() []dingo.Def {
	return []dingo.Def{
		{
			Name: txHandlerLayer,
			Build: func(
				cfg *config.Config,
			) (persistence.TransactionHandler, error) {
				// hello

				// Spawn the connection here
				return nil, nil
			},
		},
	}
}
