好的，我已经为您整理了 LI.FI 跨链操作 API 的详细文档内容。

### **LI.FI 跨链桥接 API 详细文档**

-----

### **核心 API (Core APIs)**

#### **1. 获取可用的桥接和交易所 (Get Available Bridges and Exchanges)**

此端点用于获取 LI.FI 服务可用的桥接和交易所的信息。

* **HTTP 请求:** `GET https://li.quest/v1/tools`
* **功能:** 返回 LI.FI 当前集成的所有可用工具（桥接、DEXes 等）的列表。
* **用途:** 了解 LI.FI 的底层流动性来源。

**响应 (200 OK):**

响应中包含 `bridges` 和 `exchanges` 两个对象数组。

* `bridges`:
    * `key`: 桥接的唯一标识符 (例如: "hop", "cbridge", "across")。
    * `name`: 桥接的通用名称 (例如: "Connext", "Hop Protocol")。
    * `logoURI`: 桥接的徽标 URL。
    * `supportedChains`: 此桥接支持的链 ID 数组。
* `exchanges`:
    * `key`: 交易所的唯一标识符 (例如: "1inch", "paraswap", "0x")。
    * `name`: 交易所的通用名称 (例如: "0x", "1inch")。
    * `logoURI`: 交易所的徽标 URL。
    * `supportedChains`: 此交易所支持的链 ID 数组。

**示例 cURL 请求:**

```bash
curl --request GET \
     --url https://li.quest/v1/tools
```

-----

#### **2. 获取跨链转账报价 (Get a Quote for a Token Transfer)**

此端点用于请求将一种代币兑换为另一种代币的报价，支持跨链或同链操作。

* **HTTP 请求:** `GET https://li.quest/v1/quote`
* **功能:** 根据起始链、目标链、代币和金额，计算出最优的跨链转账路径和报价。
* **用途:** 为用户提供跨链交易的预估结果和执行方案。

**查询参数:**

* `fromChain` (必需): 起始链的链 ID 或名称 (例如: "ETH", "137")。
* `toChain` (必需): 目标链的链 ID 或名称 (例如: "POL", "1")。
* `fromToken` (必需): 起始代币的地址。
* `toToken` (必需): 目标代币的地址。
* `fromAmount` (必需): 要兑换的起始代币数量（以最小单位表示）。
* `fromAddress` (可选): 用户的钱包地址。

**响应 (200 OK):**

返回一个 `Step` 对象，其中包含预估结果和可以直接发送到钱包的 `transactionRequest`。`estimate` 字段包含预估的 `toAmount` 和 `toAmountMin`（包含滑点在内的最低保证到账金额）。

**示例 cURL 请求:**

```bash
curl --request GET \
     --url 'https://li.quest/v1/quote?fromChain=POL&toChain=ETH&fromToken=DAI&toToken=USDC&fromAmount=1000000000000000000'
```

-----

#### **3. 检查跨链转账状态 (Check the Status of a Cross-Chain Transfer)**

跨链转账可能需要一些时间才能完成。此端点用于检查转账的状态。

* **HTTP 请求:** `GET https://li.quest/v1/status`
* **功能:** 查询特定跨链交易的当前状态。
* **用途:** 向用户实时反馈交易进度，提升用户体验。

**查询参数:**

* `txHash`: 交易哈希。
* `fromChain`: 起始链的链 ID（可选，但建议提供以加快请求速度）。
* `toChain`: 目标链的链 ID（可选）。
* `bridge`: 使用的桥接名称（可选）。

**响应 (200 OK):**

* **重要提示:** 即使找不到交易，此端点也会返回 200 成功响应。这是为了处理交易哈希有效但尚未被矿工打包的情况。
* `status`: 转账的当前状态，可能的值为 `PENDING`、`DONE`、`NOT_FOUND` 或 `FAILED`。
* `substatus`: 更具体的子状态，可用于 `PENDING` 和 `DONE` 状态。
* `tool`: 此转账使用的工具。
* `fromAddress`: 发送方地址。
* `toAddress`: 接收方地址。

**示例 cURL 请求:**

```bash
curl --request GET \
     --url 'https://li.quest/v1/status?txHash=0x...'
```

-----

#### **4. 获取链间连接 (Get Cross-Chain Connections)**

此端点提供有关链之间所有可能转账的信息。

* **HTTP 请求:** `GET https://li.quest/v1/connections`
* **功能:** 根据指定的起始链 (`fromChain`) 或目标链 (`toChain`)，返回所有可能的链间连接组合。
* **用途:** 用于构建用户界面，例如动态地显示支持的目标链下拉列表。

**查询参数:**

由于结果可能非常大，因此需要通过以下至少一个参数进行过滤：

* `fromChain`: 起始链的链 ID。
* `toChain`: 目标链的链 ID。
* `fromToken`: 起始代币的地址。
* `toToken`: 目标代币的地址。

**响应 (200 OK):**

返回一个包含 `connections` 数组的对象，其中每个连接对象都包含 `fromChainId`、`toChainId` 以及支持的 `fromTokens` 和 `toTokens`。

**示例 cURL 请求:**

```bash
curl --request GET \
     --url 'https://li.quest/v1/connections?fromChain=137&toChain=1'
```

-----

### **辅助 API (Auxiliary APIs)**

#### **1. 获取支持的链信息 (Get Supported Chains)**

* **HTTP 请求:** `GET https://li.quest/v1/chains`
* **功能:** 返回 LI.FI 平台支持的所有区块链的详细信息列表。

**响应 (200 OK):**

返回一个包含 `chains` 数组的对象，其中每个链对象都包含 `key`、`name`、`chainType`、`coin`、`id`、`mainnet`、`logoURI` 等信息。

**示例 cURL 请求:**

```bash
curl --request GET \
     --url https://li.quest/v1/chains
```

-----

#### **2. 获取代币信息 (Fetch All Known Tokens)**

* **HTTP 请求:** `GET https://li.quest/v1/tokens`
* **功能:** 返回 LI.FI 已知的所有代币的列表。

**查询参数:**

* `chains` (可选): 将结果限制在给定的链上。
* `chainTypes` (可选): 将结果限制在给定的链类型上。
* `minPriceUSD` (可选): 按最低美元价格过滤结果。

**响应 (200 OK):**

返回一个对象，其中键是链 ID，值是该链上代币对象的数组。每个代币对象都包含 `address`、`decimals`、`symbol`、`chainId`、`name`、`coinKey`、`logoURI` 和 `priceUSD` 等信息。

**示例 cURL 请求:**

```bash
curl --request GET \
     --url 'https://li.quest/v1/tokens?chains=1,137'
```