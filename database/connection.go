package database

import (
	"fmt"
	"github.com/mmtaee/go-oc-utils.git/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"time"
)

// DBConfig database configs
type DBConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

var (
	db      *gorm.DB
	DBDebug bool
)

// Connect connecting to database with configs.
func Connect(cfg *DBConfig, debug bool) {
	var (
		err     error
		gormCfg *gorm.Config
	)
	if debug {
		DBDebug = true
		gormCfg = &gorm.Config{
			Logger: gormLogger.Default.LogMode(gormLogger.Info),
		}
	} else {
		gormCfg = &gorm.Config{
			Logger: gormLogger.Discard,
		}
	}

	db, err = gorm.Open(
		postgres.Open(
			fmt.Sprintf(
				"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
				cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name,
			),
		),
		gormCfg,
	)

	if err != nil {
		logger.Log(logger.CRITICAL, err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		logger.Log(logger.CRITICAL, err)
	}
	for i := 0; i < 5; i++ {
		err = sqlDB.Ping()
		if err == nil {
			break
		}
		logger.Log(logger.WARNING, fmt.Sprintf("Database is not ready: %v\n", err))
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		logger.Log(logger.CRITICAL, fmt.Sprintf("Database connection failed: %v\n", err))
	}
	logger.Log(logger.INFO, "Database connection succeeded!")
}

// Connection get database connection
func Connection() *gorm.DB {
	if DBDebug {
		db.Debug()
	}
	return db
}

// Close database connection closing
func Close() {
	sqlDB, err := db.DB()
	if err != nil {
		logger.Log(logger.CRITICAL, err)
	}
	err = sqlDB.Close()
	if err != nil {
		logger.Log(logger.CRITICAL, err)
	}
}
