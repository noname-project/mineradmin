package mockstore

import (
	"github.com/boomstarternetwork/mineradmin/store"
	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock
}

func New() *MockStore {
	return &MockStore{}
}

func (s *MockStore) ProjectGet(projectID string) (store.Project, bool, error) {
	args := s.Called(projectID)
	return args.Get(0).(store.Project), args.Get(1).(bool), args.Error(2)
}

func (s *MockStore) ProjectAdd(projectID string, projectName string) error {
	args := s.Called(projectID, projectName)
	return args.Error(0)
}

func (s *MockStore) ProjectSet(projectID string, newProjectName string) error {
	args := s.Called(projectID, newProjectName)
	return args.Error(0)
}

func (s *MockStore) ProjectRemove(projectID string) error {
	args := s.Called(projectID)
	return args.Error(0)
}

func (s *MockStore) ProjectsBalances() ([]store.ProjectBalance, error) {
	args := s.Called()
	projects := args.Get(0)
	if projects == nil {
		return nil, args.Error(1)
	}
	return projects.([]store.ProjectBalance), args.Error(1)
}

func (s *MockStore) ProjectUsersBalances(projectID string) (
	[]store.UserBalance, error) {
	args := s.Called(projectID)
	projects := args.Get(0)
	if projects == nil {
		return nil, args.Error(1)
	}
	return projects.([]store.UserBalance), args.Error(1)
}

func (s *MockStore) AdminAdd(adminLogin string) (string, error) {
	args := s.Called(adminLogin)
	return args.Get(0).(string), args.Error(1)
}

func (s *MockStore) AdminRemove(adminID uint64) error {
	args := s.Called(adminID)
	return args.Error(0)
}

func (s *MockStore) AdminCheckPassword(adminLogin string,
	adminPassword string) (bool, error) {
	args := s.Called(adminLogin, adminPassword)
	return args.Get(0).(bool), args.Error(1)
}

func (s *MockStore) AdminResetPassword(adminID uint64) (string, error) {
	args := s.Called(adminID)
	return args.Get(0).(string), args.Error(1)
}

func (s *MockStore) AdminsList() ([]store.Admin, error) {
	args := s.Called()
	admins := args.Get(0)
	if admins == nil {
		return nil, args.Error(1)
	}
	return admins.([]store.Admin), args.Error(1)
}

func (s *MockStore) UserGet(userID uint64) (store.User, bool, error) {
	args := s.Called(userID)
	return args.Get(0).(store.User), args.Get(1).(bool), args.Error(2)
}

func (s *MockStore) UserAdd(userEmail string, userName string) (uint64, error) {
	args := s.Called(userEmail, userName)
	return args.Get(0).(uint64), args.Error(1)
}

func (s *MockStore) UserAddresses(userID uint64) (map[string][]string, error) {
	args := s.Called(userID)
	addrs := args.Get(0)
	if addrs == nil {
		return nil, args.Error(1)
	}
	return addrs.(map[string][]string), args.Error(1)
}

func (s *MockStore) UsersList() ([]store.User, error) {
	args := s.Called()
	users := args.Get(0)
	if users == nil {
		return nil, args.Error(1)
	}
	return users.([]store.User), args.Error(1)
}

func (s *MockStore) UserAddressAdd(userID uint64, coin string,
	address string) error {
	args := s.Called(userID, coin, address)
	return args.Error(0)
}

func (s *MockStore) UserAddressRemove(userID uint64, coin string,
	address string) error {
	args := s.Called(userID, coin, address)
	return args.Error(0)
}
