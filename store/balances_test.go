package store

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Test_dbBalancesStore_ProjectsBalances_success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub"+
			" database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"projectid", "p.name", "coin", "SUM(amount)"}).
		AddRow("id-1", "name-1", "BTC", "1").
		AddRow("id-1", "name-1", "LTC", "2").
		AddRow("id-2", "name-2", "BTC", "3.100").
		AddRow("id-2", "name-2", "LTC", "40")

	wantBalances := []ProjectBalance{
		{
			ProjectID:   "id-1",
			ProjectName: "name-1",
			Coins: []CoinAmount{
				{Coin: "BTC", Amount: "1"},
				{Coin: "LTC", Amount: "2"},
			},
		},
		{
			ProjectID:   "id-2",
			ProjectName: "name-2",
			Coins: []CoinAmount{
				{Coin: "BTC", Amount: "3.1"},
				{Coin: "LTC", Amount: "40"},
			},
		},
	}

	mock.ExpectQuery(`^SELECT .+ FROM balances`).
		WillReturnRows(rows)

	bs := NewDBBalancesStore(db)

	balances, err := bs.ProjectsBalances()

	if assert.NoError(t, err) {
		assert.Equal(t, wantBalances, balances)
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func Test_dbBalancesStore_ProjectsBalances_error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub"+
			" database connection", err)
	}
	defer db.Close()

	wantError := errors.New("some error")

	mock.ExpectQuery(`^SELECT (.+) FROM balances`).
		WillReturnError(wantError)

	bs := NewDBBalancesStore(db)

	_, err = bs.ProjectsBalances()

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func Test_dbBalancesStore_ProjectUsersBalances_success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub"+
			" database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"address", "coin", "SUM(amount)"}).
		AddRow("addr-1", "BTC", "1").
		AddRow("addr-1", "LTC", "2").
		AddRow("addr-2", "BTC", "3.100").
		AddRow("addr-2", "LTC", "40")

	wantBalances := []UserBalance{
		{
			Address: "addr-1",
			Coins: []CoinAmount{
				{Coin: "BTC", Amount: "1"},
				{Coin: "LTC", Amount: "2"},
			},
		},
		{
			Address: "addr-2",
			Coins: []CoinAmount{
				{Coin: "BTC", Amount: "3.1"},
				{Coin: "LTC", Amount: "40"},
			},
		},
	}

	mock.ExpectQuery(`^SELECT .+ FROM balances`).
		WithArgs("id-1").
		WillReturnRows(rows)

	bs := NewDBBalancesStore(db)

	balances, err := bs.ProjectUsersBalances("id-1")

	if assert.NoError(t, err) {
		assert.Equal(t, wantBalances, balances)
	}

	assert.NoError(t, mock.ExpectationsWereMet())
}

func Test_dbBalancesStore_ProjectUsersBalances_error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub"+
			" database connection", err)
	}
	defer db.Close()

	wantError := errors.New("some error")

	mock.ExpectQuery(`^SELECT (.+) FROM balances`).
		WithArgs("id-1").
		WillReturnError(wantError)

	bs := NewDBBalancesStore(db)

	_, err = bs.ProjectUsersBalances("id-1")

	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
