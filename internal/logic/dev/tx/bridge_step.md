很简单，**用原生库签名 + Li.Fi API 做一次跨链**，通常就这几步：

1. **拿报价（决定走哪条桥/DEX & 估算 gas）**

    * `GET /v1/quote`（单步）或 `GET /v1/routes`→选一条→对该步 `GET /v1/stepTransaction`（多步）。
    * 传入 `fromChain / toChain / fromToken / toToken / fromAmount (/toAddress)`，可加 `order=FASTEST|CHEAPEST`、`slippage`。
    * 响应里会给出 **`transactionRequest`**（to、data、value、gasLimit、gasPrice…）和（如需）**`approvalAddress`**。

2. **（可选）批准 ERC-20**

    * 若 `fromToken` 不是原生币，检查 allowance；不足则对 **`approvalAddress`** 发一笔 `approve(amount)`。
    * 用你自己的库（go-ethereum/ethers/viem）**本地签名 + 发送**，等 1–N 个确认。

3. **发送跨链交易**

    * 直接用 `transactionRequest` 字段构造交易（链是 **源链**），**本地签名 + 广播**。
    * `gasPrice/gasLimit` 若返回里没有，就用节点 `eth_estimateGas / suggestGasPrice` 回退。

4. **跟踪进度直到落地**

    * 用 `GET /v1/status?txHash=<源链txHash>` 轮询，直到状态变为 DONE（或拿到目标链 txLink）。
    * （部分路线会在目标链继续自动 swap / 合约调用，由桥的 relayer完成，你这边只需等状态完成）

5. **（可选）目的链后续动作**

    * 若你的路线设计了“到达后再做某合约调用/换币”，SDK/Routes 会把它作为下一步（仍通过 `stepTransaction` 生成 tx，由 relayer或你执行；视路线而定）。

---

### 极简伪流程

```txt
quote = GET /v1/quote(...)
if quote.estimate.approvalAddress && allowance < fromAmount:
    tx1 = buildApproveTx(spender=approvalAddress, amount=∞)  // 本地签名+发送
    waitReceipt(tx1)

tx2 = buildTx(quote.transactionRequest) // 本地签名+发送到源链
waitReceipt(tx2)

loop:
    s = GET /v1/status?txHash=tx2.hash
    if s.status in {DONE, FAILED}: break
```

**要点**

* **同链同币(如 BSC→BSC, BNB→BNB)** 属于普通转账，别用 `/quote`；直接 21,000 gas 裸转即可。
* 原生币地址用 `0xEeeee…`，ERC-20 用合约地址；金额按各自 `decimals` 换算。
* `toAddress` 可把落地资产直接打给目标链收款人。

如果你给我一组具体参数（链/币/数量/接收地址），我可以按这 4 步出一份可直接粘贴运行的 cURL + Go 调用示例。
