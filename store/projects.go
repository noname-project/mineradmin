package store

import (
	"database/sql"

	"github.com/stretchr/testify/mock"
)

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProjectsStore interface {
	Get(projectID string) (Project, error)
}

type dbProjectsStore struct {
	db *sql.DB
}

func NewDBProjectsStore(db *sql.DB) ProjectsStore {
	return dbProjectsStore{db: db}
}

func (ps dbProjectsStore) Get(projectID string) (Project, error) {
	var project Project
	err := ps.db.QueryRow("SELECT id, name FROM projects WHERE id=$1",
		projectID).Scan(&project.ID, &project.Name)
	return project, err
}

type MockProjectsStore struct {
	mock.Mock
}

func NewMockProjectsStore() *MockProjectsStore {
	return &MockProjectsStore{}
}

func (ps *MockProjectsStore) Get(projectID string) (Project, error) {
	args := ps.Called()
	return args.Get(0).(Project), args.Error(1)
}
