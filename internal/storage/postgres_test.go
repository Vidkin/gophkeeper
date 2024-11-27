package storage

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Vidkin/gophkeeper/internal/model"
)

func setupTestDB(t *testing.T) (*PostgresStorage, string) {
	connStr := "user=postgres password=postgres dbname=postgres host=127.0.0.1 port=5432 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)

	tempDBName := "test_db"
	_, err = db.Exec("SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = $1", tempDBName)
	require.NoError(t, err)
	db.Exec("DROP DATABASE " + tempDBName)
	_, err = db.Exec("CREATE DATABASE " + tempDBName)
	require.NoError(t, err)

	tempConnStr := "user=postgres password=postgres dbname=" + tempDBName + " host=127.0.0.1 port=5432 sslmode=disable"
	storage, err := NewPostgresStorage(tempConnStr)
	require.NoError(t, err)

	return storage, tempDBName
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

func TestPostgresStorage_Ping(t *testing.T) {
	db, dbName := setupTestDB(t)
	defer teardownTestDB(t, db.Conn, dbName)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name: "test ping ok",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := db.Ping(context.Background())
			assert.NoError(t, err)
		})
	}
}

func TestPostgresStorage_Close(t *testing.T) {
	db, dbName := setupTestDB(t)
	defer teardownTestDB(t, db.Conn, dbName)

	tests := []struct {
		name string
	}{
		{
			name: "test close ok",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := db.Close()
			assert.NoError(t, err)
		})
	}
}

func TestPostgresStorage_AddUser(t *testing.T) {
	db, dbName := setupTestDB(t)
	defer teardownTestDB(t, db.Conn, dbName)

	tests := []struct {
		name string
	}{
		{
			name: "test add user ok",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := db.AddUser(context.Background(), "login", "password")
			assert.NoError(t, err)
		})
	}
}

func TestPostgresStorage_AddFile(t *testing.T) {
	db, dbName := setupTestDB(t)
	defer teardownTestDB(t, db.Conn, dbName)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "test add file error",
			wantErr: true,
		},
		{
			name:    "test add and update file ok",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				err := db.AddFile(context.Background(), "bucketName", "fileName", "description", 1, 12)
				assert.Error(t, err)
			} else {
				err := db.AddUser(context.Background(), "login", "password")
				require.NoError(t, err)
				err = db.AddFile(context.Background(), "bucketName", "fileName", "description", 1, 12)
				assert.NoError(t, err)
				err = db.AddFile(context.Background(), "bucketName", "fileName", "description", 1, 12)
				assert.NoError(t, err)
			}
		})
	}
}

func TestPostgresStorage_GetFile(t *testing.T) {
	db, dbName := setupTestDB(t)
	defer teardownTestDB(t, db.Conn, dbName)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "test get file error",
			wantErr: true,
		},
		{
			name:    "test get file ok",
			wantErr: false,
		},
	}

	err := db.AddUser(context.Background(), "login", "password")
	require.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				err = db.AddFile(context.Background(), "bucketName", "fileName", "description", 1, 12)
				_, err = db.GetFile(context.Background(), "badName")
				assert.Error(t, err)
			} else {
				err = db.AddFile(context.Background(), "bucketName", "goodFileName", "description", 1, 12)
				assert.NoError(t, err)
				file, err := db.GetFile(context.Background(), "goodFileName")
				assert.NoError(t, err)
				assert.Equal(t, "goodFileName", file.FileName)
				assert.Equal(t, "bucketName", file.BucketName)
				assert.Equal(t, "description", file.Description)
				assert.Equal(t, int64(12), file.FileSize)
			}
		})
	}
}

func TestPostgresStorage_GetFiles(t *testing.T) {
	db, dbName := setupTestDB(t)
	defer teardownTestDB(t, db.Conn, dbName)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "test get files error",
			wantErr: true,
		},
		{
			name:    "test get files ok",
			wantErr: false,
		},
	}

	err := db.AddUser(context.Background(), "login", "password")
	require.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				creds, err := db.GetFiles(context.Background(), 2)
				assert.NoError(t, err)
				assert.Empty(t, creds)
			} else {
				err = db.AddFile(context.Background(), "bucketName", "goodFileName", "description", 1, 12)
				assert.NoError(t, err)
				file, err := db.GetFiles(context.Background(), 1)
				assert.NoError(t, err)
				assert.Equal(t, "goodFileName", file[0].FileName)
				assert.Equal(t, "bucketName", file[0].BucketName)
				assert.Equal(t, "description", file[0].Description)
				assert.Equal(t, int64(12), file[0].FileSize)
			}
		})
	}
}

func TestPostgresStorage_RemoveFile(t *testing.T) {
	db, dbName := setupTestDB(t)
	defer teardownTestDB(t, db.Conn, dbName)

	tests := []struct {
		name string
	}{
		{
			name: "test remove file ok",
		},
	}

	err := db.AddUser(context.Background(), "login", "password")
	require.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = db.AddFile(context.Background(), "bucketName", "goodFileName", "description", 1, 12)
			assert.NoError(t, err)
			err = db.RemoveFile(context.Background(), "goodFileName")
			assert.NoError(t, err)
			_, err = db.GetFile(context.Background(), "goodFileName")
			assert.Equal(t, "sql: no rows in result set", err.Error())
		})
	}
}

func TestPostgresStorage_GetUser(t *testing.T) {
	db, dbName := setupTestDB(t)
	defer teardownTestDB(t, db.Conn, dbName)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "test get user error",
			wantErr: true,
		},
		{
			name:    "test get user ok",
			wantErr: false,
		},
	}

	err := db.AddUser(context.Background(), "login", "password")
	require.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				_, err = db.GetUser(context.Background(), "badLogin")
				assert.Error(t, err)
			} else {
				user, err := db.GetUser(context.Background(), "login")
				assert.NoError(t, err)
				assert.Equal(t, "login", user.Login)
				assert.Equal(t, "password", user.Password)
				assert.Equal(t, int64(1), user.ID)
			}
		})
	}
}

func TestPostgresStorage_AddUserCredentials(t *testing.T) {
	db, dbName := setupTestDB(t)
	defer teardownTestDB(t, db.Conn, dbName)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "test add user credentials error",
			wantErr: true,
		},
		{
			name:    "test add user credentials ok",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				err := db.AddUserCredentials(context.Background(), &model.Credentials{
					Login:       "login",
					Password:    "password",
					Description: "description",
					UserID:      1,
				})
				assert.Error(t, err)
			} else {
				err := db.AddUser(context.Background(), "login", "password")
				require.NoError(t, err)
				err = db.AddUserCredentials(context.Background(), &model.Credentials{
					Login:       "login",
					Password:    "password",
					Description: "description",
					UserID:      1,
				})
				assert.NoError(t, err)
			}
		})
	}
}

func TestPostgresStorage_GetUserCredential(t *testing.T) {
	db, dbName := setupTestDB(t)
	defer teardownTestDB(t, db.Conn, dbName)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "test get user credential error",
			wantErr: true,
		},
		{
			name:    "test get user credential ok",
			wantErr: false,
		},
	}

	err := db.AddUser(context.Background(), "login", "password")
	require.NoError(t, err)
	err = db.AddUserCredentials(context.Background(), &model.Credentials{
		Login:       "login",
		Password:    "password",
		Description: "description",
		UserID:      1,
	})
	assert.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				_, err = db.GetUserCredential(context.Background(), 2)
				assert.Error(t, err)
			} else {
				cred, err := db.GetUserCredential(context.Background(), 1)
				assert.NoError(t, err)
				assert.Equal(t, "login", cred.Login)
				assert.Equal(t, "password", cred.Password)
				assert.Equal(t, "description", cred.Description)
				assert.Equal(t, int64(1), cred.UserID)
			}
		})
	}
}

func TestPostgresStorage_GetUserCredentials(t *testing.T) {
	db, dbName := setupTestDB(t)
	defer teardownTestDB(t, db.Conn, dbName)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "test get user credentials error",
			wantErr: true,
		},
		{
			name:    "test get user credentials ok",
			wantErr: false,
		},
	}

	err := db.AddUser(context.Background(), "login", "password")
	require.NoError(t, err)
	err = db.AddUserCredentials(context.Background(), &model.Credentials{
		Login:       "login",
		Password:    "password",
		Description: "description",
		UserID:      1,
	})
	assert.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				creds, err := db.GetUserCredentials(context.Background(), 2)
				assert.NoError(t, err)
				assert.Empty(t, creds)
			} else {
				creds, err := db.GetUserCredentials(context.Background(), 1)
				assert.NoError(t, err)
				assert.Equal(t, "login", creds[0].Login)
				assert.Equal(t, "password", creds[0].Password)
				assert.Equal(t, "description", creds[0].Description)
				assert.Equal(t, int64(1), creds[0].UserID)
			}
		})
	}
}

func TestPostgresStorage_RemoveUserCredential(t *testing.T) {
	db, dbName := setupTestDB(t)
	defer teardownTestDB(t, db.Conn, dbName)

	tests := []struct {
		name string
	}{
		{
			name: "test remove user credential ok",
		},
	}

	err := db.AddUser(context.Background(), "login", "password")
	require.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = db.AddUserCredentials(context.Background(), &model.Credentials{
				Login:       "login",
				Password:    "password",
				Description: "description",
				UserID:      1,
			})
			assert.NoError(t, err)
			err = db.RemoveUserCredential(context.Background(), 1)
			assert.NoError(t, err)
			_, err = db.GetUserCredential(context.Background(), 1)
			assert.Equal(t, "sql: no rows in result set", err.Error())
		})
	}
}

func TestPostgresStorage_AddNote(t *testing.T) {
	db, dbName := setupTestDB(t)
	defer teardownTestDB(t, db.Conn, dbName)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "test add note error",
			wantErr: true,
		},
		{
			name:    "test add note ok",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				err := db.AddNote(context.Background(), &model.Note{
					Text:        "test",
					Description: "description",
					UserID:      1,
				})
				assert.Error(t, err)
			} else {
				err := db.AddUser(context.Background(), "login", "password")
				require.NoError(t, err)
				err = db.AddNote(context.Background(), &model.Note{
					Text:        "test",
					Description: "description",
					UserID:      1,
				})
				assert.NoError(t, err)
			}
		})
	}
}

func TestPostgresStorage_GetNotes(t *testing.T) {
	db, dbName := setupTestDB(t)
	defer teardownTestDB(t, db.Conn, dbName)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "test get notes error",
			wantErr: true,
		},
		{
			name:    "test get notes ok",
			wantErr: false,
		},
	}

	err := db.AddUser(context.Background(), "login", "password")
	require.NoError(t, err)
	err = db.AddNote(context.Background(), &model.Note{
		Text:        "test",
		Description: "description",
		UserID:      1,
	})
	assert.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				notes, err := db.GetNotes(context.Background(), 2)
				assert.NoError(t, err)
				assert.Empty(t, notes)
			} else {
				notes, err := db.GetNotes(context.Background(), 1)
				assert.NoError(t, err)
				assert.Equal(t, "test", notes[0].Text)
				assert.Equal(t, "description", notes[0].Description)
				assert.Equal(t, int64(1), notes[0].UserID)
			}
		})
	}
}

func TestPostgresStorage_GetNote(t *testing.T) {
	db, dbName := setupTestDB(t)
	defer teardownTestDB(t, db.Conn, dbName)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "test get note error",
			wantErr: true,
		},
		{
			name:    "test get note ok",
			wantErr: false,
		},
	}

	err := db.AddUser(context.Background(), "login", "password")
	require.NoError(t, err)
	err = db.AddNote(context.Background(), &model.Note{
		Text:        "test",
		Description: "description",
		UserID:      1,
	})
	assert.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				_, err = db.GetNote(context.Background(), 2)
				assert.Error(t, err)
			} else {
				note, err := db.GetNote(context.Background(), 1)
				assert.NoError(t, err)
				assert.Equal(t, "test", note.Text)
				assert.Equal(t, "description", note.Description)
				assert.Equal(t, int64(1), note.UserID)
			}
		})
	}
}

func TestPostgresStorage_RemoveNote(t *testing.T) {
	db, dbName := setupTestDB(t)
	defer teardownTestDB(t, db.Conn, dbName)

	tests := []struct {
		name string
	}{
		{
			name: "test remove note ok",
		},
	}

	err := db.AddUser(context.Background(), "login", "password")
	require.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = db.AddNote(context.Background(), &model.Note{
				Text:        "test",
				Description: "description",
				UserID:      1,
			})
			assert.NoError(t, err)
			err = db.RemoveNote(context.Background(), 1)
			assert.NoError(t, err)
			_, err = db.GetNote(context.Background(), 1)
			assert.Equal(t, "sql: no rows in result set", err.Error())
		})
	}
}

func TestPostgresStorage_AddCard(t *testing.T) {
	db, dbName := setupTestDB(t)
	defer teardownTestDB(t, db.Conn, dbName)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "test add card error",
			wantErr: true,
		},
		{
			name:    "test add card ok",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				err := db.AddCard(context.Background(), &model.BankCard{
					ExpireDate:  "12.02.2024",
					Owner:       "owner",
					CVV:         "123",
					Number:      "1234",
					Description: "desc",
					UserID:      1,
				})
				assert.Error(t, err)
			} else {
				err := db.AddUser(context.Background(), "login", "password")
				require.NoError(t, err)
				err = db.AddCard(context.Background(), &model.BankCard{
					ExpireDate:  "12.02.2024",
					Owner:       "owner",
					CVV:         "123",
					Number:      "1234",
					Description: "desc",
					UserID:      1,
				})
				assert.NoError(t, err)
			}
		})
	}
}

func TestPostgresStorage_GetCards(t *testing.T) {
	db, dbName := setupTestDB(t)
	defer teardownTestDB(t, db.Conn, dbName)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "test get cards error",
			wantErr: true,
		},
		{
			name:    "test get cards ok",
			wantErr: false,
		},
	}

	err := db.AddUser(context.Background(), "login", "password")
	require.NoError(t, err)
	err = db.AddCard(context.Background(), &model.BankCard{
		ExpireDate:  "12.02.2024",
		Owner:       "owner",
		CVV:         "123",
		Number:      "1234",
		Description: "desc",
		UserID:      1,
	})
	assert.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				cards, err := db.GetBankCards(context.Background(), 2)
				assert.NoError(t, err)
				assert.Empty(t, cards)
			} else {
				cards, err := db.GetBankCards(context.Background(), 1)
				assert.NoError(t, err)
				assert.Equal(t, "12.02.2024", cards[0].ExpireDate)
				assert.Equal(t, "owner", cards[0].Owner)
				assert.Equal(t, int64(1), cards[0].UserID)
				assert.Equal(t, "desc", cards[0].Description)
				assert.Equal(t, "123", cards[0].CVV)
				assert.Equal(t, "1234", cards[0].Number)
			}
		})
	}
}

func TestPostgresStorage_GetBankCard(t *testing.T) {
	db, dbName := setupTestDB(t)
	defer teardownTestDB(t, db.Conn, dbName)

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "test get card error",
			wantErr: true,
		},
		{
			name:    "test get card ok",
			wantErr: false,
		},
	}

	err := db.AddUser(context.Background(), "login", "password")
	require.NoError(t, err)
	err = db.AddCard(context.Background(), &model.BankCard{
		ExpireDate:  "12.02.2024",
		Owner:       "owner",
		CVV:         "123",
		Number:      "1234",
		Description: "desc",
		UserID:      1,
	})
	assert.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				_, err = db.GetBankCard(context.Background(), 2)
				assert.Error(t, err)
			} else {
				card, err := db.GetBankCard(context.Background(), 1)
				assert.NoError(t, err)
				assert.Equal(t, "12.02.2024", card.ExpireDate)
				assert.Equal(t, "owner", card.Owner)
				assert.Equal(t, int64(1), card.UserID)
				assert.Equal(t, "desc", card.Description)
				assert.Equal(t, "123", card.CVV)
				assert.Equal(t, "1234", card.Number)
			}
		})
	}
}

func TestPostgresStorage_RemoveBankCard(t *testing.T) {
	db, dbName := setupTestDB(t)
	defer teardownTestDB(t, db.Conn, dbName)

	tests := []struct {
		name string
	}{
		{
			name: "test remove card ok",
		},
	}

	err := db.AddUser(context.Background(), "login", "password")
	require.NoError(t, err)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = db.AddCard(context.Background(), &model.BankCard{
				ExpireDate:  "12.02.2024",
				Owner:       "owner",
				CVV:         "123",
				Number:      "1234",
				Description: "desc",
				UserID:      1,
			})
			assert.NoError(t, err)
			err = db.RemoveBankCard(context.Background(), 1)
			assert.NoError(t, err)
			_, err = db.GetBankCard(context.Background(), 1)
			assert.Equal(t, "sql: no rows in result set", err.Error())
		})
	}
}
