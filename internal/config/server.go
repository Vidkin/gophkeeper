package config

import (
	"encoding/json"
	"errors"
	"flag"
	"os"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"

	"github.com/Vidkin/gophkeeper/internal/logger"
)

// ServerConfig holds the configuration settings for the server.
//
// This struct contains various fields that define how the server operates,
// including its address, storage settings, logging preferences, and more.
// The fields can be populated from environment variables, allowing for
// flexible configuration without hardcoding values.
type ServerConfig struct {
	ServerAddress    *ServerAddress `json:"address"`
	LogLevel         string
	ConfigPath       string `env:"CONFIG"`
	DatabaseDSN      string `env:"DATABASE_DSN" json:"database_dsn"`
	Key              string `env:"KEY" json:"hash_key"`
	CryptoKeyPublic  string `env:"CRYPTO_KEY_PUBLIC"`
	CryptoKeyPrivate string `env:"CRYPTO_KEY_PRIVATE"`
	RetryCount       int
}

// NewServerConfig initializes a new ServerConfig instance with default values
// and parses command-line flags and environment variables to populate its fields
//
// Returns:
// - A pointer to the newly created and initialized ServerConfig instance.
// - An error if the configuration parsing fails; otherwise, nil.
func NewServerConfig() (*ServerConfig, error) {
	var config ServerConfig
	config.ServerAddress = NewServerAddress()
	config.RetryCount = 3
	err := config.parseFlags()
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (config *ServerConfig) parseFlags() error {
	fs := flag.NewFlagSet("serverFlagSet", flag.ContinueOnError)

	fs.Var(config.ServerAddress, "a", "Net address host:port")
	fs.StringVar(&config.ConfigPath, "c", "", "Path to json config file")
	fs.StringVar(&config.ConfigPath, "config", "", "Path to json config file")
	fs.StringVar(&config.LogLevel, "l", "info", "Log level")
	fs.StringVar(&config.DatabaseDSN, "d", "", "Database DSN")
	fs.StringVar(&config.Key, "k", "", "Hash key")
	fs.StringVar(&config.CryptoKeyPublic, "crypto-key-public", "", "Path to public key pem file")
	fs.StringVar(&config.CryptoKeyPrivate, "crypto-key-private", "", "Path to private key pem file")

	if err := fs.Parse(os.Args[1:]); err != nil {
		logger.Log.Error("error parse server flags", zap.Error(err))
		return err
	}
	if config.ConfigPath != "" {
		if err := config.loadJSONConfig(config.ConfigPath); err != nil {
			logger.Log.Error("error parse json config file", zap.Error(err))
		}
	}

	err := env.Parse(config)
	if err != nil {
		return err
	}

	if config.ServerAddress.Address == "" {
		config.ServerAddress.Address = config.ServerAddress.String()
	}

	if config.CryptoKeyPublic == "" || config.CryptoKeyPrivate == "" {
		return errors.New("you should pass the path to public and private keys pem files, see --help")
	}

	return nil
}

func (config *ServerConfig) loadJSONConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var jsonServerConfig ServerConfig
	if err = json.Unmarshal(data, &jsonServerConfig); err != nil {
		return err
	}

	if config.ServerAddress.Address == "" {
		config.ServerAddress = jsonServerConfig.ServerAddress
	}

	dbDSNPassed := false
	hashKeyPassed := false

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--d", "-d":
			dbDSNPassed = true
		case "--k", "-k":
			hashKeyPassed = true
		}
	}

	if !dbDSNPassed {
		config.DatabaseDSN = jsonServerConfig.DatabaseDSN
	}

	if !hashKeyPassed {
		config.Key = jsonServerConfig.Key
	}

	return nil
}
