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
	"github.com/Vidkin/gophkeeper/internal/model"
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

func (p *PostgresStorage) AddUser(ctx context.Context, login, password string) error {
	_, err := p.Conn.ExecContext(ctx, "INSERT INTO users (login, password) VALUES ($1, $2)", login, password)
	return err
}

func (p *PostgresStorage) AddFile(ctx context.Context, bucketName, fileName, contentType, description string, userID int64, fileSize int64) error {
	var count int
	row := p.Conn.QueryRowContext(
		ctx,
		"SELECT count(*) FROM files WHERE file_name=$1 and user_id=$2",
		fileName, userID)
	err := row.Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		_, err = p.Conn.ExecContext(
			ctx,
			"UPDATE files SET content_type=$1, file_size=$2, description=$3 WHERE file_name=$4",
			contentType, fileSize, description, fileName)
		return err
	}

	_, err = p.Conn.ExecContext(
		ctx,
		"INSERT INTO files (user_id, bucket_name, file_name, content_type, file_size, description) VALUES ($1, $2, $3, $4, $5, $6)",
		userID, bucketName, fileName, contentType, fileSize, description)
	return err
}

func (p *PostgresStorage) GetFile(ctx context.Context, fileID int64) (*model.File, error) {
	row := p.Conn.QueryRowContext(
		ctx,
		"SELECT user_id, id, file_name, bucket_name, content_type, description, file_size, created_at id FROM files WHERE id = $1",
		fileID)

	var f model.File
	if err := row.Scan(&f.UserID, &f.ID, &f.FileName, &f.BucketName, &f.ContentType, &f.Description, &f.FileSize, &f.CreatedAt); err != nil {
		return nil, err
	}
	return &f, nil
}

func (p *PostgresStorage) GetFiles(ctx context.Context, userID int64) ([]*model.File, error) {
	rows, err := p.Conn.QueryContext(
		ctx,
		"SELECT user_id, id, file_name, bucket_name, content_type, description, file_size, created_at id FROM files WHERE user_id = $1",
		userID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			logger.Log.Error("error close rows", zap.Error(err))
		}
	}(rows)

	var files []*model.File
	for rows.Next() {
		var f model.File
		if err = rows.Scan(&f.UserID, &f.ID, &f.FileName, &f.BucketName, &f.ContentType, &f.Description, &f.FileSize, &f.CreatedAt); err != nil {
			return nil, err
		}
		files = append(files, &f)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}

func (p *PostgresStorage) GetUser(ctx context.Context, login string) (*model.User, error) {
	row := p.Conn.QueryRowContext(ctx, "SELECT login, password, id FROM users WHERE login = $1", login)

	var u model.User
	if err := row.Scan(&u.Login, &u.Password, &u.ID); err != nil {
		return nil, err
	}
	return &u, nil
}

func (p *PostgresStorage) AddUserCredentials(ctx context.Context, cred *model.Credentials) error {
	_, err := p.Conn.ExecContext(
		ctx,
		"INSERT INTO user_credentials (login, password, description, user_id) VALUES ($1, $2, $3, $4)",
		cred.Login, cred.Password, cred.Description, cred.UserID)
	return err
}

func (p *PostgresStorage) GetUserCredentials(ctx context.Context, userID int64) ([]*model.Credentials, error) {
	rows, err := p.Conn.QueryContext(ctx, "SELECT id, user_id, login, password, description FROM user_credentials WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			logger.Log.Error("error close rows", zap.Error(err))
		}
	}(rows)

	var creds []*model.Credentials
	for rows.Next() {
		var c model.Credentials
		if err = rows.Scan(&c.ID, &c.UserID, &c.Login, &c.Password, &c.Description); err != nil {
			return nil, err
		}
		creds = append(creds, &c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return creds, nil
}

func (p *PostgresStorage) GetUserCredential(ctx context.Context, credID int64) (*model.Credentials, error) {
	row := p.Conn.QueryRowContext(ctx, "SELECT id, user_id, login, password, description FROM user_credentials WHERE id = $1", credID)

	var cred model.Credentials
	if err := row.Scan(&cred.ID, &cred.UserID, &cred.Login, &cred.Password, &cred.Description); err != nil {
		return nil, err
	}

	return &cred, nil
}

func (p *PostgresStorage) AddCard(ctx context.Context, card *model.BankCard) error {
	_, err := p.Conn.ExecContext(
		ctx,
		"INSERT INTO bank_cards (user_id, card_number, expiration_date, cvv, owner, description) "+
			"VALUES ($1, $2, $3, $4, $5, $6)", card.UserID, card.Number, card.ExpireDate, card.CVV, card.Owner, card.Description)
	return err
}

func (p *PostgresStorage) GetBankCards(ctx context.Context, userID int64) ([]*model.BankCard, error) {
	rows, err := p.Conn.QueryContext(ctx, "SELECT id, user_id, owner, card_number, expiration_date, cvv, description FROM bank_cards WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			logger.Log.Error("error close rows", zap.Error(err))
		}
	}(rows)

	var cards []*model.BankCard
	for rows.Next() {
		var b model.BankCard
		if err = rows.Scan(&b.ID, &b.UserID, &b.Owner, &b.Number, &b.ExpireDate, &b.CVV, &b.Description); err != nil {
			return nil, err
		}
		cards = append(cards, &b)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cards, nil
}

func (p *PostgresStorage) GetBankCard(ctx context.Context, cardID int64) (*model.BankCard, error) {
	row := p.Conn.QueryRowContext(ctx, "SELECT id, user_id, owner, card_number, expiration_date, cvv, description FROM bank_cards WHERE id = $1", cardID)

	var card model.BankCard
	if err := row.Scan(&card.ID, &card.UserID, &card.Owner, &card.Number, &card.ExpireDate, &card.CVV, &card.Description); err != nil {
		return nil, err
	}

	return &card, nil
}
