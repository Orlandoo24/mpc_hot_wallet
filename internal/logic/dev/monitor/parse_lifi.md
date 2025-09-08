LI.FI 可以帮助您解析交易操作类型！我们提供了一个专门的 API 端点来解析交易调用数据。

## 解析交易数据的 API

LI.FI 提供了 `/v1/calldata/parse` 端点（目前为 BETA 版本），可以解析交易调用数据并提供 JSON 格式的内部交易信息概览。

### 使用方法：

```http
GET /v1/calldata/parse?chainId={chainId}&callData={callData}
```

**参数：**
- `chainId`: 交易所在的链 ID
- `callData`: 要解析的调用数据

这个端点会基于已知的和链上的 ABI 来解析调用数据，比您当前基于事件签名的方法更全面。

## 交易状态追踪

另外，LI.FI 还提供 `/v1/status` 端点来检查交易状态，可以通过交易哈希识别操作类型：

```http
GET /v1/status?txHash={txHash}&fromChain={chainId}
```

这会返回详细的交易信息，包括使用的工具（swap、bridge等）和操作类型。

## 优势

相比您现在的方法，LI.FI 的解决方案：
- 支持更多协议和操作类型
- 不依赖硬编码的事件签名
- 可以识别复杂的多步骤操作
- 提供标准化的操作分类

这样您就不需要维护大量的事件签名常量了。

```suggestions
(Parse transaction call data API)[/api-reference/parse-transaction-call-data-beta]
(Check transaction status)[/api-reference/check-the-status-of-a-cross-chain-transfer]
```