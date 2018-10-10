package dbstore

import (
	"github.com/go-xorm/xorm"
)

type DBStore struct {
	xdb *xorm.Engine
}

func New(xdb *xorm.Engine) DBStore {
	return DBStore{xdb: xdb}
}
