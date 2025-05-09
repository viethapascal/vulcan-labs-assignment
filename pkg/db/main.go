package db

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"go.uber.org/zap"
	"vulcanlabs-assignment/pkg/utils"
)

var (
	logger *zap.SugaredLogger
	db     *leveldb.DB
)

type Queryable interface {
	Bytes() ([]byte, error)
	FromBytes([]byte) error
}

func InitDB() error {
	var err error
	logger = utils.Logger("app-db")
	if db, err = leveldb.OpenFile("./app-db", nil); err != nil {
		logger.Errorw("Failed to open database", "error", err)
		return err
	}
	return nil
}

func GetKey(key string, queryable Queryable) error {
	k := fmt.Sprintf("%v", key)
	dat, err := db.Get([]byte(k), nil)
	if err != nil {
		return err
	}
	logger.Debugw("Get value", "key", k, "dat", string(dat))
	return queryable.FromBytes(dat)
}

func SetKey(key string, value Queryable) error {
	k := fmt.Sprintf("%v", key)
	data, err := value.Bytes()
	if err != nil {
		logger.Errorw("Failed to marshal value", "key", k, "value", value, "error", err)
		return err
	}
	logger.Debugw("Set value", "key", k, "value", string(data))
	return db.Put([]byte(k), data, nil)
}
