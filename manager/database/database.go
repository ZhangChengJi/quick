package database

import (
	"database/sql"
	"gorm.io/gorm"
	"os"
	"quick/conf"
	"quick/manager/model"
	mysql "quick/pkg/gorm"
	"quick/pkg/log"
	cache "quick/pkg/redis"
	"quick/pkg/tdengine"
)

type Database struct {
	DB  *gorm.DB
	RDB *cache.RedisClient
	TDB *sql.DB
}

func New() *Database {
	db, err := mysql.New(conf.MysqlConfig)
	if err != nil {
		log.Sugar.Errorf("数据库连接失败%s", err)
		os.Exit(1)
		return nil
	}

	rdb := cache.NewClient(conf.RedisConfig)
	tdb, err := tdengine.New(conf.TdengineConfig)
	if err != nil {
		log.Sugar.Errorf("tdengine连接失败%s", err)
		os.Exit(1)
		return nil
	}
	return &Database{
		DB:  db,
		RDB: rdb,
		TDB: tdb,
	}
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.PigDevice{},
	)

}
