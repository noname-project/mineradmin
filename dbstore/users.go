package dbstore

import (
	"github.com/boomstarternetwork/mineradmin/dbstore/model"
	"github.com/boomstarternetwork/mineradmin/store"
)

func (s DBStore) UserGet(userID uint64) (store.User, bool, error) {
	user := &model.User{}

	exists, err := s.xdb.ID(userID).Get(user)
	if err != nil {
		return store.User{}, false, err
	}
	if !exists {
		return store.User{}, false, nil
	}

	return store.User{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}, true, nil
}

func (s DBStore) UserAdd(userEmail string, userName string) (uint64, error) {
	user := &model.User{
		Email: userEmail,
		Name:  userName,
	}

	_, err := s.xdb.Insert(user)
	if err != nil {
		return 0, nil
	}

	return user.ID, nil
}

func (s DBStore) UserAddresses(userID uint64) (map[string][]string, error) {
	var modelAddrs []model.UserAddress

	err := s.xdb.Find(&modelAddrs, model.UserAddress{UserID: userID})
	if err != nil {
		return nil, err
	}

	addrs := map[string][]string{}

	for _, ma := range modelAddrs {
		addrs[ma.Coin] = append(addrs[ma.Coin], ma.Address)
	}

	return addrs, nil
}

func (s DBStore) UsersList() ([]store.User, error) {
	var modelUsers []model.User

	err := s.xdb.Find(&modelUsers)
	if err != nil {
		return nil, err
	}

	var users []store.User

	for _, mu := range modelUsers {
		users = append(users, store.User{
			ID:    mu.ID,
			Email: mu.Email,
			Name:  mu.Name,
		})
	}

	return users, nil
}

func (s DBStore) UserAddressAdd(userID uint64, coin string,
	address string) error {
	_, err := s.xdb.Insert(model.UserAddress{
		UserID:  userID,
		Coin:    coin,
		Address: address,
	})
	return err
}

func (s DBStore) UserAddressRemove(userID uint64, coin string, address string) error {
	_, err := s.xdb.Delete(model.UserAddress{
		UserID:  userID,
		Coin:    coin,
		Address: address,
	})
	return err
}
