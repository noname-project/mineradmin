package dbstore

import (
	"errors"

	"github.com/boomstarternetwork/mineradmin/dbstore/model"
	"github.com/boomstarternetwork/mineradmin/store"
	"github.com/sethvargo/go-password/password"
	"golang.org/x/crypto/bcrypt"
)

func generatePassword() (string, error) {
	return password.Generate(12, 4, 4, false, true)
}

func (s DBStore) AdminAdd(adminLogin string) (string, error) {
	pswd, err := generatePassword()
	if err != nil {
		return "", errors.New("failed to generate password: " + err.Error())
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pswd), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("failed to hash password: " + err.Error())
	}

	admin := &model.Admin{
		Login:    adminLogin,
		Password: hash,
	}

	_, err = s.xdb.Insert(admin)

	return pswd, err
}

func (s DBStore) AdminRemove(adminID uint64) error {
	_, err := s.xdb.ID(adminID).Delete(model.Admin{})
	return err
}

func (s DBStore) AdminCheckPassword(adminLogin string,
	adminPassword string) (bool, error) {

	admin := &model.Admin{Login: adminLogin}

	exists, err := s.xdb.Get(admin)
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}

	err = bcrypt.CompareHashAndPassword(admin.Password, []byte(adminPassword))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

func (s DBStore) AdminResetPassword(adminID uint64) (string, error) {
	pswd, err := generatePassword()
	if err != nil {
		return "", errors.New("failed to generate password: " + err.Error())
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pswd), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("failed to hash password: " + err.Error())
	}

	admin := &model.Admin{Password: hash}

	s.xdb.ID(adminID).Update(admin)

	return pswd, err
}

func (s DBStore) AdminsList() ([]store.Admin, error) {
	var (
		admins      []store.Admin
		modelAdmins []model.Admin
	)

	err := s.xdb.OrderBy("login ASC").Find(&modelAdmins)
	if err != nil {
		return nil, err
	}

	for _, ma := range modelAdmins {
		admins = append(admins, store.Admin{
			ID:    ma.ID,
			Login: ma.Login,
		})
	}

	return admins, nil
}
