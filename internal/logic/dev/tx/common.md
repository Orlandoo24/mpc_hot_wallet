# 通用函数抽取分析

## 需要抽取到 transaction_logic.go 的通用函数

### 1. 代币相关通用函数

#### 1.1 原生代币判断
**位置**: `send_logic.go:264-277`, `bridge_logic.go:470-483`, `swap_logic.go:133-144`
**函数名**: `isNativeToken(token string) bool`
**功能**: 判断代币地址是否为原生代币
**参数**: 
- `token string` - 代币地址
**返回值**: 
- `bool` - 是否为原生代币
**重复次数**: 3次

#### 1.2 检查是否需要 approve
**位置**: `swap_logic.go:132-144`
**函数名**: `needsApproval(fromToken string) bool`
**功能**: 检查 ERC20 代币是否需要 approve
**参数**:
- `fromToken string` - 源代币地址
**返回值**:
- `bool` - 是否需要 approve
**重复次数**: 1次（但逻辑与 isNativeToken 相反）

### 2. ERC20 相关通用函数

#### 2.1 构建 ERC20 transfer 调用数据
**位置**: `send_logic.go:316-334`
**函数名**: `buildERC20TransferData(toAddress string, amount *big.Int) ([]byte, error)`
**功能**: 构建 ERC20 transfer 函数的调用数据
**参数**:
- `toAddress string` - 接收地址
- `amount *big.Int` - 转账金额
**返回值**:
- `[]byte` - 编码后的调用数据
- `error` - 错误信息
**重复次数**: 1次

#### 2.2 构建 ERC20 approve 调用数据
**位置**: `bridge_logic.go:280-297`, `swap_logic.go:160-175`
**函数名**: `buildERC20ApproveData(spenderAddress string, amount *big.Int) []byte`
**功能**: 构建 ERC20 approve 函数的调用数据
**参数**:
- `spenderAddress string` - 授权地址
- `amount *big.Int` - 授权金额
**返回值**:
- `[]byte` - 编码后的调用数据
**重复次数**: 2次

#### 2.3 检查 ERC20 allowance
**位置**: `bridge_logic.go:541-571`
**函数名**: `checkAllowance(client *ethclient.Client, tokenAddress, owner, spender string) (*big.Int, error)`
**功能**: 检查 ERC20 代币的当前 allowance
**参数**:
- `client *ethclient.Client` - RPC 客户端
- `tokenAddress string` - 代币合约地址
- `owner string` - 拥有者地址
- `spender string` - 被授权者地址
**返回值**:
- `*big.Int` - 当前 allowance
- `error` - 错误信息
**重复次数**: 1次

### 3. 交易相关通用函数

#### 3.1 构建并发送 approve 交易
**位置**: `bridge_logic.go:280-356`, `swap_logic.go:147-202`
**函数名**: `executeApproveTransaction(client *ethclient.Client, privateKey *ecdsa.PrivateKey, tokenAddress, spenderAddress string, amount *big.Int, chainId int64) (string, error)`
**功能**: 执行 ERC20 approve 交易
**参数**:
- `client *ethclient.Client` - RPC 客户端
- `privateKey *ecdsa.PrivateKey` - 私钥
- `tokenAddress string` - 代币合约地址
- `spenderAddress string` - 授权地址
- `amount *big.Int` - 授权金额
- `chainId int64` - 链ID
**返回值**:
- `string` - 交易哈希
- `error` - 错误信息
**重复次数**: 2次

#### 3.2 等待交易确认
**位置**: `bridge_logic.go:359-382`
**函数名**: `waitForTransactionReceipt(client *ethclient.Client, txHash common.Hash, timeout time.Duration) (*evmTypes.Receipt, error)`
**功能**: 等待交易确认
**参数**:
- `client *ethclient.Client` - RPC 客户端
- `txHash common.Hash` - 交易哈希
- `timeout time.Duration` - 超时时间
**返回值**:
- `*evmTypes.Receipt` - 交易收据
- `error` - 错误信息
**重复次数**: 1次

#### 3.3 构建并发送交易
**位置**: `bridge_logic.go:385-467`, `swap_logic.go:265-380`
**函数名**: `buildAndSendTransaction(client *ethclient.Client, privateKey *ecdsa.PrivateKey, to common.Address, value *big.Int, data []byte, gasLimit uint64, gasPrice *big.Int, chainId int64) (string, error)`
**功能**: 构建并发送交易
**参数**:
- `client *ethclient.Client` - RPC 客户端
- `privateKey *ecdsa.PrivateKey` - 私钥
- `to common.Address` - 目标地址
- `value *big.Int` - 转账金额
- `data []byte` - 调用数据
- `gasLimit uint64` - Gas 限制
- `gasPrice *big.Int` - Gas 价格
- `chainId int64` - 链ID
**返回值**:
- `string` - 交易哈希
- `error` - 错误信息
**重复次数**: 2次

### 4. Gas 估算通用函数

#### 4.1 估算原生代币转账 Gas
**位置**: `send_logic.go:279-313`
**函数名**: `estimateNativeTransferGas(client *ethclient.Client, fromAddress, toAddress common.Address, value *big.Int) (uint64, *big.Int, error)`
**功能**: 估算原生代币转账的 Gas
**参数**:
- `client *ethclient.Client` - RPC 客户端
- `fromAddress common.Address` - 发送地址
- `toAddress common.Address` - 接收地址
- `value *big.Int` - 转账金额
**返回值**:
- `uint64` - Gas 限制
- `*big.Int` - Gas 价格
- `error` - 错误信息
**重复次数**: 1次

#### 4.2 估算 ERC20 转账 Gas
**位置**: `send_logic.go:337-370`
**函数名**: `estimateERC20TransferGas(client *ethclient.Client, fromAddress, tokenAddress common.Address, data []byte) (uint64, *big.Int, error)`
**功能**: 估算 ERC20 转账的 Gas
**参数**:
- `client *ethclient.Client` - RPC 客户端
- `fromAddress common.Address` - 发送地址
- `tokenAddress common.Address` - 代币合约地址
- `data []byte` - 调用数据
**返回值**:
- `uint64` - Gas 限制
- `*big.Int` - Gas 价格
- `error` - 错误信息
**重复次数**: 1次

### 5. 链相关通用函数

#### 5.1 根据链ID获取链名称
**位置**: `bridge_logic.go:396-414`
**函数名**: `getChainNameByID(chainId int) string`
**功能**: 根据链ID获取链名称
**参数**:
- `chainId int` - 链ID
**返回值**:
- `string` - 链名称
**重复次数**: 1次

#### 5.2 构建区块浏览器链接
**位置**: `send_logic.go:211-239`, `bridge_logic.go:417-435`
**函数名**: `buildExplorerUrl(chain string, txHash string) string`
**功能**: 根据链类型构建区块浏览器链接
**参数**:
- `chain string` - 链名称
- `txHash string` - 交易哈希
**返回值**:
- `string` - 浏览器链接
**重复次数**: 2次

#### 5.3 获取链显示名称
**位置**: `send_logic.go:242-261`
**函数名**: `getChainDisplayName(chain string) string`
**功能**: 获取链的显示名称
**参数**:
- `chain string` - 链名称
**返回值**:
- `string` - 显示名称
**重复次数**: 1次

### 6. 重试机制通用函数

#### 6.1 带重试的交易发送
**位置**: `bridge_logic.go:592-607`
**函数名**: `sendTransactionWithRetry(client *ethclient.Client, privateKey *ecdsa.PrivateKey, to common.Address, value *big.Int, data []byte, gasLimit uint64, gasPrice *big.Int, chainId int64, maxRetries int) (string, error)`
**功能**: 带重试机制的交易发送
**参数**:
- `client *ethclient.Client` - RPC 客户端
- `privateKey *ecdsa.PrivateKey` - 私钥
- `to common.Address` - 目标地址
- `value *big.Int` - 转账金额
- `data []byte` - 调用数据
- `gasLimit uint64` - Gas 限制
- `gasPrice *big.Int` - Gas 价格
- `chainId int64` - 链ID
- `maxRetries int` - 最大重试次数
**返回值**:
- `string` - 交易哈希
- `error` - 错误信息
**重复次数**: 1次

### 7. 钱包相关通用函数

#### 7.1 获取钱包私钥
**位置**: `send_logic.go:42-55`, `bridge_logic.go:152-163`, `swap_logic.go:93-104`
**函数名**: `getWalletPrivateKey(fromAddress string) (*ecdsa.PrivateKey, error)`
**功能**: 从数据库获取钱包私钥
**参数**:
- `fromAddress string` - 钱包地址
**返回值**:
- `*ecdsa.PrivateKey` - 私钥
- `error` - 错误信息
**重复次数**: 3次

## 抽取优先级

### 高优先级（重复次数多或核心功能）
1. `isNativeToken` - 3次重复
2. `buildERC20ApproveData` - 2次重复
3. `executeApproveTransaction` - 2次重复
4. `buildAndSendTransaction` - 2次重复
5. `buildExplorerUrl` - 2次重复
6. `getWalletPrivateKey` - 3次重复

### 中优先级（功能重要但重复次数少）
1. `buildERC20TransferData` - 1次重复
2. `checkAllowance` - 1次重复
3. `waitForTransactionReceipt` - 1次重复
4. `estimateNativeTransferGas` - 1次重复
5. `estimateERC20TransferGas` - 1次重复

### 低优先级（工具函数）
1. `getChainNameByID` - 1次重复
2. `getChainDisplayName` - 1次重复
3. `sendTransactionWithRetry` - 1次重复

## 注意事项

1. 所有函数都需要添加适当的错误处理
2. 需要统一日志记录格式
3. 考虑添加配置参数（如重试次数、超时时间等）
4. 保持函数签名的一致性
5. 添加详细的函数注释和文档
