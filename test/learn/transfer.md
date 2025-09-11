根据 `send_logic.go` 的代码分析，构建一笔 EVM 链交易需要做以下事情：

[1 tool called]

## 🔧 EVM 链交易构建完整流程分析

根据 `send_logic.go` 代码，构建一笔 EVM 链交易需要经过以下步骤：

### 🎯 **核心流程概览** (WrapSend 函数)

```go
// 位置: send_logic.go:39-200
func (l *TransactionLogic) WrapSend(req *types.TransactionReq) (resp *types.TransactionResp, err error)
```

### 📋 **详细步骤分解**

#### **步骤1: 链类型检测与配置获取**
```go
// 位置: send_logic.go:42-59
// 1. 检测特殊链类型
if l.isSolanaChain(req.Chain) {
    return l.handleSolanaTransfer(req)
}
if l.isBTCChain(req.Chain) {
    return l.handleBTCTransfer(req)
}

// 2. 获取 EVM 链配置
chainConfig, ok := l.svcCtx.Config.Chains[req.Chain]
if !ok {
    return nil, fmt.Errorf("unsupported chain: %s", req.Chain)
}
```

**作用**: 确定目标链并获取必要的配置信息（ChainId、RPC URL等）

#### **步骤2: RPC 客户端连接**
```go
// 位置: send_logic.go:61-69
client, err := ethclient.Dial(chainConfig.RpcUrl)
if err != nil {
    return nil, errors.New("failed to connect to chain")
}
defer client.Close()
```

**作用**: 建立与区块链网络的连接

#### **步骤3: 私钥获取与验证**
```go
// 位置: send_logic.go:71-78
privateKey, err := l.GetWalletPrivateKey(req.FromAddress)
if err != nil {
    return nil, err
}
```

**作用**: 从数据库获取发送方的私钥用于签名

#### **步骤4: 收款地址验证**
```go
// 位置: send_logic.go:79-94
toAddr := common.HexToAddress(req.ToAddress)
code, err := client.CodeAt(l.ctx, toAddr, nil)
if err != nil {
    return nil, errors.New("failed to check to address")
}

if len(code) > 0 {
    l.Infof("警告：收款地址是合约地址，可能存在 gas 不足风险")
} else {
    l.Infof("收款地址是 EOA 地址，安全")
}
```

**作用**: 检查收款地址类型，区分 EOA 地址和合约地址

#### **步骤5: 金额解析**
```go
// 位置: send_logic.go:95-98
amount := new(big.Int)
amount.SetString(req.Amount, 10)
```

**作用**: 将字符串金额转换为 big.Int 格式

#### **步骤6: 获取 Nonce**
```go
// 位置: send_logic.go:100-106
nonce, err := client.PendingNonceAt(l.ctx, common.HexToAddress(req.FromAddress))
if err != nil {
    return nil, errors.New("failed to get nonce")
}
```

**作用**: 获取发送方账户的下一个交易序号

#### **步骤7: 交易类型判断与构建**

**7A. 原生代币转账** (BNB/ETH等)
```go
// 位置: send_logic.go:112-135
if l.IsNativeToken(req.FromToken) {
    // Gas 估算
    fromAddr := common.HexToAddress(req.FromAddress)
    gasLimit, gasPrice, err := l.EstimateNativeTransferGas(client, fromAddr, toAddr, amount)
    
    // 构建交易
    tx = evmTypes.NewTx(&evmTypes.LegacyTx{
        Nonce:    nonce,
        To:       &toAddr,
        Value:    amount,
        Gas:      gasLimit,
        GasPrice: gasPrice,
        Data:     nil,  // 原生转账无需 data
    })
}
```

**7B. ERC20 代币转账**
```go
// 位置: send_logic.go:136-167
else {
    // 构建 ERC20 transfer 调用数据
    data, err := l.BuildERC20TransferData(req.ToAddress, amount)
    
    // Gas 估算
    fromAddr := common.HexToAddress(req.FromAddress)
    tokenAddr := common.HexToAddress(req.FromToken)
    gasLimit, gasPrice, err := l.EstimateERC20TransferGas(client, fromAddr, tokenAddr, data)
    
    // 构建交易
    tx = evmTypes.NewTx(&evmTypes.LegacyTx{
        Nonce:    nonce,
        To:       &tokenAddr,  // 目标是代币合约
        Value:    big.NewInt(0),  // ERC20 转账 value 为 0
        Gas:      gasLimit,
        GasPrice: gasPrice,
        Data:     data,  // 包含 transfer 调用数据
    })
}
```

#### **步骤8: 交易签名**
```go
// 位置: send_logic.go:169-176
signedTx, err := evmTypes.SignTx(tx, evmTypes.NewEIP155Signer(big.NewInt(chainConfig.ChainId)), privateKey)
if err != nil {
    return nil, errors.New("failed to sign transaction")
}
```

**作用**: 使用私钥对交易进行数字签名

#### **步骤9: 异步发送交易**
```go
// 位置: send_logic.go:178-199
txHash := signedTx.Hash().Hex()

// 立即返回响应
resp = &types.TransactionResp{
    TxHash:      txHash,
    Message:     l.buildSuccessMessage(req),
    ExplorerUrl: l.BuildExplorerUrl(req.Chain, txHash),
    Chain:       req.Chain,
    Status:      "pending",
}

// 异步发送交易（不阻塞响应）
go func() {
    asyncCtx := context.Background()
    l.sendTransactionAsync(asyncCtx, client, signedTx, txHash)
}()
```

**作用**: 立即返回交易哈希，异步发送到区块链网络

[1 tool called]

### 🔑 **关键技术细节**

#### **Gas 估算机制**

**原生代币转账 Gas 估算:**
```go
// 调用: EstimateNativeTransferGas()
// 特点: 简单转账，Gas 相对固定
gasLimit, gasPrice, err := l.EstimateNativeTransferGas(client, fromAddr, toAddr, amount)
```

**ERC20 代币转账 Gas 估算:**
```go
// 调用: EstimateERC20TransferGas()
// 特点: 合约调用，Gas 消耗更高
gasLimit, gasPrice, err := l.EstimateERC20TransferGas(client, fromAddr, tokenAddr, data)
```

#### **ERC20 Transfer 数据构建**
```go
// 调用: BuildERC20TransferData()
// 功能: 构建 transfer(address,uint256) 的 ABI 编码数据
data, err := l.BuildERC20TransferData(req.ToAddress, amount)
```

#### **异步发送机制**
```go
// 位置: send_logic.go:202-235
func (l *TransactionLogic) sendTransactionAsync(ctx context.Context, client *ethclient.Client, signedTx *evmTypes.Transaction, txHash string) {
    // 重试机制: 最多重试 3 次
    maxRetries := 3
    for i := 0; i < maxRetries; i++ {
        err := client.SendTransaction(ctx, signedTx)
        if err != nil {
            // 指数退避: 1s, 2s, 3s
            time.Sleep(time.Duration(i+1) * time.Second)
            continue
        } else {
            // 发送成功
            return
        }
    }
}
```

### 📊 **交易构建参数对比**

| 交易类型 | To 地址 | Value | Data | Gas 特点 |
|---------|---------|-------|------|----------|
| **原生代币** | 接收方地址 | 转账金额 | `nil` | 21,000 Gas |
| **ERC20代币** | 代币合约地址 | `0` | `transfer()` 调用数据 | ~65,000 Gas |

### 🎯 **核心要点总结**

1. **配置驱动**: 通过链配置获取 ChainId 和 RPC URL
2. **类型区分**: 原生代币 vs ERC20 代币的不同处理逻辑
3. **安全验证**: 收款地址类型检查，避免 OOG 问题
4. **Gas 优化**: 针对不同交易类型的精确 Gas 估算
5. **异步处理**: 立即返回交易哈希，异步发送交易
6. **重试机制**: 3次重试 + 指数退避策略
7. **错误处理**: 完整的错误捕获和日志记录

[1 tool called]

这个流程展示了一个完整的、生产级别的 EVM 链交易构建系统，包含了所有必要的安全检查、错误处理和性能优化措施。