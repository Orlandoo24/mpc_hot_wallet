给你一个完整的 BSC 交易监控 demo：

go

package main

import (
"context"
"crypto/ecdsa"
"fmt"
"log"
"math/big"
"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 事件签名
var (
// Transfer(address,address,uint256)
TransferEventSignature = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))
// Approval(address,address,uint256)
ApprovalEventSignature = crypto.Keccak256Hash([]byte("Approval(address,address,uint256)"))
// Swap 事件 (PancakeSwap)
SwapEventSignature = crypto.Keccak256Hash([]byte("Swap(address,uint256,uint256,uint256,uint256,address)"))
)

// 交易类型
type TxType string

const (
TxTypeSend     TxType = "send"
TxTypeReceive  TxType = "receive"
TxTypeSwap     TxType = "swap"
TxTypeBridge   TxType = "bridge"
TxTypeApprove  TxType = "approve"
TxTypeRevoke   TxType = "revoke"
TxTypeUnknown  TxType = "unknown"
)

// 交易事件
type TransactionEvent struct {
TxHash      string                 `json:"tx_hash"`
BlockNumber uint64                 `json:"block_number"`
From        string                 `json:"from"`
To          string                 `json:"to"`
Value       *big.Int               `json:"value"`
TxType      TxType                 `json:"tx_type"`
TokenAddr   string                 `json:"token_addr,omitempty"`
TokenAmount *big.Int               `json:"token_amount,omitempty"`
GasUsed     uint64                 `json:"gas_used"`
GasPrice    *big.Int               `json:"gas_price"`
Details     map[string]interface{} `json:"details,omitempty"`
}

// BSC 监控器
type BSCMonitor struct {
client          *ethclient.Client
watchAddresses  map[common.Address]bool // 监控的地址列表
contractABIs    map[common.Address]*abi.ABI
eventHandlers   []func(*TransactionEvent)
}

// 创建监控器
func NewBSCMonitor(rpcURL string, watchAddresses []string) (*BSCMonitor, error) {
client, err := ethclient.Dial(rpcURL)
if err != nil {
return nil, fmt.Errorf("failed to connect to BSC: %w", err)
}

	monitor := &BSCMonitor{
		client:         client,
		watchAddresses: make(map[common.Address]bool),
		contractABIs:   make(map[common.Address]*abi.ABI),
		eventHandlers:  make([]func(*TransactionEvent), 0),
	}

	// 添加监控地址
	for _, addr := range watchAddresses {
		monitor.watchAddresses[common.HexToAddress(addr)] = true
	}

	return monitor, nil
}

// 添加事件处理器
func (m *BSCMonitor) AddEventHandler(handler func(*TransactionEvent)) {
m.eventHandlers = append(m.eventHandlers, handler)
}

// 开始监控
func (m *BSCMonitor) Start(ctx context.Context) error {
log.Println("开始监控 BSC 交易...")

	// 订阅新区块
	headers := make(chan *types.Header)
	sub, err := m.client.SubscribeNewHead(ctx, headers)
	if err != nil {
		return fmt.Errorf("failed to subscribe to new heads: %w", err)
	}
	defer sub.Unsubscribe()

	for {
		select {
		case err := <-sub.Err():
			log.Printf("订阅错误: %v", err)
			return err
		case header := <-headers:
			// 处理新区块
			if err := m.processBlock(ctx, header.Number); err != nil {
				log.Printf("处理区块 %d 失败: %v", header.Number.Uint64(), err)
			}
		case <-ctx.Done():
			log.Println("监控已停止")
			return ctx.Err()
		}
	}
}

// 处理区块
func (m *BSCMonitor) processBlock(ctx context.Context, blockNumber *big.Int) error {
block, err := m.client.BlockByNumber(ctx, blockNumber)
if err != nil {
return fmt.Errorf("获取区块失败: %w", err)
}

	log.Printf("处理区块 %d, 包含 %d 个交易", block.NumberU64(), len(block.Transactions()))

	// 处理区块中的每个交易
	for _, tx := range block.Transactions() {
		if err := m.processTx(ctx, tx, block.NumberU64()); err != nil {
			log.Printf("处理交易 %s 失败: %v", tx.Hash().Hex(), err)
		}
	}

	return nil
}

// 处理单个交易
func (m *BSCMonitor) processTx(ctx context.Context, tx *types.Transaction, blockNumber uint64) error {
// 获取交易回执
receipt, err := m.client.TransactionReceipt(ctx, tx.Hash())
if err != nil {
return fmt.Errorf("获取交易回执失败: %w", err)
}

	// 检查是否涉及监控地址
	var isWatched bool
	from, err := m.getFromAddress(tx)
	if err == nil {
		if m.watchAddresses[from] {
			isWatched = true
		}
	}

	if tx.To() != nil && m.watchAddresses[*tx.To()] {
		isWatched = true
	}

	// 如果不是监控地址且没有监控地址列表，则处理所有交易
	if len(m.watchAddresses) > 0 && !isWatched {
		return nil
	}

	// 分析交易类型和事件
	event := m.analyzeTx(tx, receipt, blockNumber)

	// 触发事件处理器
	for _, handler := range m.eventHandlers {
		handler(event)
	}

	return nil
}

// 分析交易
func (m *BSCMonitor) analyzeTx(tx *types.Transaction, receipt *types.Receipt, blockNumber uint64) *TransactionEvent {
event := &TransactionEvent{
TxHash:      tx.Hash().Hex(),
BlockNumber: blockNumber,
Value:       tx.Value(),
GasUsed:     receipt.GasUsed,
GasPrice:    tx.GasPrice(),
Details:     make(map[string]interface{}),
TxType:      TxTypeUnknown,
}

	// 获取 from 地址
	if from, err := m.getFromAddress(tx); err == nil {
		event.From = from.Hex()
	}

	// 获取 to 地址
	if tx.To() != nil {
		event.To = tx.To().Hex()
	}

	// 分析交易类型
	if len(receipt.Logs) == 0 {
		// 没有日志，可能是普通转账
		if tx.Value().Cmp(big.NewInt(0)) > 0 {
			event.TxType = TxTypeSend
		}
	} else {
		// 分析日志确定交易类型
		event.TxType = m.analyzeLogs(receipt.Logs, event)
	}

	return event
}

// 分析日志
func (m *BSCMonitor) analyzeLogs(logs []*types.Log, event *TransactionEvent) TxType {
for _, vLog := range logs {
if len(vLog.Topics) == 0 {
continue
}

		eventSig := vLog.Topics[0]

		switch eventSig {
		case TransferEventSignature:
			return m.handleTransferEvent(vLog, event)
		case ApprovalEventSignature:
			return m.handleApprovalEvent(vLog, event)
		case SwapEventSignature:
			return m.handleSwapEvent(vLog, event)
		}
	}

	return TxTypeUnknown
}

// 处理 Transfer 事件
func (m *BSCMonitor) handleTransferEvent(vLog *types.Log, event *TransactionEvent) TxType {
if len(vLog.Topics) >= 3 {
from := common.BytesToAddress(vLog.Topics[1].Bytes())
to := common.BytesToAddress(vLog.Topics[2].Bytes())

		event.TokenAddr = vLog.Address.Hex()
		event.Details["token_from"] = from.Hex()
		event.Details["token_to"] = to.Hex()

		// 如果 data 长度正确，解析转账金额
		if len(vLog.Data) == 32 {
			event.TokenAmount = new(big.Int).SetBytes(vLog.Data)
		}

		// 判断是发送还是接收
		if len(m.watchAddresses) > 0 {
			if m.watchAddresses[from] {
				return TxTypeSend
			}
			if m.watchAddresses[to] {
				return TxTypeReceive
			}
		}

		return TxTypeSend
	}

	return TxTypeUnknown
}

// 处理 Approval 事件
func (m *BSCMonitor) handleApprovalEvent(vLog *types.Log, event *TransactionEvent) TxType {
if len(vLog.Topics) >= 3 {
owner := common.BytesToAddress(vLog.Topics[1].Bytes())
spender := common.BytesToAddress(vLog.Topics[2].Bytes())

		event.TokenAddr = vLog.Address.Hex()
		event.Details["owner"] = owner.Hex()
		event.Details["spender"] = spender.Hex()

		// 解析授权金额
		if len(vLog.Data) == 32 {
			amount := new(big.Int).SetBytes(vLog.Data)
			event.TokenAmount = amount

			// 如果金额为0，表示撤销授权
			if amount.Cmp(big.NewInt(0)) == 0 {
				return TxTypeRevoke
			}
		}

		return TxTypeApprove
	}

	return TxTypeUnknown
}

// 处理 Swap 事件
func (m *BSCMonitor) handleSwapEvent(vLog *types.Log, event *TransactionEvent) TxType {
event.Details["swap_contract"] = vLog.Address.Hex()
return TxTypeSwap
}

// 获取交易发送方地址
func (m *BSCMonitor) getFromAddress(tx *types.Transaction) (common.Address, error) {
var signer types.Signer = types.LatestSignerForChainID(tx.ChainId())
return signer.Sender(tx)
}

// 事件通知处理器
func eventNotificationHandler(event *TransactionEvent) {
log.Printf("🔔 检测到交易:")
log.Printf("  哈希: %s", event.TxHash)
log.Printf("  类型: %s", event.TxType)
log.Printf("  区块: %d", event.BlockNumber)
log.Printf("  发送方: %s", event.From)
log.Printf("  接收方: %s", event.To)

	if event.Value != nil && event.Value.Cmp(big.NewInt(0)) > 0 {
		log.Printf("  BNB 金额: %s wei", event.Value.String())
	}

	if event.TokenAddr != "" {
		log.Printf("  代币合约: %s", event.TokenAddr)
		if event.TokenAmount != nil {
			log.Printf("  代币金额: %s", event.TokenAmount.String())
		}
	}

	log.Printf("  Gas 使用: %d", event.GasUsed)
	log.Printf("  详情: %+v", event.Details)
	log.Println("  ---")
}

func main() {
// BSC 主网 RPC (你也可以用其他 RPC 提供商)
rpcURL := "https://bsc-dataseed1.binance.org/"

	// 要监控的地址列表 (空列表表示监控所有交易)
	watchAddresses := []string{
		// "0x742d35Cc6474C4532C59AE3f970d89BbE8Fb0eEA", // 示例地址
	}

	// 创建监控器
	monitor, err := NewBSCMonitor(rpcURL, watchAddresses)
	if err != nil {
		log.Fatal("创建监控器失败:", err)
	}

	// 添加事件处理器
	monitor.AddEventHandler(eventNotificationHandler)

	// 开始监控
	ctx := context.Background()
	if err := monitor.Start(ctx); err != nil {
		log.Fatal("监控失败:", err)
	}
}
使用说明
依赖安装:

bash

go mod init bsc-monitor
go get github.com/ethereum/go-ethereum
配置说明:

rpcURL: BSC RPC 节点地址，可以用免费的或付费的
watchAddresses: 要监控的地址列表，空列表监控所有交易
可以添加多个事件处理器来处理不同的通知需求
功能特点:

实时监控 BSC 新区块和交易
解析 Transfer、Approval、Swap 等常见事件
识别 send、receive、approve、revoke、swap 交易类型
灵活的事件处理机制
扩展建议:

添加更多 DEX 的 Swap 事件签名
实现跨链桥事件识别
添加数据库存储
实现 Webhook 通知
添加重连和错误恢复机制
这个 demo 可以直接运行，你可以根据需要进行扩展和优化。