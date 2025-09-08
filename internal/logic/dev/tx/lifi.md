好的！你想实现一笔完整操作：使用一个 EVM 私钥钱包（具有 address + 私钥），在 BSC 上 **先执行 approve，再 send（比如 swap 或桥接）**。下面我帮你从 LI.FI 的文档中筛选出你最需要的接口与说明。

---

## 需要的文档链接与API功能

### 1. LI.FI SDK — 授权管理（Approve 操作）

* 文档页面：**Token Management – LI.FI SDK**
  包含一系列有关代币授权的接口：

    * `getTokenAllowance`：查询当前授权额度
    * `setTokenAllowance`：发起 approve
    * `revokeTokenApproval`：撤销授权
      ([LI.FI][1])

### 2. LI.FI API — 获取交易数据（Approve + Send）

* 文档页面：**Overview – LI.FI API Reference**
  提供 Base-URL（`https://li.quest/v1`）及授权方式说明（可用 `x-lifi-api-key` 提高速率）
  ([LI.FI][2])

* API 接口：

    * **`/quote`**：获取单步交易（同链 swap 或跨链 bridge）的报价和 `transactionRequest`，里面包含发送交易所需的签名数据。
    * **`/quote/contractCalls`**（Beta，多重调用）：可以组合多个合约调用，支持跨链 + 执行多个步骤合并。
      ([LI.FI][3], [LI.FI帮助中心][4])

### 3. FAQ — 为什么 API 不直接提供 approve 数据？

* 文档页面：**LI.FI 技术帮助中心问答**
  说明：API 不额外暴露 approve 的交易数据，建议通过 SDK 来获取和执行授权操作。
  ([LI.FI帮助中心][5])

---

## 使用流程总览

| 操作步骤                 | 方法                                                                     | 来源               |
| -------------------- | ---------------------------------------------------------------------- | ---------------- |
| 查询/设置授权（approve）     | 使用 LI.FI SDK 的 `getTokenAllowance` / `setTokenAllowance`               | ([LI.FI][1])     |
| 获取执行交易信息             | 调用 LI.FI API `/quote` 或 `/quote/contractCalls` 获取 `transactionRequest` | ([LI.FI][2])     |
| 执行 approve + send    | 使用钱包私钥签名并发送授权后，再按返回的 `transactionRequest` 发送交易                         | —                |
| 执行整体流程（可参考 SDK 自动流程） | LI.FI SDK 内部的 `executeRoute` 会自动处理授权与交易逻辑                              | ([LI.FI帮助中心][4]) |

---

**核心说明：**

* LI.FI 的 HTTP API 专注获取交易数据，不处理 private key 签名，也不提供 approve tx 构建接口。
* 授权步骤建议使用 LI.FI SDK 完成（这些操作通常会生成签名交易数据留给你）。
* 主要 API 接口为 `/quote` 或 `/quote/contractCalls`，它们会包含 send 操作的数据。

---

如果你需要，我可以继续帮你写一个 **Go 示例（使用 EVM-Go SDK + LI.FI）**，完成 authorize（approve）+ fetch transactionRequest 然后私钥签名 + 发送，整个 end-to-end 的流程。

[1]: https://docs.li.fi/sdk/token-management?utm_source=chatgpt.com "Token Management"
[2]: https://docs.li.fi/api-reference/introduction?utm_source=chatgpt.com "Overview - LI.FI"
[3]: https://docs.li.fi/api-reference/perform-a-contract-call-across-blockchains-beta?utm_source=chatgpt.com "Perform a contract call across blockchains (BETA)"
[4]: https://help.li.fi/hc/en-us/articles/16509431041307-How-do-I-use-the-API-endpoints-and-interprete-the-responses?utm_source=chatgpt.com "How do I use the API endpoints and interprete ..."
[5]: https://help.li.fi/hc/en-us/articles/13305031013275--How-can-I-find-get-approve-tx-data-and-check-for-allowance?utm_source=chatgpt.com "How can I find get approve tx data and check for allowance?"
