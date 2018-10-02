package store

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Test_dbProjectsStore_Get_success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub"+
			" database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow("id-1", "name-1")

	wantProject := Project{ID: "id-1", Name: "name-1"}

	mock.ExpectQuery(`^SELECT .+ FROM projects`).
		WithArgs("id-1").
		WillReturnRows(rows)

	ps := NewDBProjectsStore(db)

	project, err := ps.Get("id-1")

	if assert.NoError(t, err) {
		assert.Equal(t, wantProject, project)
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func Test_dbProjectsStore_Get_error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub"+
			" database connection", err)
	}
	defer db.Close()

	wantError := errors.New("some error")

	mock.ExpectQuery(`^SELECT (.+) FROM projects`).
		WithArgs("id-1").
		WillReturnError(wantError)

	ps := NewDBProjectsStore(db)

	_, err = ps.Get("id-1")

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
