package service

import (
	"log/slog"
	"os"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/gbadopt/gbadopt/internal/config"
	"github.com/gbadopt/gbadopt/internal/repository"
)

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
}

func newServiceDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock: %v", err)
	}
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("gorm open: %v", err)
	}
	return db, mock
}

func TestUserServiceLogin(t *testing.T) {
	db, mock := newServiceDB(t)
	cfg := config.Load()
	repo := repository.NewUserRepository(db)
	svc := NewUserService(repo, newTestLogger(), cfg)
	hash, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.MinCost)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs("alice", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password_hash", "nickname", "role"}).
			AddRow(1, "alice", "a@b.c", string(hash), "Alice", "user"))
	u, token, err := svc.Login("alice", "secret123")
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	if u.Username != "alice" || token == "" {
		t.Errorf("unexpected result: %+v", u)
	}
}

func TestUserServiceLoginWrongPassword(t *testing.T) {
	db, mock := newServiceDB(t)
	cfg := config.Load()
	repo := repository.NewUserRepository(db)
	svc := NewUserService(repo, newTestLogger(), cfg)
	hash, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.MinCost)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 ORDER BY "users"."id" LIMIT $2`)).
		WithArgs("alice", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password_hash", "nickname", "role"}).
			AddRow(1, "alice", "a@b.c", string(hash), "Alice", "user"))
	if _, _, err := svc.Login("alice", "wrong"); err == nil {
		t.Error("expected login error for wrong password")
	}
}
