package alchemy

import (
	"context"
	"fmt"

	"github.com/Miklend/lib-CG/client/node"
	"github.com/Miklend/lib-CG/common"
	"github.com/Miklend/lib-CG/common/logging"
	"github.com/Miklend/lib-CG/models/configs"
	"github.com/Miklend/lib-CG/models/dto"
	"github.com/ethereum/go-ethereum/rpc"
)

type alchemyClient struct {
	networkName string
	apiKey      string
	baseURL     string
	rpcClient   *rpc.Client
	logger      *logging.Logger
}

func NewAlchemyClient(cfg configs.Provider, logger *logging.Logger) (node.Provider, error) {

	fullURL := fmt.Sprintf("%s%s", cfg.BaseURL, cfg.ApiKey)

	rpcClient, err := rpc.Dial(fullURL)
	if err != nil {
		return nil, fmt.Errorf("failed connect to %s: %w", cfg.NetworkName, err)
	}

	return &alchemyClient{
		networkName: cfg.NetworkName,
		apiKey:      cfg.ApiKey,
		baseURL:     cfg.BaseURL,
		rpcClient:   rpcClient,
		logger:      logger,
	}, nil
}
func (a *alchemyClient) Close() error {
	a.rpcClient.Close()
	return nil
}

// Basic
func (a *alchemyClient) BlockByNumber(ctx context.Context, param string) (dto.BlockDTO, error) {
	var result dto.BlockRPCDTO

	err := a.rpcClient.CallContext(ctx, &result, "eth_getBlockByNumber", param, true)
	if err != nil {
		return dto.BlockDTO{}, err
	}

	block := common.ConvertBlockDTO(result)

	return block, nil
}
func (a *alchemyClient) TxByHash(ctx context.Context, param string) (dto.TxDTO, error) {
	var result dto.TxRpcDTO

	err := a.rpcClient.CallContext(ctx, &result, "eth_getTransactionByHash", param)
	if err != nil {
		return dto.TxDTO{}, err
	}

	tx := common.ConvertTxDTO(result)

	return tx, nil
}
func (a *alchemyClient) ReceiptByTxHash(ctx context.Context, param string) (dto.ReceiptDTO, error) {
	var result dto.ReceiptRpcDTO

	err := a.rpcClient.CallContext(ctx, &result, "eth_getTransactionReceipt", param)
	if err != nil {
		return dto.ReceiptDTO{}, err
	}

	receipt := common.ConvertReceiptDTO(result)

	return receipt, nil
}
func (a *alchemyClient) ReceiptByBlockNumber(ctx context.Context, param string) ([]dto.ReceiptDTO, error) {
	var result []dto.ReceiptRpcDTO

	err := a.rpcClient.CallContext(ctx, &result, "eth_getBlockReceipts", param)
	if err != nil {
		return nil, err
	}

	var receipts []dto.ReceiptDTO

	for _, r := range result {
		receipts = append(receipts, common.ConvertReceiptDTO(r))
	}

	return receipts, nil
}

// Subscribe
func (a *alchemyClient) SubscribeBlockWithReceipts(
	ctx context.Context,
	blockCh chan<- *dto.BlockDTO,
) (*rpc.ClientSubscription, error) {

	internalCh := make(chan map[string]interface{})

	sub, err := a.rpcClient.Subscribe(ctx, "eth", internalCh, "newHeads", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to newHeads: %w", err)
	}

	go func() {
		defer close(blockCh)

		for msg := range internalCh {
			select {
			case <-ctx.Done():
				return
			default:
			}

			numberHex, ok := msg["number"].(string)
			if !ok {
				a.logger.Warn("no number field in newHeads message")
				continue
			}

			block, err := a.BlockByNumber(ctx, numberHex)
			if err != nil {
				a.logger.Errorf("fetching block %s: %v", numberHex, err)
				continue
			}

			receipts, err := a.ReceiptByBlockNumber(ctx, numberHex)
			if err != nil {
				a.logger.Errorf("fetching receipts for block %s: %v", numberHex, err)
				continue
			}

			for i := range block.Transactions {
				for _, r := range receipts {
					if block.Transactions[i].Hash == r.TransactionHash {
						block.Transactions[i].Receipt = r
						break
					}
				}
			}

			select {
			case blockCh <- &block:
			case <-ctx.Done():
				return
			}
		}
	}()

	return sub, nil
}
func (a *alchemyClient) SubscribePendingTransactions(
	ctx context.Context,
	txsCh chan<- *dto.TxDTO,
) (*rpc.ClientSubscription, error) {

	internalCh := make(chan dto.TxDTO)

	params := map[string]interface{}{
		"hashesOnly": false,
	}

	sub, err := a.rpcClient.Subscribe(ctx, "eth", internalCh, "alchemy_pendingTransactions", params)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to alchemy_pendingTransactions: %w", err)
	}

	go func() {
		defer close(txsCh)
		for tx := range internalCh {
			select {
			case <-ctx.Done():
				return
			case txsCh <- &tx:
			}
		}
	}()

	return sub, nil
}

// Batch
func (a *alchemyClient) BatchRequest(ctx context.Context, requests []rpc.BatchElem) ([]rpc.BatchElem, error) {
	if err := a.rpcClient.BatchCallContext(ctx, requests); err != nil {
		return nil, err
	}
	return requests, nil
}
func (a *alchemyClient) BatchBlockByNumber(ctx context.Context, numbers []string) (map[dto.Hash]dto.BlockDTO, error) {
	blocks := make(map[dto.Hash]dto.BlockDTO)
	results := make([]dto.BlockRPCDTO, len(numbers))

	batch := make([]rpc.BatchElem, len(numbers))
	for i, n := range numbers {
		batch[i] = rpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []interface{}{n, true},
			Result: &results[i],
		}
	}

	res, err := a.BatchRequest(ctx, batch)
	if err != nil {
		return nil, fmt.Errorf("batch block request failed: %w", err)
	}

	for _, r := range res {
		if r.Error != nil || r.Result == nil {
			a.logger.Warnf("not found block, err: %v", r.Error)
			continue
		}
		raw := *(r.Result.(*dto.BlockRPCDTO))
		block := common.ConvertBlockDTO(raw)
		blocks[block.Hash] = block
	}

	return blocks, nil
}
func (a *alchemyClient) BatchReceiptByBlockNumber(ctx context.Context, numbers []string) (map[dto.Hash][]dto.ReceiptDTO, error) {
	receiptsMap := make(map[dto.Hash][]dto.ReceiptDTO)
	results := make([][]dto.ReceiptRpcDTO, len(numbers))

	batch := make([]rpc.BatchElem, len(numbers))
	for i, n := range numbers {
		batch[i] = rpc.BatchElem{
			Method: "eth_getBlockReceipts",
			Args:   []interface{}{n},
			Result: &results[i],
		}
	}

	res, err := a.BatchRequest(ctx, batch)
	if err != nil {
		return nil, fmt.Errorf("batch block receipt request failed: %w", err)
	}

	for i, r := range res {
		if r.Error != nil || r.Result == nil {
			a.logger.Warnf("not found block receipt, err: %v", r.Error)
			continue
		}

		blockReceipts := results[i]
		if len(blockReceipts) == 0 {
			continue
		}

		var converted []dto.ReceiptDTO
		for _, rec := range blockReceipts {
			converted = append(converted, common.ConvertReceiptDTO(rec))
		}

		receiptsMap[converted[0].BlockHash] = converted
	}

	return receiptsMap, nil
}
func (a *alchemyClient) BatchReceiptByTxHash(ctx context.Context, hashes []string) (map[dto.Hash]dto.ReceiptDTO, error) {
	receiptsMap := make(map[dto.Hash]dto.ReceiptDTO)
	results := make([]dto.ReceiptRpcDTO, len(hashes))

	batch := make([]rpc.BatchElem, len(hashes))
	for i, h := range hashes {
		batch[i] = rpc.BatchElem{
			Method: "eth_getTransactionReceipt",
			Args:   []interface{}{h},
			Result: &results[i],
		}
	}

	res, err := a.BatchRequest(ctx, batch)
	if err != nil {
		return nil, fmt.Errorf("batch transaction receipt request failed: %w", err)
	}

	for i, r := range res {
		if r.Error != nil || r.Result == nil {
			a.logger.Warnf("receipt not found for tx %s", hashes[i])
			continue
		}

		raw := results[i]
		converted := common.ConvertReceiptDTO(raw)
		receiptsMap[converted.TransactionHash] = converted
	}

	return receiptsMap, nil
}
func (a *alchemyClient) BatchBlockWithReceiptByNumber(ctx context.Context, numbers []string) (map[dto.Hash]dto.BlockDTO, error) {
	blocksMap := make(map[dto.Hash]dto.BlockDTO)
	blockResults := make([]dto.BlockRPCDTO, len(numbers))
	receiptResults := make([][]dto.ReceiptRpcDTO, len(numbers))

	blockBatch := make([]rpc.BatchElem, len(numbers))
	for i, n := range numbers {
		blockBatch[i] = rpc.BatchElem{
			Method: "eth_getBlockByNumber",
			Args:   []interface{}{n, true},
			Result: &blockResults[i],
		}
	}

	receiptBatch := make([]rpc.BatchElem, len(numbers))
	for i, n := range numbers {
		receiptBatch[i] = rpc.BatchElem{
			Method: "eth_getBlockReceipts",
			Args:   []interface{}{n},
			Result: &receiptResults[i],
		}
	}

	if _, err := a.BatchRequest(ctx, blockBatch); err != nil {
		return nil, fmt.Errorf("batch block request failed: %w", err)
	}
	if _, err := a.BatchRequest(ctx, receiptBatch); err != nil {
		return nil, fmt.Errorf("batch receipt request failed: %w", err)
	}

	for i := range numbers {
		block := common.ConvertBlockDTO(blockResults[i])

		receipts := make([]dto.ReceiptDTO, 0, len(receiptResults[i]))
		for _, r := range receiptResults[i] {
			receipts = append(receipts, common.ConvertReceiptDTO(r))
		}

		for j := range block.Transactions {
			for _, r := range receipts {
				if block.Transactions[j].Hash == r.TransactionHash {
					block.Transactions[j].Receipt = r
					break
				}
			}
		}

		blocksMap[block.Hash] = block
	}

	return blocksMap, nil
}
