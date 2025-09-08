package svc

import (
	"context"
	"log"
	"time"

	"demo/internal/config"
	"demo/internal/logic/monitor"
	"demo/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ServiceContext struct {
	Config        config.Config
	WalletsDao    model.WalletsDao
	DB            *gorm.DB
	MonitorCancel context.CancelFunc // 用于停止监控
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 重新从配置文件读取 DSN
	db, err := initDB(c.Postgres.DSN)
	if err != nil {
		log.Fatalf("failed to init db: %v", err)
	}

	svcCtx := &ServiceContext{
		Config:     c,
		WalletsDao: model.NewWalletsDao(db),
		DB:         db,
	}

	// 启动BSC监控
	svcCtx.startBSCMonitor()

	return svcCtx
}

// startBSCMonitor 启动BSC监控
func (svc *ServiceContext) startBSCMonitor() {
	// BSC WebSocket URL
	wsURL := "wss://bsc-rpc.publicnode.com"

	// 监控地址列表 (可以从配置文件读取，这里先硬编码为空，监控所有交易)
	watchAddresses := []string{
		// 可以添加需要监控的特定地址
		// "0x742d35Cc6474C4532C59AE3f970d89BbE8Fb0eEA",
	}

	// 创建监控上下文
	ctx, cancel := context.WithCancel(context.Background())
	svc.MonitorCancel = cancel

	// 在后台启动监控
	go func() {
		log.Println("🚀 启动BSC链监控服务...")
		if err := monitor.StartBSCMonitoring(ctx, wsURL, watchAddresses); err != nil {
			if err != context.Canceled {
				log.Printf("❌ BSC监控服务异常: %v", err)
			} else {
				log.Println("✅ BSC监控服务已停止")
			}
		}
	}()
}

// StopMonitor 停止监控服务
func (svc *ServiceContext) StopMonitor() {
	if svc.MonitorCancel != nil {
		log.Println("🛑 正在停止BSC监控服务...")
		svc.MonitorCancel()
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
