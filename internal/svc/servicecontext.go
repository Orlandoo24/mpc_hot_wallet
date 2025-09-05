package svc

import (
	"log"
	"time"

	"demo/internal/config"
	"demo/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ServiceContext struct {
	Config     config.Config
	WalletsDao model.WalletsDao
	DB         *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 重新从配置文件读取 DSN
	db, err := initDB(c.Postgres.DSN)
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	return &ServiceContext{
		Config:     c,
		WalletsDao: model.NewWalletsDao(db),
		DB:         db,
	}
}

func initDB(dsn string) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}

	// 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return db, nil
}
