package handlers

import (
	"crypto/tls"
	"database/sql"
	"math/rand"
	"net"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/Vidkin/gophkeeper/internal/storage"
)

const (
	TokenFileName = "gophkeeperJWT.tmp"
	expiredToken  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MzIzODczMzgsIlVzZXJJRCI6MX0.B6kBiV1YOiDZd1oxp4weHgkFtJcN5VebwWpRD70uQDw"
)

func GetTLSListener(addr, certFile, keyFile string) (net.Listener, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	cfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"h2"}}
	return tls.Listen("tcp", addr, cfg)
}

func setExpiredToken(t *testing.T) {
	err := os.Remove(path.Join(os.TempDir(), TokenFileName))
	if !os.IsNotExist(err) {
		require.NoError(t, err)
	}
	f, err := os.Create(path.Join(os.TempDir(), TokenFileName))
	require.NoError(t, err)
	defer f.Close()
	_, err = f.WriteString(expiredToken)
	require.NoError(t, err)
}

const letters = "abcdefghijklmnopqrstuvwxyz"

func randomDBName(length int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func setupTestDB(t *testing.T) (*storage.PostgresStorage, string) {
	connStr := "user=postgres password=postgres dbname=postgres host=127.0.0.1 port=5432 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)

	tempDBName := randomDBName(10)
	_, err = db.Exec("SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = $1", tempDBName)
	require.NoError(t, err)
	db.Exec("DROP DATABASE " + tempDBName)
	_, err = db.Exec("CREATE DATABASE " + tempDBName)
	require.NoError(t, err)

	tempConnStr := "user=postgres password=postgres dbname=" + tempDBName + " host=127.0.0.1 port=5432 sslmode=disable"
	st, err := storage.NewPostgresStorage(tempConnStr)
	require.NoError(t, err)

	return st, tempDBName
}

func teardownTestDB(t *testing.T, db *sql.DB, dbName string) {
	db.Close()

	connStr := "user=postgres password=postgres dbname=postgres host=127.0.0.1 port=5432 sslmode=disable"
	mainDB, err := sql.Open("postgres", connStr)
	require.NoError(t, err)

	_, err = mainDB.Exec("SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = $1", dbName)
	require.NoError(t, err)

	_, err = mainDB.Exec("DROP DATABASE " + dbName)
	require.NoError(t, err)

	mainDB.Close()
}
