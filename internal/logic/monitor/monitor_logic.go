package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// TokenEvent 代币事件结构
type TokenEvent struct {
	BlockNumber uint64 `json:"blockNumber"`
	TxHash      string `json:"txHash"`
	Timestamp   int64  `json:"timestamp"`
	EventType   string `json:"eventType"` // ABI标准事件名称: Transfer/Approval/Swap/Bridge/Deposit/Withdrawal/NativeTransfer
	Direction   string `json:"direction"` // IN/OUT/NONE - 资金流向标记
	FromAddr    string `json:"fromAddr"`
	ToAddr      string `json:"toAddr"`
	TokenAddr   string `json:"tokenAddr"` // 代币合约地址
	Amount      string `json:"amount"`    // 使用string存储以避免精度问题
	ChainId     uint64 `json:"chainId"`   // 支持跨链场景
}

// BSCMonitor BSC监控器
type BSCMonitor struct {
	client         *ethclient.Client
	watchAddresses map[common.Address]bool // 监控的地址列表
	eventHandlers  []func(*TokenEvent)
	chainId        uint64
	logParser      *LogParser
}

// NewBSCMonitor 创建BSC监控器
func NewBSCMonitor(wsURL string, watchAddresses []string) (*BSCMonitor, error) {
	client, err := ethclient.Dial(wsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to BSC WebSocket: %w", err)
	}

	// 获取链ID
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %w", err)
	}

	monitor := &BSCMonitor{
		client:         client,
		watchAddresses: make(map[common.Address]bool),
		eventHandlers:  make([]func(*TokenEvent), 0),
		chainId:        chainID.Uint64(),
		logParser:      NewLogParser(),
	}

	// 添加监控地址
	for _, addr := range watchAddresses {
		monitor.watchAddresses[common.HexToAddress(addr)] = true
	}

	return monitor, nil
}

// AddEventHandler 添加事件处理器
func (m *BSCMonitor) AddEventHandler(handler func(*TokenEvent)) {
	m.eventHandlers = append(m.eventHandlers, handler)
}

// Start 开始监控
func (m *BSCMonitor) Start(ctx context.Context) error {
	log.Println("开始监控 BSC 交易...")

	// 订阅新区块头
	headers := make(chan *types.Header)
	sub, err := m.client.SubscribeNewHead(ctx, headers)
	if err != nil {
		return fmt.Errorf("failed to subscribe to new heads: %w", err)
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			log.Printf("WebSocket订阅错误: %v", err)
			return err
		case header := <-headers:
			// 处理新区块
			if err := m.processBlock(ctx, header); err != nil {
				log.Printf("处理区块 %d 失败: %v", header.Number.Uint64(), err)
			}
		case <-ctx.Done():
			log.Println("BSC监控已停止")
			return ctx.Err()
		}
	}
}

// processBlock 处理区块
func (m *BSCMonitor) processBlock(ctx context.Context, header *types.Header) error {
	blockNumber := header.Number.Uint64()
	timestamp := int64(header.Time)

	// 获取完整区块信息
	block, err := m.client.BlockByNumber(ctx, header.Number)
	if err != nil {
		return fmt.Errorf("获取区块失败: %w", err)
	}

	// 处理区块中的每个交易
	var relevantTxCount int
	for _, tx := range block.Transactions() {
		// 检查是否为监控的交易
		if m.isWatchedTransaction(tx) {
			// 第一个相关交易时输出区块日志
			if relevantTxCount == 0 {
				log.Printf("🔍 发现相关区块 %d (包含监控地址的交易)", blockNumber)
			}
			relevantTxCount++

			// 处理相关交易
			if err := m.processTx(ctx, tx, blockNumber, timestamp); err != nil {
				log.Printf("处理交易 %s 失败: %v", tx.Hash().Hex(), err)
			}
		}
	}

	// 输出相关交易统计
	if relevantTxCount > 0 {
		log.Printf("✅ 区块 %d 处理完成，找到 %d 个相关交易", blockNumber, relevantTxCount)
	}

	return nil
}

// processTx 处理单个交易（调用前已确认是监控的交易）
func (m *BSCMonitor) processTx(ctx context.Context, tx *types.Transaction, blockNumber uint64, timestamp int64) error {
	// 获取交易回执
	receipt, err := m.client.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		return fmt.Errorf("获取交易回执失败: %w", err)
	}

	// 检查交易状态
	if receipt.Status != types.ReceiptStatusSuccessful {
		log.Printf("⚠️  跳过失败交易: %s", tx.Hash().Hex())
		return nil
	}

	// 解析交易日志，生成TokenEvent
	events := m.logParser.ParseLogs(receipt.Logs, tx, blockNumber, timestamp, m.chainId)

	// 触发事件处理器
	for _, event := range events {
		for _, handler := range m.eventHandlers {
			handler(event)
		}
	}

	return nil
}

// isWatchedTransaction 检查是否为监控的交易
func (m *BSCMonitor) isWatchedTransaction(tx *types.Transaction) bool {
	// 如果没有设置监控地址，不监控任何交易
	if len(m.watchAddresses) == 0 {
		return false
	}

	// 检查发送方地址
	// 修复ChainID为0的问题，BSC链ID为56
	txChainID := tx.ChainId()
	if txChainID == nil || txChainID.Uint64() == 0 {
		txChainID = big.NewInt(56) // BSC Chain ID
	}
	signer := types.LatestSignerForChainID(txChainID)
	if from, err := signer.Sender(tx); err == nil {
		if m.watchAddresses[from] {
			return true
		}
	}

	// 检查接收方地址
	if tx.To() != nil && m.watchAddresses[*tx.To()] {
		return true
	}

	return false
}

// MockKafkaProducer 模拟Kafka生产者
func MockKafkaProducer(event *TokenEvent) {
	eventJSON, _ := json.MarshalIndent(event, "", "  ")
	log.Printf("📤 发送到Kafka: %s", string(eventJSON))

	// TODO: 实际的Kafka生产者实现
	// producer.Send("token-events", eventJSON)
}

// StartBSCMonitoring 启动BSC监控 (对外接口)
func StartBSCMonitoring(ctx context.Context, wsURL string, watchAddresses []string) error {
	// 带重连机制的监控启动
	return StartBSCMonitoringWithReconnect(ctx, wsURL, watchAddresses)
}

// StartBSCMonitoringWithReconnect 带自动重连的BSC监控
func StartBSCMonitoringWithReconnect(ctx context.Context, wsURL string, watchAddresses []string) error {
	for {
		select {
		case <-ctx.Done():
			log.Println("✅ BSC监控服务已停止")
			return ctx.Err()
		default:
			log.Println("🔄 尝试连接BSC监控...")

			monitor, err := NewBSCMonitor(wsURL, watchAddresses)
			if err != nil {
				log.Printf("❌ 创建BSC监控失败: %v, 5秒后重试...", err)
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(5 * time.Second):
					continue
				}
			}

			// 添加Kafka事件处理器
			monitor.AddEventHandler(MockKafkaProducer)

			// 添加日志事件处理器
			monitor.AddEventHandler(func(event *TokenEvent) {
				formattedAmount := FormatTokenAmount(event.Amount, event.TokenAddr, event.ChainId)

				// 为不同方向的事件添加不同的emoji
				var emoji string
				switch event.Direction {
				case "IN":
					emoji = "📥" // 接收
				case "OUT":
					emoji = "📤" // 发送
				default:
					emoji = "🔔" // 其他事件
				}

				// 构建方向标记
				var directionTag string
				if event.Direction != "NONE" && event.Direction != "" {
					directionTag = fmt.Sprintf("-%s", event.Direction)
				}

				log.Printf("%s EVM事件: %s%s | 金额: %s | 哈希: %s",
					emoji, event.EventType, directionTag, formattedAmount, event.TxHash[:10]+"...")
			})

			// 启动监控
			err = monitor.Start(ctx)
			if err != nil {
				if err == context.Canceled {
					log.Println("✅ BSC监控服务已停止")
					return err
				}
				log.Printf("❌ BSC监控连接异常: %v, 3秒后重连...", err)
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(3 * time.Second):
					continue
				}
			}
		}
	}
}
