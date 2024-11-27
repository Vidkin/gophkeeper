// Package storage provides functionality for interacting with a PostgreSQL database.
//
// This package includes the PostgresStorage struct, which implements methods for managing users, files,
// notes, bank cards, and user credentials in a PostgreSQL database.
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

// PostgresStorage represents a storage backend using PostgreSQL.
type PostgresStorage struct {
	Conn *sql.DB
}

// NewPostgresStorage initializes a new PostgresStorage instance and applies database migrations.
//
// Parameters:
//   - dbDSN: A string representing the Data Source Name (DSN) for connecting to the PostgreSQL database.
//
// Returns:
//   - A pointer to a PostgresStorage instance.
//   - An error if the connection to the database could not be established or if migrations fail.
//
// The function opens a connection to the PostgreSQL database, creates a migration instance, and applies
// any pending migrations. If successful, it returns a PostgresStorage instance with an active database connection.
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

// Ping checks the connection to the PostgreSQL database.
//
// Parameters:
//   - ctx: The context for the operation.
//
// Returns:
//   - An error if the connection is not alive or if there is an issue with the database.
func (p *PostgresStorage) Ping(ctx context.Context) error {
	return p.Conn.PingContext(ctx)
}

// Close closes the database connection.
//
// Returns:
//   - An error if there is an issue closing the connection.
func (p *PostgresStorage) Close() error {
	return p.Conn.Close()
}

// AddUser adds a new user to the database.
//
// Parameters:
//   - ctx: The context for the operation.
//   - login: A string representing the user's login name.
//   - password: A string representing the user's password.
//
// Returns:
//   - An error if the operation fails.
func (p *PostgresStorage) AddUser(ctx context.Context, login, password string) error {
	_, err := p.Conn.ExecContext(ctx, "INSERT INTO users (login, password) VALUES ($1, $2)", login, password)
	return err
}

// AddFile adds a new file or updates an existing file for a user.
//
// Parameters:
//   - ctx: The context for the operation.
//   - bucketName: A string representing the name of the bucket where the file is stored.
//   - fileName: A string representing the name of the file.
//   - description: A string providing additional information about the file.
//   - userID: An int64 representing the unique identifier of the user.
//   - fileSize: An int64 representing the size of the file in bytes.
//
// Returns:
//   - An error if the operation fails.
func (p *PostgresStorage) AddFile(ctx context.Context, bucketName, fileName, description string, userID int64, fileSize int64) error {
	var count int
	row := p.Conn.QueryRowContext(
		ctx,
		"SELECT count(*) FROM files WHERE file_name=$1 and user_id=$2",
		fileName, userID)
	if err := row.Scan(&count); err != nil {
		return err
	}

	if count > 0 {
		_, err := p.Conn.ExecContext(
			ctx,
			"UPDATE files SET file_size=$1, description=$2 WHERE file_name=$3", fileSize, description, fileName)
		return err
	}

	_, err := p.Conn.ExecContext(
		ctx,
		"INSERT INTO files (user_id, bucket_name, file_name, file_size, description) VALUES ($1, $2, $3, $4, $5)",
		userID, bucketName, fileName, fileSize, description)
	return err
}

// GetFile retrieves a file by its name from the database.
//
// Parameters:
//   - ctx: The context for the operation.
//   - fileName: A string representing the name of the file to retrieve.
//
// Returns:
//   - A pointer to a model.File instance containing the file information.
//   - An error if the operation fails or if the file is not found.
func (p *PostgresStorage) GetFile(ctx context.Context, fileName string) (*model.File, error) {
	row := p.Conn.QueryRowContext(
		ctx,
		"SELECT user_id, id, file_name, bucket_name, description, file_size, created_at FROM files WHERE file_name = $1",
		fileName)

	var f model.File
	if err := row.Scan(&f.UserID, &f.ID, &f.FileName, &f.BucketName, &f.Description, &f.FileSize, &f.CreatedAt); err != nil {
		return nil, err
	}
	return &f, nil
}

// GetFiles retrieves all files associated with a user.
//
// Parameters:
//   - ctx: The context for the operation.
//   - userID: An int64 representing the unique identifier of the user.
//
// Returns:
//   - A slice of pointers to model.File instances containing the user's files.
//   - An error if the operation fails.
func (p *PostgresStorage) GetFiles(ctx context.Context, userID int64) ([]*model.File, error) {
	rows, err := p.Conn.QueryContext(
		ctx,
		"SELECT user_id, id, file_name, bucket_name, description, file_size, created_at id FROM files WHERE user_id = $1",
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
		if err = rows.Scan(&f.UserID, &f.ID, &f.FileName, &f.BucketName, &f.Description, &f.FileSize, &f.CreatedAt); err != nil {
			return nil, err
		}
		files = append(files, &f)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}

// RemoveFile deletes a file from the database by its name.
//
// Parameters:
//   - ctx: The context for the operation.
//   - fileName: A string representing the name of the file to delete.
//
// Returns:
//   - An error if the operation fails.
func (p *PostgresStorage) RemoveFile(ctx context.Context, fileName string) error {
	_, err := p.Conn.ExecContext(ctx, "DELETE FROM files WHERE file_name = $1", fileName)
	return err
}

// GetUser retrieves a user by their login from the database.
//
// Parameters:
//   - ctx: The context for the operation.
//   - login: A string representing the user's login name.
//
// Returns:
//   - A pointer to a model.User instance containing the user information.
//   - An error if the operation fails or if the user is not found.
func (p *PostgresStorage) GetUser(ctx context.Context, login string) (*model.User, error) {
	row := p.Conn.QueryRowContext(ctx, "SELECT login, password, id FROM users WHERE login = $1", login)

	var u model.User
	if err := row.Scan(&u.Login, &u.Password, &u.ID); err != nil {
		return nil, err
	}
	return &u, nil
}

// AddUserCredentials adds new user credentials to the database.
//
// Parameters:
//   - ctx: The context for the operation.
//   - cred: A pointer to a model.Credentials instance containing the credentials to add.
//
// Returns:
//   - An error if the operation fails.
func (p *PostgresStorage) AddUserCredentials(ctx context.Context, cred *model.Credentials) error {
	_, err := p.Conn.ExecContext(
		ctx,
		"INSERT INTO user_credentials (login, password, description, user_id) VALUES ($1, $2, $3, $4)",
		cred.Login, cred.Password, cred.Description, cred.UserID)
	return err
}

// GetUserCredentials retrieves all credentials associated with a user.
//
// Parameters:
//   - ctx: The context for the operation.
//   - userID: An int64 representing the unique identifier of the user.
//
// Returns:
//   - A slice of pointers to model.Credentials instances containing the user's credentials.
//   - An error if the operation fails.
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

// GetUserCredential retrieves a specific user credential by its ID.
//
// Parameters:
//   - ctx: The context for the operation.
//   - id: An int64 representing the unique identifier of the credential.
//
// Returns:
//   - A pointer to a model.Credentials instance containing the credential information.
//   - An error if the operation fails or if the credential is not found.
func (p *PostgresStorage) GetUserCredential(ctx context.Context, id int64) (*model.Credentials, error) {
	row := p.Conn.QueryRowContext(ctx, "SELECT id, user_id, login, password, description FROM user_credentials WHERE id = $1", id)

	var cred model.Credentials
	if err := row.Scan(&cred.ID, &cred.UserID, &cred.Login, &cred.Password, &cred.Description); err != nil {
		return nil, err
	}

	return &cred, nil
}

// RemoveUserCredential deletes a user credential from the database by its ID.
//
// Parameters:
//   - ctx: The context for the operation.
//   - id: An int64 representing the unique identifier of the credential to delete.
//
// Returns:
//   - An error if the operation fails.
func (p *PostgresStorage) RemoveUserCredential(ctx context.Context, id int64) error {
	_, err := p.Conn.ExecContext(ctx, "DELETE FROM user_credentials WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

// AddNote adds a new note to the database.
//
// Parameters:
//   - ctx: The context for the operation.
//   - note: A pointer to a model.Note instance containing the note information to add.
//
// Returns:
//   - An error if the operation fails.
func (p *PostgresStorage) AddNote(ctx context.Context, note *model.Note) error {
	_, err := p.Conn.ExecContext(
		ctx,
		"INSERT INTO notes (text, description, user_id) VALUES ($1, $2, $3)",
		note.Text, note.Description, note.UserID)
	return err
}

// GetNotes retrieves all notes associated with a user.
//
// Parameters:
//   - ctx: The context for the operation.
//   - userID: An int64 representing the unique identifier of the user.
//
// Returns:
//   - A slice of pointers to model.Note instances containing the user's notes.
//   - An error if the operation fails.
func (p *PostgresStorage) GetNotes(ctx context.Context, userID int64) ([]*model.Note, error) {
	rows, err := p.Conn.QueryContext(ctx, "SELECT id, user_id, text, description FROM notes WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			logger.Log.Error("error close rows", zap.Error(err))
		}
	}(rows)

	var notes []*model.Note
	for rows.Next() {
		var n model.Note
		if err = rows.Scan(&n.ID, &n.UserID, &n.Text, &n.Description); err != nil {
			return nil, err
		}
		notes = append(notes, &n)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}

// GetNote retrieves a specific note by its ID from the database.
//
// Parameters:
//   - ctx: The context for the operation
//   - id: An int64 representing the unique identifier of the note to retrieve.
//
// Returns:
//   - A pointer to a model.Note instance containing the note information.
//   - An error if the operation fails or if the note is not found.
func (p *PostgresStorage) GetNote(ctx context.Context, id int64) (*model.Note, error) {
	row := p.Conn.QueryRowContext(ctx, "SELECT id, user_id, text, description FROM notes WHERE id = $1", id)

	var note model.Note
	if err := row.Scan(&note.ID, &note.UserID, &note.Text, &note.Description); err != nil {
		return nil, err
	}

	return &note, nil
}

// RemoveNote deletes a note from the database by its ID.
//
// Parameters:
//   - ctx: The context for the operation.
//   - id: An int64 representing the unique identifier of the note to delete.
//
// Returns:
//   - An error if the operation fails.
func (p *PostgresStorage) RemoveNote(ctx context.Context, id int64) error {
	_, err := p.Conn.ExecContext(ctx, "DELETE FROM notes WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

// AddCard adds a new bank card to the database.
//
// Parameters:
//   - ctx: The context for the operation.
//   - card: A pointer to a model.BankCard instance containing the card information to add.
//
// Returns:
//   - An error if the operation fails.
func (p *PostgresStorage) AddCard(ctx context.Context, card *model.BankCard) error {
	_, err := p.Conn.ExecContext(
		ctx,
		"INSERT INTO bank_cards (user_id, card_number, expiration_date, cvv, owner, description) "+
			"VALUES ($1, $2, $3, $4, $5, $6)", card.UserID, card.Number, card.ExpireDate, card.CVV, card.Owner, card.Description)
	return err
}

// GetBankCards retrieves all bank cards associated with a user.
//
// Parameters:
//   - ctx: The context for the operation.
//   - userID: An int64 representing the unique identifier of the user.
//
// Returns:
//   - A slice of pointers to model.BankCard instances containing the user's bank cards.
//   - An error if the operation fails.
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

// GetBankCard retrieves a specific bank card by its ID from the database.
//
// Parameters:
//   - ctx: The context for the operation.
//   - id: An int64 representing the unique identifier of the bank card to retrieve.
//
// Returns:
//   - A pointer to a model.BankCard instance containing the bank card information.
//   - An error if the operation fails or if the bank card is not found.
func (p *PostgresStorage) GetBankCard(ctx context.Context, id int64) (*model.BankCard, error) {
	row := p.Conn.QueryRowContext(ctx, "SELECT id, user_id, owner, card_number, expiration_date, cvv, description FROM bank_cards WHERE id = $1", id)

	var card model.BankCard
	if err := row.Scan(&card.ID, &card.UserID, &card.Owner, &card.Number, &card.ExpireDate, &card.CVV, &card.Description); err != nil {
		return nil, err
	}

	return &card, nil
}

// RemoveBankCard deletes a bank card from the database by its ID.
//
// Parameters:
//   - ctx: The context for the operation.
//   - id: An int64 representing the unique identifier of the bank card to delete.
//
// Returns:
//   - An error if the operation fails.
func (p *PostgresStorage) RemoveBankCard(ctx context.Context, id int64) error {
	_, err := p.Conn.ExecContext(ctx, "DELETE FROM bank_cards WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
