package main

import (
	"context"
	"encoding/json"
	"fmt"

	"os"
	"os/signal"
	"syscall"

	fabricClient "github.com/Miklend/lib-CG/client/fabric"
	"github.com/Miklend/lib-CG/common/logging"
	"github.com/Miklend/lib-CG/models/configs"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	logger := logging.GetLogger()

	cfg := configs.Provider{
		ProviderType: "alchemy",
		NetworkName:  "eth",
		BaseURL:      "https://eth-mainnet.g.alchemy.com/v2/",
		ApiKey:       "Uzb0sDJmjCuvs21_OyLbn",
	}

	provaider, err := fabricClient.NewProvider(cfg, logger)
	if err != nil {
		logger.Fatal(err)
	}

	blksN := []string{"0x169F2E5", "0x169F2E6"}
	blocks, err := provaider.BatchBlockByNumber(ctx, blksN)
	if err != nil {
		logger.Fatal(err)
	}

	if err := saveToJSONFile("blocks.json", blocks); err != nil {
		logger.Fatal(err)
	}
}

func saveToJSONFile(filename string, data any) error {
	// сериализуем с отступами для читаемости
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data to JSON: %w", err)
	}

	// создаём файл (или перезаписываем)
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filename, err)
	}
	defer file.Close()

	_, err = file.Write(bytes)
	if err != nil {
		return fmt.Errorf("failed to write JSON to file: %w", err)
	}

	return nil
}
