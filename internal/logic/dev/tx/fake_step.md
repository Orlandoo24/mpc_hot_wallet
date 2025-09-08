好的，我们来将您提供的这份精炼、专业的跨链操作流程，整理成一个更加结构化和易于理解的指南。

### **使用原生库签名 + LI.FI API 完成跨链交易的完整指南**

本指南详细描述了如何结合本地钱包签名与 LI.FI 的强大 API 来安全、可靠地执行一次跨链资产转移。整个过程的核心思想是：**LI.FI 负责计算最优路径并生成交易参数，而私钥签名和交易广播则完全在您本地完成，确保了私钥的绝对安全。**

-----

#### **第一步：获取报价与交易方案 (Get Quote)**

这是所有操作的起点。我们需要向 LI.FI 请求一个最优的跨链方案，它会综合考虑费用、速度和滑点，并返回执行该方案所需的具体交易数据。

* **API 端点:**

    * **单步交易:** `GET /v1/quote` - 适用于大部分简单的 A→B 跨链场景。
    * **多步交易:** `GET /v1/routes` → 选择一条路由 → `GET /v1/stepTransaction` - 适用于更复杂的路径，例如需要多次兑换或经过多个协议的场景。

* **核心请求参数:**

    * `fromChain`: 源链的 ID 或名称 (例如: `1`, `ETH`)。
    * `toChain`: 目标链的 ID 或名称 (例如: `137`, `POL`)。
    * `fromToken`: 源代币的合约地址 (原生币使用 `0x0000000000000000000000000000000000000000` 或 `0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE`)。
    * `toToken`: 目标代币的合约地址。
    * `fromAmount`: 要转移的源代币数量 (注意：需要是 **wei** 格式，即根据代币的 `decimals` 进行换算)。
    * `toAddress` (可选): 在目标链上接收资产的地址。如果留空，则默认为发起交易的地址 (`fromAddress`)。

* **可选优化参数:**

    * `order`: `FASTEST` (最快) | `CHEAPEST` (最便宜) | `RECOMMENDED` (推荐)。
    * `slippage`: 可接受的滑点百分比 (例如: `0.03` 代表 3%)。

* **关键响应内容:**

    * `transactionRequest`: 一个 JSON 对象，包含了在 **源链** 上发起交易所需的所有参数，如 `to` (交互的合约地址), `data` (已编码的函数调用数据), `value` (发送的原生币数量), `gasLimit`, `gasPrice` 等。
    * `estimate.approvalAddress`: 如果转移的是 ERC-20 代币，这里会提供需要进行 `approve` 授权的合约地址（通常是 LI.FI 的代理合约或桥接合约）。

-----

#### **第二步：(可选) 批准 ERC-20 代币授权 (Approve ERC-20)**

如果您的 `fromToken` 是一个 ERC-20 代币（而不是 ETH、BNB 等原生币），您必须先授权给 LI.FI 的智能合约，允许其从您的钱包中转移指定数量的代币。

1.  **检查授权额度 (Allowance):** 调用该 ERC-20 合约的 `allowance(owner, spender)` 方法，其中 `owner` 是您的钱包地址，`spender` 是上一步 `quote` 响应中返回的 `approvalAddress`。
2.  **发起授权交易:**
    * 如果现有额度低于 `fromAmount`，您需要构建一笔 `approve` 交易。
    * **目标合约:** `fromToken` 的合约地址。
    * **函数调用:** `approve(spender, amount)`，其中 `spender` 是 `approvalAddress`，`amount` 通常设置为一个极大值（例如 `uint256 max`）以避免未来重复授权。
3.  **本地签名与发送:** 使用您的原生库 (如 ethers.js, viem, go-ethereum) 在本地对这笔 `approve` 交易进行签名，然后将其广播到 **源链** 网络。
4.  **等待确认:** 等待该授权交易被矿工打包确认（至少 1 个区块确认）。

-----

#### **第三步：发送主跨链交易 (Send Cross-Chain Transaction)**

现在，所有前置条件都已满足，可以正式发起跨链交易了。

1.  **构建交易:** 直接使用第一步 `quote` 响应中获取的 `transactionRequest` 对象来构建您的交易。所有字段（`to`, `data`, `value` 等）都已经为您准备好了。
2.  **Gas 费用回退机制:** 如果 `transactionRequest` 中没有提供 `gasPrice` 或 `gasLimit`，您应该调用连接节点的 RPC 方法（如 `eth_estimateGas` 和 `eth_gasPrice`）来获取估算值，以确保交易能成功上链。
3.  **本地签名与广播:** 同样地，使用您的原生库在本地对这笔主交易进行签名，并将其广播到 **源链** 网络。
4.  **记录交易哈希:** 交易广播成功后，您会得到一个交易哈希 (`txHash`)。**这个哈希是跟踪后续进度的唯一凭证。**

-----

#### **第四步：轮询跟踪交易状态 (Track Progress)**

跨链交易并非瞬间完成，您需要持续跟踪其状态，直到资产在目标链上成功落地。

* **API 端点:** `GET /v1/status?txHash=<源链txHash>`
* **轮询逻辑:**
    * 使用第三步获取的 `txHash`，定期（例如每 10-15 秒）调用此 API。
    * 检查响应中的 `status` 字段。
    * `PENDING`: 交易仍在处理中，继续轮询。
    * `DONE`: 交易已成功完成！响应中通常会包含目标链的交易链接 (`receiving.txLink`)。
    * `FAILED`: 交易失败。
* **完成标志:** 当状态变为 `DONE` 时，代表整个跨链流程已结束，资产已到达 `toAddress`。

-----

#### **第五步：(可选) 目的链的后续操作 (Post-Arrival Actions)**

在某些复杂的路由设计中（例如，跨链后立即进行一次兑换），LI.FI 的 `routes` 可能会包含多个步骤。如果当前步骤完成后还有下一步，您需要：

1.  再次调用 `GET /v1/stepTransaction` 获取下一步的 `transactionRequest`。
2.  根据路由设计，这个 `transactionRequest` 可能由您自己执行，也可能由桥的 Relayer 自动完成。您只需继续通过 `status` 接口跟踪最终状态即可。

-----

### **极简伪代码流程**

下面是一个简化的代码逻辑，清晰地展示了上述步骤的交互流程。

```
// 1. 获取报价
quote = GET_request("https://li.quest/v1/quote", {fromChain, toChain, ...})

// 2. 检查并进行 ERC-20 授权 (如果需要)
if quote.estimate.approvalAddress:
    current_allowance = check_allowance(owner=my_address, spender=quote.estimate.approvalAddress)
    if current_allowance < fromAmount:
        // 构建、本地签名并发送 approve 交易
        approve_tx = build_approve_transaction(
            spender=quote.estimate.approvalAddress,
            amount=MAX_UINT256
        )
        approve_tx_hash = sign_and_send(approve_tx)
        wait_for_confirmation(approve_tx_hash) // 等待交易上链

// 3. 发送主跨链交易
// 使用 quote 返回的参数构建交易
cross_chain_tx = build_transaction(quote.transactionRequest)
cross_chain_tx_hash = sign_and_send(cross_chain_tx) // 本地签名并发送到源链
wait_for_confirmation(cross_chain_tx_hash)

// 4. 轮询状态直到完成
while True:
    status_response = GET_request("https://li.quest/v1/status", {txHash: cross_chain_tx_hash})
    if status_response.status == "DONE":
        print("跨链成功! 目标链交易:", status_response.receiving.txLink)
        break
    elif status_response.status == "FAILED":
        print("跨链失败!")
        break
    sleep(15) // 等待 15 秒后再次查询
```