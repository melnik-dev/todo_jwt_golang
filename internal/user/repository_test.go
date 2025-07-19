package user_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/melnik-dev/go_todo_jwt/internal/user"
	"github.com/melnik-dev/go_todo_jwt/pkg/db"
	"github.com/sirupsen/logrus"
	"io"
	"regexp"
	"testing"
)

func mockDB() (*user.Repository, sqlmock.Sqlmock, error) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	pgDB := sqlx.NewDb(mockDb, "sqlMock")

	testLogger := logrus.New()
	testLogger.SetOutput(io.Discard)
	repo := user.NewRepository(&db.Db{
		DB: pgDB,
	}, testLogger)

	return repo, mock, err
}

func TestUserRepository_Create_Success(t *testing.T) {
	repo, mock, err := mockDB()
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`)).
		WithArgs("test_user", "test_pass").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	result, err := repo.Create(&user.User{
		Name:     "test_user",
		Password: "test_pass",
	})
	if err != nil {
		t.Fatal(err)
	}

	if result.ID != 1 {
		t.Errorf("Expected ID %d, got %d", 1, result.ID)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestUserRepository_Create_Fail(t *testing.T) {
	repo, mock, err := mockDB()
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs("test_user", "test_pass").
		WillReturnError(sqlmock.ErrCancelled)

	_, err = repo.Create(&user.User{
		Name:     "test_user",
		Password: "test_pass",
	})
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestUserRepository_Get_Success(t *testing.T) {
	repo, mock, err := mockDB()
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users WHERE username = $1`)).
		WithArgs("test_user").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password"}).
			AddRow(1, "test_user", "test_pass"))

	result, err := repo.Get("test_user")
	if err != nil {
		t.Fatal(err)
	}

	if result.ID != 1 || result.Name != "test_user" {
		t.Errorf("Unexpected result: %+v", result)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestUserRepository_Get_Fail(t *testing.T) {
	repo, mock, err := mockDB()
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM users WHERE username = $1`)).
		WithArgs("test_user").
		WillReturnError(sqlmock.ErrCancelled)

	_, err = repo.Get("test_user")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
