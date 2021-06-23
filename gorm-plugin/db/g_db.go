package db

import (
	"context"
	"github.com/laxiaohong/lollipop/gorm-plugin/db/config"
	"github.com/laxiaohong/lollipop/gorm-plugin/log"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

type DB struct {
	db *gorm.DB
}

func NewDB(cfg *config.MySQLConfig, logger *zap.Logger) (*DB, func()) {
	const _prefixLog = "NewDB: "
	logrus.Infof("NewDB config:%#v", cfg.String())
	namingStrategy := schema.NamingStrategy{SingularTable: true, NoLowerCase: false}
	logrus.Info(namingStrategy)
	db, err := gorm.Open(mysql.Open(cfg.GetDsn()),
		&gorm.Config{
			Logger: log.NewLogger(logger, gormLogger.Config{
				SlowThreshold: time.Duration(cfg.LogConfig.GetSlowLogSecond()) * time.Second,
				LogLevel:      gormLogger.Info,
			}),
			NamingStrategy: namingStrategy,
		})
	if err != nil {
		panic(err)
	}

	_db, err := db.DB()
	if err != nil {
		panic(err)
	}
	if cfg.GetConnConfig().MaxOpenConn != nil {
		logrus.Info(_prefixLog, "maxIdleConn", cfg.GetConnConfig().GetMaxIdleConn())
		_db.SetMaxOpenConns(int(cfg.GetConnConfig().GetMaxOpenConn()))
	} else {
		logrus.Warn(_prefixLog, "maxIdleConn_default", 100)
		_db.SetMaxOpenConns(100)
	}

	if cfg.GetConnConfig().ConnMaxLifeSecond != nil {
		logrus.Info(_prefixLog, "maxIdleSecond", time.Duration(cfg.GetConnConfig().GetConnMaxIdleSecond()*int64(time.Second)).String())
		_db.SetConnMaxLifetime(time.Duration(cfg.GetConnConfig().GetConnMaxIdleSecond()) * time.Second)
	} else {
		logrus.Warn(_prefixLog, "maxIdleSecond", time.Duration(3600*time.Second).String())
		_db.SetConnMaxLifetime(3600 * time.Second)
	}

	if cfg.GetConnConfig().ConnMaxIdleSecond != nil {
		logrus.Info(_prefixLog, "maxIdleConn", cfg.GetConnConfig().GetMaxIdleConn())
		_db.SetMaxIdleConns(int(cfg.GetConnConfig().GetMaxIdleConn()))
	} else {
		logrus.Warn(_prefixLog, "maxIdleConn", 100)
		_db.SetMaxIdleConns(100)
	}

	if cfg.GetConnConfig().ConnMaxIdleSecond != nil {
		logrus.Info(_prefixLog, "connMaxIdleSecond", time.Duration(cfg.GetConnConfig().GetConnMaxIdleSecond()*int64(time.Second)).String())
		_db.SetConnMaxIdleTime(time.Duration(cfg.GetConnConfig().GetConnMaxIdleSecond()) * time.Second)
	} else {
		logrus.Warn(_prefixLog, "connMaxIdleSecond", (3600 * time.Second).String())
		_db.SetConnMaxIdleTime(3600 * time.Second)
	}

	return &DB{db: db}, func() {
		err := _db.Close()
		logrus.Warn("close db ", "mysql close , err info ", err)
	}
}

func (c *DB) WithContext(ctx context.Context) *gorm.DB {
	return c.db.WithContext(ctx)
}
