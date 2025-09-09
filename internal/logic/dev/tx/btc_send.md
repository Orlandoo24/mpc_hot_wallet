以下是完整的 Cursor 提示词，帮你快速开发 BTC 转账功能：

## Cursor 提示词

```
# 任务：为现有的 Go 转账系统添加 Bitcoin (BTC) 支持

## 背景
我有一个现有的 Go 转账系统，已经支持 EVM 链和 Solana。现在需要添加 Bitcoin 转账支持，使用 LI.FI API 来简化 BTC UTXO 处理。

## 现有代码结构
- `WrapSend` 函数是主入口，已有 EVM 和 Solana 分支
- 使用 `isSolanaChain()` 检测链类型
- 通过 LI.FI API 获取 quote 和交易数据
- 异步发送交易并返回 TxHash

## 需要实现的功能

### 1. 添加 BTC 链检测函数
```go
func (l *TransactionLogic) isBTCChain(chain string) bool {
    btcChains := []string{"BTC", "Bitcoin", "btc", "bitcoin", "20000000000001"}
    for _, c := range btcChains {
        if strings.EqualFold(chain, c) {
            return true
        }
    }
    return false
}
```

### 2. 在 WrapSend 中添加 BTC 分支
在现有的 Solana 检测后添加：
```go
// 检测是否为 Bitcoin 链
if l.isBTCChain(req.Chain) {
    return l.handleBTCTransfer(req)
}
```

### 3. 实现 handleBTCTransfer 函数
需要实现以下步骤：
- 调用 LI.FI quote API 获取 BTC 转账数据
- 解析 PSBT (Partially Signed Bitcoin Transaction)
- 从 PSBT 中提取 memo 信息
- 签名并广播 BTC 交易
- 返回交易响应

### 4. LI.FI API 集成要点

#### BTC Chain ID 和 Token 地址
- Bitcoin Chain ID: `20000000000001`
- 原生 BTC token: `"bitcoin"`
- 只支持 SegWit 地址 (bc1q 开头)

#### Quote API 调用
```go
params := url.Values{}
params.Set("fromChain", "20000000000001")  // Bitcoin Chain ID
params.Set("toChain", "20000000000001")    // 同链转账
params.Set("fromToken", "bitcoin")         // BTC 原生代币
params.Set("toToken", "bitcoin")
params.Set("fromAmount", req.Amount)
params.Set("fromAddress", req.FromAddress)
params.Set("toAddress", req.ToAddress)
params.Set("integrator", "mpc_go-demo")
```

#### PSBT 处理
- LI.FI 返回的 `transactionRequest.data` 是 PSBT hex 格式
- 需要解码 PSBT 提取 memo 信息
- memo 存储在 OP_RETURN 输出中

### 5. 需要的 Go 依赖包
```go
import (
    "github.com/btcsuite/btcd/btcutil"
    "github.com/btcsuite/btcd/chaincfg"
    "github.com/btcsuite/btcd/txscript"
    "github.com/btcsuite/btcd/wire"
    // 或者使用更现代的库
    // "github.com/btcsuite/btcd/btcec/v2"
)
```

### 6. 实现要求
- 保持与现有 EVM/Solana 代码风格一致
- 使用相同的错误处理模式
- 保持异步发送交易的设计
- 添加适当的日志记录
- 构建正确的区块浏览器链接 (mempool.space)

### 7. 关键注意事项
- BTC 地址必须有足够的 UTXO 才能获取 quote
- 交易需要在 30 分钟内发送，建议获取最新 quote
- PSBT 需要用 BTC 钱包签名
- memo 用于跨链转账时指定目标链和地址

## 期望输出
请生成完整的 `handleBTCTransfer` 函数和相关辅助函数，包括：
- BTC quote 获取逻辑
- PSBT 解析和 memo 提取
- BTC 交易签名和广播
- 错误处理和日志记录
- 与现有代码风格保持一致

请确保代码可以直接集成到现有系统中，并提供必要的注释说明。
```

这个提示词包含了：
1. **完整的背景信息** - 让 Cursor 理解现有系统
2. **具体的技术要求** - LI.FI API 集成细节
3. **代码结构指导** - 保持一致的代码风格
4. **关键参数和配置** - BTC Chain ID、token 地址等
5. **实现步骤** - 清晰的开发路径
6. **注意事项** - 避免常见陷阱

使用这个提示词，Cursor 应该能够生成可直接使用的 BTC 转账代码。