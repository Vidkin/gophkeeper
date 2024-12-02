package commands

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitConfig_ValidConfig(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "*config.json")
	tmpFile.Write([]byte(`{
	  "address": "127.0.0.1:8080",
	  "crypto_key_public_path": "/test/public.crt",
	  "hash_key": "test_hash_key",
	  "secret_key": "test_secret_key"}`))
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	cfgFilePath = tmpFile.Name()
	hashKey = "test_hash_key"
	secretKey = "test_secret_key"

	initConfig()

	assert.Equal(t, "test_hash_key", viper.GetString("hash_key"))
	assert.Equal(t, "test_secret_key", viper.GetString("secret_key"))
}

func TestInitConfig_MissingConfigFile(t *testing.T) {
	cfgFilePath = "non_existing_file.json"
	hashKey = "test_hash_key"
	secretKey = "test_secret_key"

	assert.Panics(t, func() { initConfig() }, "Expected panic when config file is missing")
}

func TestInitConfig_ConfigNameIsEmpty(t *testing.T) {
	cfgFilePath = ""
	hashKey = "test_hash_key"
	secretKey = "test_secret_key"

	assert.Panics(t, func() { initConfig() }, "Expected panic when config file name is empty")
}

func TestInitConfig_MissingHashKey(t *testing.T) {
	viper.Reset()
	hashKey = ""
	secretKey = "test_secret_key"
	tmpFile, err := os.CreateTemp("", "*config.json")
	tmpFile.Write([]byte(`{
	  "address": "127.0.0.1:8080",
	  "crypto_key_public_path": "/test/public.crt",
	  "secret_key": "test_secret_key"}`))
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	cfgFilePath = tmpFile.Name()

	assert.Panics(t, func() {
		initConfig()
		fmt.Println(viper.AllSettings())
	}, "Expected panic when hash_key is missing")
}

func TestInitConfig_MissingSecretKey(t *testing.T) {
	viper.Reset()
	secretKey = ""
	hashKey = "test_hash_key"
	tmpFile, err := os.CreateTemp("", "*config.json")
	tmpFile.Write([]byte(`{
	  "address": "127.0.0.1:8080",
	  "crypto_key_public_path": "/test/public.crt",
	  "hash_key": "test_hash_key"}`))
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	cfgFilePath = tmpFile.Name()

	assert.Panics(t, func() { initConfig() }, "Expected panic when secret_key is missing")
}
