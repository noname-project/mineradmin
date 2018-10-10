package dbstore

import (
	"database/sql"

	"github.com/boomstarternetwork/mineradmin/store"
	"github.com/shopspring/decimal"
)

func (s DBStore) ProjectsBalances() ([]store.ProjectBalance, error) {
	var balances []store.ProjectBalance

	rows, err := s.xdb.DB().Query(`
		SELECT p.id, p.name, bs.coin, SUM(bs.amount)
		FROM projects AS p LEFT JOIN balances AS bs ON p.id = bs.projectid
  		GROUP BY p.id, p.name, bs.coin ORDER BY p.name ASC, bs.coin ASC;
	`)
	if err != nil {
		return balances, err
	}
	defer rows.Close()

	var (
		projectID   string
		projectName string
		coin        sql.NullString
		amount      sql.NullString
	)

	for rows.Next() {
		err = rows.Scan(&projectID, &projectName, &coin, &amount)
		if err != nil {
			return balances, err
		}

		if !coin.Valid || !amount.Valid {
			// If balances for project is empty we can get null coin and amount.
			// In that case just add project id and name without coins.
			balances = append(balances, store.ProjectBalance{
				ProjectID:   projectID,
				ProjectName: projectName,
			})
			continue
		}

		// Use decimal to properly truncate trailing zeros from string.
		amountDec, err := decimal.NewFromString(amount.String)
		if err != nil {
			return balances, err
		}

		if len(balances) == 0 ||
			balances[len(balances)-1].ProjectID != projectID {
			// If this is first project or next project we initing its
			// in balances array.
			balances = append(balances, store.ProjectBalance{
				ProjectID:   projectID,
				ProjectName: projectName,
				Coins: []store.CoinAmount{
					{Coin: coin.String, Amount: amountDec.String()},
				},
			})
		} else {
			// Otherwise, we adding next coin data.
			balances[len(balances)-1].Coins = append(
				balances[len(balances)-1].Coins, store.CoinAmount{
					Coin: coin.String, Amount: amountDec.String()})
		}
	}

	return balances, rows.Err()
}

func (s DBStore) ProjectUsersBalances(projectID string) (
	[]store.UserBalance, error) {
	var balances []store.UserBalance

	rows, err := s.xdb.DB().Query(`
		SELECT u.email, b.coin, SUM(b.amount)
		FROM users AS u
			LEFT JOIN user_addresses AS ua ON u.id = ua.userid
			LEFT JOIN balances AS b ON ua.address = b.address
		WHERE b.projectid = $1
		GROUP BY u.email, b.coin
		ORDER BY u.email ASC, b.coin ASC;
	`, projectID)
	if err != nil {
		return balances, err
	}
	defer rows.Close()

	var (
		email  string
		coin   sql.NullString
		amount sql.NullString
	)

	for rows.Next() {
		err = rows.Scan(&email, &coin, &amount)
		if err != nil {
			return balances, err
		}

		if !coin.Valid || !amount.Valid {
			// If balances for user is empty we can get null coin and amount.
			// In that case just add user email without coins.
			balances = append(balances, store.UserBalance{
				Email: email,
			})
			continue
		}

		// Use decimal to properly truncate trailing zeros from string.
		amountDec, err := decimal.NewFromString(amount.String)
		if err != nil {
			return balances, err
		}

		if len(balances) == 0 || balances[len(balances)-1].Email != email {
			// If this is first user or next user we initing its
			// in balances array.
			balances = append(balances, store.UserBalance{
				Email: email,
				Coins: []store.CoinAmount{
					{Coin: coin.String, Amount: amountDec.String()},
				},
			})
		} else {
			// Otherwise, we adding next coin data.
			balances[len(balances)-1].Coins = append(
				balances[len(balances)-1].Coins,
				store.CoinAmount{Coin: coin.String, Amount: amountDec.String()})
		}
	}

	return balances, rows.Err()
}
