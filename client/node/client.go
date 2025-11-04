package node

import (
	"context"

	"github.com/Miklend/lib-CG/models/dto"
	"github.com/ethereum/go-ethereum/rpc"
)

type Provider interface {
	BlockByNumber(ctx context.Context, param string) (dto.BlockDTO, error)
	TxByHash(ctx context.Context, param string) (dto.TxDTO, error)
	ReceiptByTxHash(ctx context.Context, param string) (dto.ReceiptDTO, error)
	ReceiptByBlockNumber(ctx context.Context, param string) ([]dto.ReceiptDTO, error)

	BatchRequest(ctx context.Context, requests []rpc.BatchElem) ([]rpc.BatchElem, error)
	BatchBlockByNumber(ctx context.Context, numbers []string) (map[dto.Hash]dto.BlockDTO, error)
	BatchReceiptByBlockNumber(ctx context.Context, numbers []string) (map[dto.Hash][]dto.ReceiptDTO, error)
	BatchReceiptByTxHash(ctx context.Context, hashes []string) (map[dto.Hash]dto.ReceiptDTO, error)
	BatchBlockWithReceiptByNumber(ctx context.Context, numbers []string) (map[dto.Hash]dto.BlockDTO, error)

	SubscribeBlockWithReceipts(ctx context.Context, blockCh chan<- *dto.BlockDTO) (*rpc.ClientSubscription, error)
	SubscribePendingTransactions(ctx context.Context, txsCh chan<- *dto.TxDTO) (*rpc.ClientSubscription, error)

	Close() error
}
