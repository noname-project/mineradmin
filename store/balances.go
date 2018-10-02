package store

import (
	"database/sql"

	"github.com/stretchr/testify/mock"
)

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
	Address string
	Coins   []CoinAmount
}

type BalancesStore interface {
	ProjectsBalances() ([]ProjectBalance, error)
	ProjectUsersBalances(projectID string) ([]UserBalance, error)
}

type dbBalancesStore struct {
	db *sql.DB
}

func NewDBBalancesStore(db *sql.DB) BalancesStore {
	return dbBalancesStore{db: db}
}

func (bs dbBalancesStore) ProjectsBalances() ([]ProjectBalance, error) {
	var balances []ProjectBalance

	rows, err := bs.db.Query(`
		SELECT projectid, p.name, coin, SUM(amount)
		FROM balances as b LEFT JOIN projects AS p ON b.projectid = p.id
  		GROUP BY projectid, p.name, coin ORDER BY p.name ASC;
	`)
	if err != nil {
		return balances, err
	}
	defer rows.Close()

	var (
		projectID   string
		projectName string
		coin        string
		amount      string
	)

	for rows.Next() {
		err = rows.Scan(&projectID, &projectName, &coin, &amount)
		if err != nil {
			return balances, err
		}

		if len(balances) == 0 || balances[len(balances)-1].
			ProjectID != projectID {
			balances = append(balances, ProjectBalance{
				ProjectID:   projectID,
				ProjectName: projectName,
				Coins: []CoinAmount{
					{Coin: coin, Amount: amount},
				},
			})
		} else {
			balances[len(balances)-1].Coins = append(balances[len(
				balances)-1].Coins, CoinAmount{Coin: coin, Amount: amount})
		}
	}

	return balances, rows.Err()
}

func (bs dbBalancesStore) ProjectUsersBalances(projectID string) (
	[]UserBalance, error) {
	var balances []UserBalance

	rows, err := bs.db.Query(`
		SELECT address, coin, SUM(amount)
		FROM balances WHERE projectid = $1
  		GROUP BY address, coin ORDER BY address ASC, coin ASC;
	`, projectID)
	if err != nil {
		return balances, err
	}
	defer rows.Close()

	var (
		address string
		coin    string
		amount  string
	)

	for rows.Next() {
		err = rows.Scan(&address, &coin, &amount)
		if err != nil {
			return balances, err
		}

		if len(balances) == 0 || balances[len(balances)-1].
			Address != address {
			balances = append(balances, UserBalance{
				Address: address,
				Coins: []CoinAmount{
					{Coin: coin, Amount: amount},
				},
			})
		} else {
			balances[len(balances)-1].Coins = append(balances[len(
				balances)-1].Coins, CoinAmount{Coin: coin, Amount: amount})
		}
	}

	return balances, rows.Err()
}

type MockBalancesStore struct {
	mock.Mock
}

func NewMockBalancesStore() *MockBalancesStore {
	return &MockBalancesStore{}
}

func (bs *MockBalancesStore) ProjectsBalances() ([]ProjectBalance, error) {
	args := bs.Called()
	projects := args.Get(0)
	if projects == nil {
		return nil, args.Error(1)
	}
	return projects.([]ProjectBalance), args.Error(1)
}

func (bs *MockBalancesStore) ProjectUsersBalances() ([]UserBalance, error) {
	args := bs.Called()
	projects := args.Get(0)
	if projects == nil {
		return nil, args.Error(1)
	}
	return projects.([]UserBalance), args.Error(1)
}
