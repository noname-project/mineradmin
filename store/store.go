package store

type Store interface {
	ProjectGet(projectID string) (Project, bool, error)
	ProjectAdd(projectID string, projectName string) error
	ProjectSet(projectID string, newProjectName string) error
	ProjectRemove(projectID string) error

	ProjectsBalances() ([]ProjectBalance, error)
	ProjectUsersBalances(projectID string) ([]UserBalance, error)

	AdminAdd(adminLogin string) (string, error)
	AdminRemove(adminID uint64) error
	AdminsList() ([]Admin, error)

	AdminCheckPassword(adminLogin string, adminPassword string) (bool, error)
	AdminResetPassword(adminID uint64) (string, error)

	UserGet(userID uint64) (User, bool, error)
	UserAdd(userEmail string, userName string) (uint64, error)
	UserAddresses(userID uint64) (map[string][]string, error)
	UsersList() ([]User, error)

	UserAddressAdd(userID uint64, coin string, address string) error
	UserAddressRemove(userID uint64, coin string, address string) error
}

type Admin struct {
	ID    uint64
	Login string
}

type CoinAmount struct {
	Coin   string
	Amount string
}

type ProjectBalance struct {
	ProjectID   string
	ProjectName string
	Coins       []CoinAmount
}

type UserBalance struct {
	Email string
	Coins []CoinAmount
}

type Project struct {
	ID   string
	Name string
}

type User struct {
	ID    uint64
	Email string
	Name  string
}
