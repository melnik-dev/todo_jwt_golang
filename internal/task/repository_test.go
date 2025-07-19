package task_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/melnik-dev/go_todo_jwt/internal/task"
	"github.com/melnik-dev/go_todo_jwt/internal/user"
	"github.com/melnik-dev/go_todo_jwt/pkg/db"
	"github.com/sirupsen/logrus"
	"io"
	"reflect"
	"regexp"
	"testing"
)

func mockDB() (*task.Repository, sqlmock.Sqlmock, error) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	pgDB := sqlx.NewDb(mockDb, "sqlMock")

	testLogger := logrus.New()
	testLogger.SetOutput(io.Discard)
	repo := task.NewRepository(&db.Db{
		DB: pgDB,
	}, testLogger)

	return repo, mock, err
}

func TestTaskRepository_Create_Success(t *testing.T) {
	repo, mock, err := mockDB()
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO tasks (user_id, title, description)
				VALUES ($1, $2, $3) 
				RETURNING id`)).
		WithArgs(42, "test_title", "test_desc").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	exp, err := repo.Create(&task.Task{
		UserID:      42,
		Title:       "test_title",
		Description: "test_desc",
	})
	if err != nil {
		t.Fatal(err)
	}

	if exp.ID != 1 {
		t.Errorf("Expected ID %d, got %d", 1, exp.ID)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestTaskRepository_Create_Fail(t *testing.T) {
	repo, mock, err := mockDB()
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectQuery(`INSERT INTO tasks`).
		WithArgs(42, "test_title", "test_desc").
		WillReturnError(sqlmock.ErrCancelled)

	_, err = repo.Create(&task.Task{
		UserID:      42,
		Title:       "test_title",
		Description: "test_desc",
	})
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestTaskRepository_Update_Success(t *testing.T) {
	repo, mock, err := mockDB()
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE tasks 
				SET title = $1, description = $2, completed = $3 
				WHERE id = $4 AND user_id = $5`)).
		WithArgs("test_title", "test_desc", true, 1, 42).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Update(&task.Task{
		ID:          1,
		UserID:      42,
		Title:       "test_title",
		Description: "test_desc",
		Completed:   true,
	})
	if err != nil {
		t.Fatal(err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestTaskRepository_Get_Fail(t *testing.T) {
	repo, mock, err := mockDB()
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectExec(`UPDATE tasks`).
		WithArgs("test_title", "test_desc", true, 1, 42).
		WillReturnError(sqlmock.ErrCancelled)

	err = repo.Update(&task.Task{
		ID:          1,
		UserID:      42,
		Title:       "test_title",
		Description: "test_desc",
		Completed:   true,
	})
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestTaskRepository_DeleteById_Success(t *testing.T) {
	repo, mock, err := mockDB()
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM tasks WHERE id = $1 AND user_id = $2`)).
		WithArgs(1, 42).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.DeleteById(&task.Task{
		ID:     1,
		UserID: 42,
	})
	if err != nil {
		t.Fatal(err)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestTaskRepository_DeleteById_Fail(t *testing.T) {
	repo, mock, err := mockDB()
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectExec(`DELETE FROM tasks`).
		WithArgs(1, 42).
		WillReturnError(sqlmock.ErrCancelled)

	err = repo.DeleteById(&task.Task{
		ID:     1,
		UserID: 42,
	})
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestTaskRepository_GetById_Success(t *testing.T) {
	repo, mock, err := mockDB()
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM tasks WHERE id = $1 AND user_id = $2`)).
		WithArgs(1, 42).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "description", "completed"}).
			AddRow(1, 42, "test_title", "test_desc", true))

	exp, err := repo.GetById(&task.Task{
		ID:          1,
		UserID:      42,
		Title:       "test_title",
		Description: "test_desc",
		Completed:   true,
	})
	if err != nil {
		t.Fatal(err)
	}

	if exp.ID != 1 ||
		exp.UserID != 42 ||
		exp.Title != "test_title" ||
		exp.Description != "test_desc" ||
		exp.Completed != true {
		t.Errorf("Unexpected result: %+v", exp)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestTaskRepository_GetById_Fail(t *testing.T) {
	repo, mock, err := mockDB()
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM tasks`)).
		WithArgs(1, 42).
		WillReturnError(sqlmock.ErrCancelled)

	_, err = repo.GetById(&task.Task{
		ID:     1,
		UserID: 42,
	})
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestTaskRepository_GetAll_Success(t *testing.T) {
	repo, mock, err := mockDB()
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM tasks WHERE user_id = $1`)).
		WithArgs(42).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "description", "completed"}).
			AddRow(1, 42, "test_title_1", "test_desc_1", true).
			AddRow(2, 42, "test_title_2", "test_desc_2", false))

	exp, err := repo.GetAll(&user.User{
		ID: 42,
	})
	if err != nil {
		t.Fatal(err)
	}

	expect := []task.Task{
		{ID: 1, UserID: 42, Title: "test_title_1", Description: "test_desc_1", Completed: true},
		{ID: 2, UserID: 42, Title: "test_title_2", Description: "test_desc_2", Completed: false},
	}

	if !reflect.DeepEqual(exp, expect) {
		t.Errorf("Expected %+v, got %+v", expect, exp)
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestTaskRepository_GetAll_Fail(t *testing.T) {
	repo, mock, err := mockDB()
	if err != nil {
		t.Fatal(err)
	}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM tasks WHERE user_id = $1`)).
		WithArgs(42).
		WillReturnError(sqlmock.ErrCancelled)

	_, err = repo.GetAll(&user.User{
		ID: 42,
	})
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
