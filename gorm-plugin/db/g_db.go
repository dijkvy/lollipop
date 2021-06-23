package db

import (
	"context"
	"fmt"
	"github.com/laxiaohong/lollipop/gorm-plugin/db/config"
	"github.com/laxiaohong/lollipop/gorm-plugin/log"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"time"
)

type DB struct {
	db *gorm.DB
}

func NewDB(cfg *config.MySQLConfig, logger *zap.Logger) (*DB, func()) {
	fmt.Printf("NewDB config:%#v", cfg.String())
	db, err := gorm.Open(mysql.Open(cfg.GetDsn()), &gorm.Config{
		Logger: log.NewLogger(logger, gormLogger.Config{
			SlowThreshold: time.Duration(cfg.LogConfig.GetSlowLogSecond()) * time.Second,
			LogLevel:      gormLogger.Info,
		}),
	})
	if err != nil {
		panic(err)
	}

	_db, err := db.DB()
	if err != nil {
		panic(err)
	}
	if cfg.GetConnConfig().MaxOpenConn != nil {
		_db.SetMaxOpenConns(int(cfg.GetConnConfig().GetMaxOpenConn()))
	} else {
		_db.SetMaxOpenConns(100)
	}

	if cfg.GetConnConfig().ConnMaxLifeSecond != nil {
		_db.SetConnMaxLifetime(time.Duration(cfg.GetConnConfig().GetConnMaxIdleSecond()) * time.Second)
	} else {
		_db.SetConnMaxLifetime(3600 * time.Second)
	}

	if cfg.GetConnConfig().ConnMaxIdleSecond != nil {
		_db.SetMaxIdleConns(int(cfg.GetConnConfig().GetMaxIdleConn()))
	} else {
		_db.SetMaxIdleConns(100)
	}

	if cfg.GetConnConfig().ConnMaxIdleSecond != nil {
		_db.SetConnMaxIdleTime(time.Duration(cfg.GetConnConfig().GetConnMaxIdleSecond()) * time.Second)
	} else {
		_db.SetConnMaxIdleTime(3600 * time.Second)

	}
	return &DB{db: db}, func() {
		err := _db.Close()
		fmt.Println("mysql close , err info ", err)
	}
}

func (c *DB) WithContext(ctx context.Context) *gorm.DB {

	return c.db.WithContext(ctx)
}
