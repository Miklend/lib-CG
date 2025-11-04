package common

import (
	"strconv"
	"strings"

	"github.com/Miklend/lib-CG/models/dto"
)

func trim0x(s string) string {
	return strings.TrimPrefix(s, "0x")
}
func hexToUint64(hexStr string) uint64 {
	if hexStr == "" {
		return 0
	}
	val, err := strconv.ParseUint(trim0x(hexStr), 16, 64)
	if err != nil {
		return 0
	}
	return val
}
func toHashArray(list []string) []dto.Hash {
	out := make([]dto.Hash, len(list))
	for i, v := range list {
		out[i] = dto.Hash(v)
	}
	return out
}

func ConvertBlockDTO(b dto.BlockRPCDTO) dto.BlockDTO {
	var txs []dto.TxDTO
	for _, l := range b.Transactions {
		txs = append(txs, ConvertTxDTO(l))
	}

	return dto.BlockDTO{
		Hash:             dto.Hash(b.Hash),
		Number:           hexToUint64(b.Number),
		ParentHash:       dto.Hash(b.ParentHash),
		Nonce:            b.Nonce,
		Sha3Uncles:       dto.Hash(b.Sha3Uncles),
		LogsBloom:        b.LogsBloom,
		TransactionsRoot: dto.Hash(b.TransactionsRoot),
		StateRoot:        dto.Hash(b.StateRoot),
		ReceiptsRoot:     dto.Hash(b.ReceiptsRoot),
		Miner:            dto.Address(b.Miner),
		Difficulty:       b.Difficulty,
		TotalDifficulty:  b.TotalDifficulty,
		Size:             hexToUint64(b.Size),
		ExtraData:        b.ExtraData,
		GasLimit:         hexToUint64(b.GasLimit),
		GasUsed:          hexToUint64(b.GasUsed),
		BaseFeePerGas:    b.BaseFeePerGas,
		Timestamp:        dto.Time(hexToUint64(b.Timestamp)),
		MixHash:          dto.Hash(b.MixHash),
		Transactions:     txs,
		Uncles:           toHashArray(b.Uncles),
	}
}
func ConvertTxDTO(t dto.TxRpcDTO) dto.TxDTO {
	var to *dto.Address
	if t.To != nil {
		addr := dto.Address(*t.To)
		to = &addr
	}

	var accessList []dto.AccessListEntryDTO
	for _, a := range t.AccessList {
		accessList = append(accessList, dto.AccessListEntryDTO{
			Address:     dto.Address(a.Address),
			StorageKeys: a.StorageKeys,
		})
	}

	return dto.TxDTO{
		Hash:                 dto.Hash(t.Hash),
		BlockHash:            dto.Hash(t.BlockHash),
		BlockNumber:          hexToUint64(t.BlockNumber),
		TransactionIndex:     hexToUint64(t.TransactionIndex),
		From:                 dto.Address(t.From),
		To:                   to,
		Value:                t.Value,
		Gas:                  hexToUint64(t.Gas),
		GasPrice:             t.GasPrice,
		Input:                t.Input,
		Nonce:                hexToUint64(t.Nonce),
		Type:                 t.Type,
		MaxFeePerGas:         t.MaxFeePerGas,
		MaxPriorityFeePerGas: t.MaxPriorityFeePerGas,
		ChainID:              t.ChainID,
		V:                    t.V,
		R:                    t.R,
		S:                    t.S,
		AccessList:           accessList,
	}
}
func ConvertReceiptDTO(r dto.ReceiptRpcDTO) dto.ReceiptDTO {
	var to *dto.Address
	if r.To != nil {
		addr := dto.Address(*r.To)
		to = &addr
	}
	var contract *dto.Address
	if r.ContractAddress != nil {
		addr := dto.Address(*r.ContractAddress)
		contract = &addr
	}

	logs := make([]dto.LogDTO, len(r.Logs))
	for i, l := range r.Logs {
		logs[i] = ConvertLogDTO(l)
	}

	return dto.ReceiptDTO{
		TransactionHash:   dto.Hash(r.TransactionHash),
		TransactionIndex:  hexToUint64(r.TransactionIndex),
		BlockHash:         dto.Hash(r.BlockHash),
		BlockNumber:       hexToUint64(r.BlockNumber),
		From:              dto.Address(r.From),
		To:                to,
		ContractAddress:   contract,
		CumulativeGasUsed: hexToUint64(r.CumulativeGasUsed),
		GasUsed:           hexToUint64(r.GasUsed),
		EffectiveGasPrice: r.EffectiveGasPrice,
		Status:            r.Status,
		LogsBloom:         r.LogsBloom,
		Logs:              logs,
	}
}
func ConvertLogDTO(l dto.LogRpcDTO) dto.LogDTO {
	return dto.LogDTO{
		BlockNumber:      hexToUint64(l.BlockNumber),
		BlockHash:        dto.Hash(l.BlockHash),
		TransactionHash:  dto.Hash(l.TransactionHash),
		TransactionIndex: hexToUint64(l.TransactionIndex),
		LogIndex:         hexToUint64(l.LogIndex),
		Address:          dto.Address(l.Address),
		Data:             l.Data,
		Topics:           l.Topics,
		Removed:          l.Removed,
	}
}
