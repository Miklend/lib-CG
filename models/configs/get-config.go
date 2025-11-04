package configs

import (
	"flag"

	"log/slog"
	"os"
	"sync"

	"github.com/Miklend/lib-CG/common/logging"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Provider Provider `json:"provaider"`
	Broker   Broker   `json:"broker"`
}

const (
	flagConfigPathName = "configs"
	envConfigPathName  = "CONFIG_PATH"
	dotEnvFileName     = ".env"
)

var (
	instance *Config
	once     sync.Once
)

func GetConfig(logger *logging.Logger) *Config {
	logger.Debug("start get config")

	once.Do(func() {
		_ = godotenv.Load(dotEnvFileName)

		var configPath string
		flag.StringVar(&configPath, flagConfigPathName, "", "path to config file (e.g., ./configs/config.yaml)")
		flag.Parse()

		if path, ok := os.LookupEnv(envConfigPathName); ok && path != "" {
			configPath = path
		}

		if configPath == "" {

			possiblePaths := []string{
				"./configs/configs.yaml",
				"./configs/config.yaml",
				"../configs/configs.yaml",
				"../../configs/configs.yaml",
			}

			for _, path := range possiblePaths {
				if _, err := os.Stat(path); err == nil {
					configPath = path
					logger.Debugf("Found config at: %s", path)
					break
				}
			}

			if configPath == "" {
				configPath = "./configs/configs.yaml"
			}
		}

		instance = &Config{}

		if readErr := cleanenv.ReadConfig(configPath, instance); readErr != nil {
			description, descrErr := cleanenv.GetDescription(instance, nil)
			if descrErr != nil {
				panic(descrErr)
			}

			slog.Error(
				"Failed to read config. Ensure 'config.yaml' exists or all required env variables are set.",
				slog.String("error", readErr.Error()),
				slog.String("config_path", configPath),
				slog.String("description", description),
			)
			os.Exit(1)
		}

		logger.Debug("config loaded successfully", slog.Any("config", instance))
	})

	return instance
}
