# MPC Demo - 多链钱包 API 服务

一个基于 Go-Zero 框架构建的企业级多链钱包 API 服务，支持 **EVM 兼容链**、**Bitcoin** 和 **Solana** 三大生态系统。集成 LI.FI 跨链协议和实时区块链监控，提供完整的钱包管理、多链转账、代币交换、跨链桥接、授权管理和交易监控解决方案。

## 🌟 核心特性

### 🔗 多链生态支持
- **EVM 生态** - 以太坊、BSC、Polygon 等 EVM 兼容链
- **Bitcoin 生态** - Bitcoin 测试网原生转账（优化 UTXO 管理）
- **Solana 生态** - Solana 主网和测试网支持，已实现 devnet transfer 和 swap
- **统一接口** - 一套 API 操作所有支持的区块链
- **线下签名** - 所有链交易均采用本地线下签名，确保私钥安全

### 🔐 企业级钱包管理
- **多链钱包初始化** - 同时创建 EVM、Bitcoin、Solana 钱包
- **安全私钥管理** - 数据库加密存储和检索，支持不同链的私钥格式
- **地址验证** - 智能识别和验证不同链的地址格式
- **批量钱包操作** - 支持批量创建和管理钱包
- **离线签名保护** - 私钥永不离开服务器，所有签名操作本地完成

### 💸 全链转账功能
- **EVM 原生转账** - ETH、BNB、MATIC 等原生代币转账
- **ERC20 代币转账** - 支持所有标准 ERC20 代币
- **Bitcoin 转账** - 优化的 UTXO 选择和交易构建
- **Solana 转账** - SPL 代币和 SOL 原生转账，支持 devnet 测试
- **智能 Gas 优化** - 动态费用估算和优化
- **离线签名架构** - 所有转账交易均采用本地私钥签名，确保安全性

### 🔄 代币交换 (Swap)
- **LI.FI 集成** - 使用 LI.FI 协议获得最优交换路由
- **原生 Solana Swap** - 支持 Solana devnet 原生代币交换
- **智能 Approve** - 自动检查和管理代币授权
- **滑点保护** - 可配置的滑点保护机制
- **多 DEX 聚合** - 自动选择最优的去中心化交易所
- **离线签名安全** - 所有交换交易均在本地完成签名

### 🌉 跨链桥接 (Bridge)
- **多链资产桥接** - 支持 EVM 链之间的资产跨链转移
- **LI.FI 协议集成** - 使用业界领先的跨链聚合协议
- **最优路由选择** - 自动选择最快或最便宜的跨链路由
- **实时状态追踪** - 完整的跨链交易进度监控
- **智能重试机制** - 自动处理网络异常和交易失败
- **离线签名安全** - 跨链交易在本地完成签名，私钥不离开服务器

### 🔑 智能授权管理
- **精确授权检查** - 实时查询 ERC20 代币授权状态
- **智能授权策略** - 支持无限授权和精确金额授权
- **批量授权查询** - 一次性查询多个代币的授权状态
- **安全授权撤销** - 一键撤销不必要的代币授权
- **授权历史追踪** - 完整的用户授权操作记录

### 📊 实时链上监控
- **多链监控支持** - 实时监控 EVM 链上交易活动
- **智能事件解析** - 自动解析 Transfer、Approval、Swap 等链上事件
- **交易方向识别** - 精确标记资金流向（IN/OUT/INTERNAL）
- **LI.FI 事件增强** - 集成 LI.FI API 进行高级交易分析
- **WebSocket 实时连接** - 毫秒级的链上数据推送
- **Kafka 数据流** - 结构化事件数据推送到消息队列
- **智能钱包过滤** - 只监控数据库中管理的钱包地址

### ⚡ Bitcoin 优化特性
- **优化 UTXO 管理** - 使用 Blockstream API 替代不稳定的 RPC 连接
- **智能费用计算** - 动态矿工费估算和优化
- **多地址格式支持** - 支持 P2PKH、P2SH、SegWit 等地址格式
- **测试网完整支持** - 完整的 Bitcoin 测试网转账功能

## 🚀 快速开始

### 环境要求
- Go 1.19+
- PostgreSQL 12+
- 稳定的网络连接（用于 RPC 调用）

### 安装步骤

1. **克隆项目**
```bash
git clone <repository-url>
cd mpc_go
```

2. **安装依赖**
```bash
go mod tidy
```

3. **配置数据库**
```bash
# 创建 PostgreSQL 数据库
createdb testdb

# 配置数据库连接（编辑 etc/demo.yaml）
```

4. **启动服务**
```bash
go run main.go
```

服务将在 `http://localhost:8888` 启动，同时自动启动 BSC 区块链监控服务

## 📋 API 接口文档

### 🔐 钱包管理

#### 初始化多链钱包
```http
POST /api/wallet_init
Content-Type: application/json

{
  "address": "0x742d35Cc6634C0532925a3b8D0c5B2b8b8c8c8c8",
  "private_key": "0x1234567890abcdef..."
}
```

### 💸 多链转账操作

#### EVM 链转账（ETH/BNB/MATIC 等）
```http
POST /api/transaction/send
Content-Type: application/json

{
  "chain": "BSC",
  "from_address": "0x742d35Cc6634C0532925a3b8D0c5B2b8b8c8c8c8",
  "to_address": "0x8ba1f109551bD432803012645Hac136c22C57592",
  "from_token": "0x0000000000000000000000000000000000000000",
  "to_token": "0x0000000000000000000000000000000000000000",
  "amount": "1000000000000000000"
}
```

#### Bitcoin 转账
```http
POST /api/transaction/send
Content-Type: application/json

{
  "chain": "BTC",
  "from_address": "n21s8JW8s9YX8tQ5LkbedjVAkBJPA11XcR",
  "to_address": "mtP6cccjQPZw82sbPdSHha6N68o945y4cU",
  "from_token": "BTC",
  "to_token": "BTC",
  "amount": "10000"
}
```
*注：Bitcoin 金额单位为 satoshi (1 BTC = 100,000,000 satoshi)*

#### Solana 转账
```http
POST /api/transaction/send
Content-Type: application/json

{
  "chain": "Solana",
  "from_address": "DYw8jCTfwHNRJhhmFcbXvVDTqWMEVFBX6ZKUmG5CNSKK",
  "to_address": "Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB",
  "from_token": "11111111111111111111111111111111",
  "to_token": "11111111111111111111111111111111",
  "amount": "1000000000"
}
```
*注：Solana 金额单位为 lamports (1 SOL = 1,000,000,000 lamports)*

#### ERC20 代币转账
```http
POST /api/transaction/send
Content-Type: application/json

{
  "chain": "BSC",
  "from_address": "0x742d35Cc6634C0532925a3b8D0c5B2b8b8c8c8c8",
  "to_address": "0x8ba1f109551bD432803012645Hac136c22C57592",
  "from_token": "0x55d398326f99059fF775485246999027B3197955",
  "to_token": "0x55d398326f99059fF775485246999027B3197955",
  "amount": "1000000000000000000"
}
```

### 🔄 代币交换

#### EVM 链代币交换
```http
POST /api/transaction/swap
Content-Type: application/json

{
  "chain": "BSC",
  "from_address": "0x742d35Cc6634C0532925a3b8D0c5B2b8b8c8c8c8",
  "to_address": "0x742d35Cc6634C0532925a3b8D0c5B2b8b8c8c8c8",
  "from_token": "0x0000000000000000000000000000000000000000",
  "to_token": "0x55d398326f99059fF775485246999027B3197955",
  "amount": "100000000000000000"
}
```

### 跨链转账

#### 获取跨链报价
```http
POST /api/bridge/quote
Content-Type: application/json

{
  "from_chain": 56,
  "to_chain": 137,
  "from_token": "0x55d398326f99059fF775485246999027B3197955",
  "to_token": "0xc2132D05D31c914a87C6611C10748AEb04B58e8F",
  "from_amount": "1000000000000000",
  "from_address": "0x...",
  "to_address": "0x...",
  "order": "FASTEST"
}
```

#### 执行跨链转账
```http
POST /api/bridge/wrap
Content-Type: application/json

{
  "from_address": "0x...",
  "from_chain": 56,
  "to_chain": 137,
  "from_token": "0x55d398326f99059fF775485246999027B3197955",
  "to_token": "0xc2132D05D31c914a87C6611C10748AEb04B58e8F",
  "amount": "1000000000000000",
  "to_address": "0x...",
  "order": "CHEAPEST"
}
```

#### 查询跨链状态
```http
POST /api/bridge/status
Content-Type: application/json

{
  "tx_hash": "0x..."
}
```

### 授权管理

#### 检查授权额度
```http
POST /api/transaction/check_allowance
Content-Type: application/json

{
  "token_address": "0x55d398326f99059fF775485246999027B3197955",
  "owner_address": "0x...",
  "spender_address": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
  "chain": "BSC"
}
```

#### 授权代币
```http
POST /api/transaction/approve
Content-Type: application/json

{
  "token_address": "0x55d398326f99059fF775485246999027B3197955",
  "spender_address": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
  "owner_address": "0x...",
  "amount": "max",
  "chain": "BSC"
}
```

#### 撤销授权
```http
POST /api/transaction/revoke
Content-Type: application/json

{
  "token_address": "0x55d398326f99059fF775485246999027B3197955",
  "spender_address": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
  "owner_address": "0x...",
  "chain": "BSC"
}
```

## ⚙️ 配置文件

### etc/demo.yaml
```yaml
Name: demo-api
Host: 0.0.0.0
Port: 8888

Postgres:
  DSN: "postgres://admin:123456@localhost:5432/testdb?sslmode=disable"

Lifi:
  ApiUrl: "https://li.quest/v1"

Chains:
  BSC:
    Name: "Binance Smart Chain Mainnet"
    RpcUrl: "https://bsc-rpc.publicnode.com"
    ChainId: 56
  ETH:
    Name: "Ethereum Mainnet"
    RpcUrl: "https://mainnet.infura.io/v3/your-infura-id"
    ChainId: 1
  Polygon:
    Name: "Polygon Mainnet"
    RpcUrl: "https://polygon-rpc.com"
    ChainId: 137
```

## 🏗️ 项目结构

```
mpc_demo/
├── cli/                    # 命令行工具
├── etc/                    # 配置文件
├── internal/
│   ├── config/            # 配置管理
│   ├── constant/          # 常量定义
│   ├── handler/           # HTTP 处理器
│   ├── logic/             # 业务逻辑
│   │   ├── monitor/       # 区块链监控模块
│   │   │   ├── monitor_logic.go     # 监控核心逻辑
│   │   │   └── log_parse_logic.go   # 事件解析逻辑
│   │   └── transaction/   # 交易处理模块
│   │       ├── approve_logic.go     # 授权管理
│   │       ├── bridge_logic.go      # 跨链转账
│   │       ├── send_logic.go        # 普通转账
│   │       ├── swap_logic.go        # 代币交换
│   │       ├── transaction_logic.go # 通用交易
│   │       └── wallet_logic.go      # 钱包管理
│   ├── mid/               # 中间件
│   ├── model/             # 数据模型
│   ├── svc/               # 服务上下文
│   └── types/             # 类型定义
├── test/                  # 测试文件
├── main.go               # 入口文件
└── README.md             # 项目文档
```

## 🔧 核心技术

- **框架**: Go-Zero (微服务框架)
- **数据库**: PostgreSQL (钱包数据存储)
- **区块链**: go-ethereum (以太坊客户端)
- **跨链**: LI.FI Protocol (跨链聚合协议)
- **监控**: WebSocket + 事件解析 (实时区块链监控)
- **消息队列**: Kafka (事件数据流)
- **API**: RESTful API (JSON 格式)

## 🌟 特色功能

### 🔒 离线签名安全架构
- **本地签名**: 所有交易签名均在应用服务器本地完成
- **私钥隔离**: 私钥永不离开服务器，不依赖外部签名服务
- **多链支持**: EVM (ECDSA)、Bitcoin (ECDSA)、Solana (Ed25519) 全链签名
- **零信任**: 无需信任第三方钱包或签名服务
- **完全控制**: 应用完全控制签名过程和私钥安全

### LI.FI 集成优化
- **最佳实践参数**: 集成官方推荐的优化配置
- **智能路由**: 自动选择最快或最便宜的跨链路由
- **时间策略**: 优化的重试和等待策略
- **费用优化**: 智能 Gas 估算和滑点保护

### 智能错误处理
- **RPC 容错**: 自动处理 RPC 节点的误导性错误
- **重试机制**: 网络异常时的自动重试
- **状态追踪**: 完整的交易状态监控和轮询

### 区块链监控特性
- **实时监控**: WebSocket 实时连接 BSC 测试网
- **智能解析**: 自动解析 EVM 事件并识别交易类型
- **方向标记**: 精确识别 Transfer 事件的 IN/OUT 方向
- **LI.FI 增强**: 集成 LI.FI API 进行高级交易分析
- **自动重连**: 网络异常时自动重连机制
- **数据流**: Kafka 集成，支持事件数据流处理

### 安全特性
- **私钥保护**: 安全的私钥存储和管理
- **授权管理**: 精确的代币授权控制
- **交易验证**: 完整的交易参数验证和确认
- **地址过滤**: 只监控授权钱包地址的相关交易

## 📊 支持的区块链

| 链名称 | Chain ID | 主网/测试网 | RPC 状态 | 监控状态 |
|--------|----------|-------------|----------|----------|
| BSC | 56 | 主网 | ✅ 已配置 | ❌ 未启用 |
| BSC 测试网 | 97 | 测试网 | ✅ 已配置 | ✅ 实时监控 |
| Ethereum | 1 | 主网 | ⚠️ 需配置 | ❌ 未启用 |
| Polygon | 137 | 主网 | ⚠️ 需配置 | ❌ 未启用 |
| Arbitrum | 42161 | 主网 | ⚠️ 需配置 | ❌ 未启用 |
| Optimism | 10 | 主网 | ⚠️ 需配置 | ❌ 未启用 |
| Base | 8453 | 主网 | ⚠️ 需配置 | ❌ 未启用 |

## 🛡️ 安全提醒

⚠️ **重要安全提醒**：
- 这是一个演示项目，**不建议在生产环境中存储大量资金**
- 私钥以明文形式存储在数据库中，仅适用于开发和测试
- 生产环境请使用硬件钱包或专业的密钥管理服务
- 定期备份钱包数据和私钥信息

✅ **安全优势**：
- **离线签名架构**: 所有交易均在本地完成签名，私钥不会通过网络传输
- **多链签名支持**: 支持 EVM (ECDSA)、Bitcoin (ECDSA)、Solana (Ed25519) 等不同签名算法
- **零第三方依赖**: 无需依赖外部钱包服务或签名提供商
- **完全自主控制**: 应用拥有签名过程的完全控制权

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request



---

**⭐ 如果这个项目对你有帮助，请给它一个 Star！**
