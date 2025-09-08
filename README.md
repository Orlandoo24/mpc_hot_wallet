# MPC Demo - 本地热钱包演示项目

一个基于 Go-Zero 框架构建的多功能区块链钱包 API 服务，集成了 LI.FI 跨链协议和实时区块链监控功能，提供完整的钱包操作、代币交换、跨链转账、授权管理和交易监控功能。

## ✨ 主要功能

### 🔐 钱包管理
- **钱包初始化** - 创建和管理以太坊兼容钱包
- **私钥管理** - 安全的私钥存储和检索
- **多链支持** - 支持 BSC、ETH、Polygon、Arbitrum、Optimism、Base 等主流链

### 💸 交易功能
- **原生转账** - 支持 ETH、BNB 等原生代币转账
- **ERC20 转账** - 支持所有标准 ERC20 代币转账
- **智能 Gas 估算** - 动态 Gas 价格和限制优化
- **交易状态追踪** - 完整的交易确认和状态监控

### 🔄 代币交换 (Swap)
- **LI.FI 集成** - 使用 LI.FI 协议获得最优交换路由
- **智能 Approve** - 自动检查和管理代币授权
- **滑点保护** - 可配置的滑点保护机制
- **多 DEX 聚合** - 自动选择最优的去中心化交易所

### 🌉 跨链转账 (Bridge)
- **多链桥接** - 支持主流区块链之间的资产跨链
- **LI.FI 最佳实践** - 集成官方推荐的优化参数
- **状态轮询** - 实时跟踪跨链交易进度
- **重试机制** - 自动处理网络异常和交易失败

### 🔑 授权管理 (Approve)
- **授权检查** - 查询代币对特定合约的授权额度
- **智能授权** - 支持无限授权和指定金额授权
- **批量查询** - 一次查询多个代币的授权状态
- **授权撤销** - 安全撤销不需要的代币授权

### 📊 实时区块链监控
- **WebSocket 连接** - 实时监控 BSC 测试网区块链交易
- **智能事件解析** - 自动识别 Transfer、Approval、Swap、Bridge 等 EVM 事件
- **方向识别** - 智能标记交易方向（IN/OUT），精确识别资金流向
- **LI.FI 增强** - 集成 LI.FI API 进行高级交易类型识别
- **Kafka 集成** - 将解析的事件数据推送到 Kafka 消息队列
- **自动重连** - WebSocket 连接异常时自动重连机制
- **钱包过滤** - 只监控数据库中存储的钱包地址相关交易

## 🚀 快速开始

### 环境要求
- Go 1.19+
- PostgreSQL 12+
- 稳定的网络连接（用于 RPC 调用）

### 安装步骤

1. **克隆项目**
```bash
git clone <repository-url>
cd mpc_demo
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

### 钱包管理

#### 初始化钱包
```http
POST /api/wallet_init
Content-Type: application/json

{
  "address": "0x...",
  "private_key": "0x..."
}
```

### 交易操作

#### 发送交易
```http
POST /api/transaction/send
Content-Type: application/json

{
  "from_address": "0x...",
  "to_address": "0x...",
  "chain": "BSC",
  "from_token": "0x0000000000000000000000000000000000000000",
  "to_token": "0x0000000000000000000000000000000000000000",
  "amount": "1000000000000000000"
}
```

#### 代币交换
```http
POST /api/transaction/swap
Content-Type: application/json

{
  "from_address": "0x...",
  "to_address": "0x...",
  "chain": "BSC",
  "from_token": "0x0000000000000000000000000000000000000000",
  "to_token": "0x55d398326f99059fF775485246999027B3197955",
  "amount": "10000000000000"
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

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📄 许可证

本项目采用 MIT 许可证。

## 🆘 常见问题

### Q: 如何添加新的区块链支持？
A: 在 `etc/demo.yaml` 中添加新的链配置，并在相关逻辑文件中添加链ID映射。

### Q: 交易失败怎么办？
A: 检查 Gas 价格、网络状态和账户余额。系统包含自动重试机制。

### Q: 如何更换 RPC 节点？
A: 修改 `etc/demo.yaml` 中对应链的 `RpcUrl` 配置。

### Q: 跨链交易卡住了怎么办？
A: 使用 `/api/bridge/status` 接口查询交易状态，LI.FI 协议会自动处理大部分异常情况。

### Q: 区块链监控功能如何工作？
A: 系统会自动连接 BSC 测试网 WebSocket，实时监控数据库中钱包地址的相关交易，解析事件并推送到 Kafka。

### Q: 如何查看监控到的交易事件？
A: 监控事件会在控制台输出，格式为：`📥 EVM事件: Transfer-IN | 金额: 0.00001 BNB | 哈希: 0x1234...`

### Q: 监控服务断开了怎么办？
A: 系统具备自动重连机制，WebSocket 连接异常时会自动尝试重连。

---

**⭐ 如果这个项目对你有帮助，请给它一个 Star！**
