package storage

import (
	"context"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go.elastic.co/apm/module/apmgorm"

	_ "go.elastic.co/apm/module/apmgorm/dialects/postgres"
)

type Storage struct {
	db *gorm.DB
}

func NewStorage(dbstring string) (*Storage, error) {
	driverName := "postgres"
	db, err := apmgorm.Open(driverName, dbstring)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to connect to postgres_custom '%s'", dbstring)
	}
	return &Storage{db}, nil
}

func (s Storage) GetDBConn() *gorm.DB {
	return s.db
}

func (s Storage) WithContext(ctx context.Context) *gorm.DB {
	s.db = apmgorm.WithContext(ctx, s.db)
	return s.db
}
