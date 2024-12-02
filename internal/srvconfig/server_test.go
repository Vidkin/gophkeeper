package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServerConfig(t *testing.T) {
	os.Args = []string{
		"cmd",
		"-a", "localhost:8080",
		"-l", "debug",
		"-crypto-key-private", "path",
		"-crypto-key-public", "path",
		"-db-key", "strongDBKey2Ks5nM2J5JaI59PPEhL1x",
		"-minio-endpoint", "test",
		"-minio-secret", "test",
		"-minio-id", "test",
		"-j", "test"}
	config, err := NewServerConfig()
	require.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "localhost:8080", config.ServerAddress.Address)
	assert.Equal(t, "debug", config.LogLevel)
	assert.Equal(t, 3, config.RetryCount)
}

func TestNewServerConfig_MissingRequiredFields(t *testing.T) {
	os.Args = []string{"cmd"}
	config, err := NewServerConfig()
	require.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "you must pass the path to public and private keys pem files")
}

func TestLoadJSONConfig(t *testing.T) {
	jsonConfig := `{
		"address": {"host": "127.0.0.1", "port": "8080"},
		"database_dsn": "user:password@tcp(localhost:3306)/dbname",
		"hash_key": "defaultHashKey",
		"jwt_key": "defaultHashKey",
		"minio_endpoint": "localhost:9000",
		"minio_access_key_id": "minio_access_key",
		"minio_secret_access_key": "minio_secret_key",
		"database_key": "strongDBKey2Ks5nM2J5JaI59PPEhL1x"
	}`
	tmpFile, err := os.CreateTemp("", "config.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(jsonConfig)
	require.NoError(t, err)
	tmpFile.Close()

	os.Args = []string{
		"cmd",
		"-c", tmpFile.Name(),
		"-crypto-key-private", "path",
		"-crypto-key-public", "path",
		"-minio-endpoint", "test",
		"-minio-secret", "test",
		"-minio-id", "test",
		"-db-key", "strongDBKey2Ks5nM2J5JaI59PPEhL1x",
		"-j", "defaultHashKey",
		"-d", "user:password@tcp(localhost:3306)/dbname",
		"-k", "defaultHashKey",
	}
	config, err := NewServerConfig()
	require.NoError(t, err)
	assert.NotNil(t, config)
	assert.Equal(t, "127.0.0.1:8080", config.ServerAddress.Address)
	assert.Equal(t, "user:password@tcp(localhost:3306)/dbname", config.DatabaseDSN)
	assert.Equal(t, "defaultHashKey", config.Key)
}

func TestLoadJSONConfig_MissingFields(t *testing.T) {
	jsonConfig := `{
		"address": {"host": "localhost", "port": "8080"}
	}`
	tmpFile, err := os.CreateTemp("", "config.json")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(jsonConfig)
	require.NoError(t, err)
	tmpFile.Close()

	os.Args = []string{"cmd", "-c", tmpFile.Name()}
	_, err = NewServerConfig()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "you must pass the path to public and private keys pem files")
}

func TestServerConfig_MissingCryptoKeys(t *testing.T) {
	os.Args = []string{"cmd", "-a", "localhost:8080", "-k", "examplehashkey1234567890123456", "-j", "examplejwtkey", "-minio-endpoint", "localhost:9000", "-minio-id", "minio_access_key", "-minio-secret", "minio_secret_key", "-db-key", "exampledatabasekey1234567890123456"}
	config, err := NewServerConfig()
	require.Error(t, err)
	assert.Nil(t, config)
	assert.Contains(t, err.Error(), "you must pass the path to public and private keys pem files, see --help")
}

func TestServerConfig_Cleanup(t *testing.T) {
	os.Unsetenv("MINIO_ENDPOINT")
	os.Unsetenv("MINIO_ACCESS_KEY_ID")
	os.Unsetenv("MINIO_SECRET_ACCESS_KEY")
	os.Unsetenv("DATABASE_DSN")
	os.Unsetenv("JWT_KEY")
	os.Unsetenv("KEY")
	os.Unsetenv("DATABASE_KEY")
	os.Unsetenv("CRYPTO_KEY_PUBLIC")
	os.Unsetenv("CRYPTO_KEY_PRIVATE")
}
