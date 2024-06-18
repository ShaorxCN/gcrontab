package model

import (
	"gcrontab/constant"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func setupMockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gormDB, err := gorm.Open("postgres", db)
	if err != nil {
		return nil, nil, err
	}

	return gormDB, mock, nil
}

func TestFindEmail(t *testing.T) {
	mockdb, mock, err := setupMockDB()
	db.DB = mockdb
	close(done)
	if err != nil {
		t.Fatalf("Failed to setup mock database: %v", err)
	}
	defer mockdb.Close()

	id := uuid.New().String()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT DISTINCT email FROM \"tbl_user\" WHERE (id= $1 and fail_notify = $2 and status = $3 or role = $4)")).WithArgs(id, constant.NOTIFYONDB, constant.STATUSNORMALDB, constant.TASKADMINDB).WillReturnRows(sqlmock.NewRows([]string{"email"}).AddRow("on@test.com"))

	_, err = FindEmails(id)

	if err != nil {
		t.Errorf("error was not expected while query email: %s", err)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %v", err)
	}
}
