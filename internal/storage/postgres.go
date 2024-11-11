package storage

import (
	"context"
	"database/sql"
	"embed"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	"github.com/Vidkin/gophkeeper/internal/logger"
)

//go:embed migrations/*.sql
var Migrations embed.FS

type PostgresStorage struct {
	Conn *sql.DB
}

func NewPostgresStorage(dbDSN string) (*PostgresStorage, error) {
	var p PostgresStorage
	db, err := sql.Open("pgx", dbDSN)
	if err != nil {
		logger.Log.Fatal("error open sql connection", zap.Error(err))
		return nil, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Log.Fatal("can't create postgres driver for migrations", zap.Error(err))
		return nil, err
	}

	d, err := iofs.New(Migrations, "migrations")
	if err != nil {
		logger.Log.Fatal("can't get migrations from FS", zap.Error(err))
		return nil, err
	}

	m, err := migrate.NewWithInstance("iofs", d, "postgres", driver)
	if err != nil {
		logger.Log.Fatal("can't create new migrate instance", zap.Error(err))
		return nil, err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Log.Fatal("can't exec migrations", zap.Error(err))
		return nil, err
	}
	p.Conn = db
	return &p, nil
}

func (p *PostgresStorage) Ping(ctx context.Context) error {
	return p.Conn.PingContext(ctx)
}

func (p *PostgresStorage) Close() error {
	return p.Conn.Close()
}
