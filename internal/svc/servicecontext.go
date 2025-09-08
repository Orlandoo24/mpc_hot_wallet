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
	//wsURL := "wss://bsc-rpc.publicnode.com"
	wsURL := "wss://bsc-testnet-rpc.publicnode.com"

	// 从数据库获取所有钱包地址
	watchAddresses := svc.getWalletAddressesFromDB()

	if len(watchAddresses) == 0 {
		log.Println("⚠️  数据库中没有找到钱包地址，跳过BSC监控启动")
		return
	}

	log.Printf("📍 将监控 %d 个钱包地址", len(watchAddresses))
	for i, addr := range watchAddresses {
		if i < 5 { // 只显示前5个地址
			log.Printf("   - %s", addr)
		} else if i == 5 {
			log.Printf("   - ... 还有 %d 个地址", len(watchAddresses)-5)
			break
		}
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

// getWalletAddressesFromDB 从数据库获取钱包地址
func (svc *ServiceContext) getWalletAddressesFromDB() []string {
	// 查询所有钱包地址
	wallets, err := svc.WalletsDao.FindAll(context.Background())
	if err != nil {
		log.Printf("⚠️  获取钱包地址失败: %v", err)
		return []string{}
	}

	var addresses []string
	for _, wallet := range wallets {
		if wallet.Address != "" {
			addresses = append(addresses, wallet.Address)
		}
	}

	return addresses
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
