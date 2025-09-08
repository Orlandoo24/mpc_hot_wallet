好的，给你把 **LI.FI 的跨链 bridge 操作**相关文档和最小可跑的 Demo 摘出来了（只放“桥接”必需的那几页）：

# 文档链接 + 作用

* **`GET /v1/quote` – 获取单步跨链/换币的可执行交易（含 gas 估算、transactionRequest）**
这是发起桥接/换币的核心入口。返回的 `Step` 里带 `transactionRequest`，你直接签名发送即可。适用于从 A 链把某币换/桥到 B 链（也可同链换币）。([docs.li.fi][1])

* **Routes vs Quote 的区别 & 何时需要 `stepTransaction`**
如果你用 **Routes**（多步路由），需要在选定 route 后再对每个 step 调 `stepTransaction` 拿到可执行交易；而 **Quote** 已直接返回 `transactionRequest`，无需再调 `stepTransaction`。官方“请求 Route / 取 Quote”的说明和端到端流程里都强调了这一点。([docs.li.fi][2])

* **`GET /v1/status` – 跨链交易状态查询**
跨链可能耗时较长，用这个端点跟踪发送链/接收链的状态（传 `txHash` 即可）。([docs.li.fi][3])

* **`GET /v1/tools` – 当前可用的桥与聚合器列表**
用来了解/限制策略里允许的桥（allow/deny bridges）。([docs.li.fi][4])

* **Connections（可选，做“可连通性”/偏好桥清单）**
返回在 from/to 条件下可能的桥/交易所组合，并支持 allow/deny；/tools 返回的桥名称可配合这里使用。([docs.li.fi][5])

---

# Demo A：用 **/v1/quote** 做跨链 Bridge（cURL）

> 例：从 **BSC(56)** 把原生 **BNB** 跨到 **Base(8453)** 的 **原生 ETH**（只示范写法；金额和地址按需替换）

```bash
curl --request GET \
  --url 'https://li.quest/v1/quote?fromChain=56&toChain=8453&fromToken=0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE&toToken=0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE&fromAmount=10000000000000&fromAddress=0xYourSender&toAddress=0xRecipient&order=FASTEST&slippage=0.005'
```

响应里你会拿到：

* `estimate.gasCosts[]`：gas 估算
* `transactionRequest`：`to / data / value / gasLimit / gasPrice` —— 直接签名发送即可（原生币不需要 approve；若是 ERC-20，先按 `estimate.approvalAddress` 做 approve）。([docs.li.fi][1])

> 提示：**同链同币（BSC→BSC, BNB→BNB）不是 bridge/swap 场景**，`/quote` 会报错；那是“普通转账”，请直接用 RPC 构造 21,000 gas 的裸转账。([docs.li.fi][1])

---

# Demo B：发送 `transactionRequest`（Go，最小可跑）

```go
// go.mod 需要：github.com/ethereum/go-ethereum v1.x
package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TxReq struct {
	To       string `json:"to"`
	Data     string `json:"data"`
	Value    string `json:"value"`
	GasLimit string `json:"gasLimit"`
	GasPrice string `json:"gasPrice"`
}
type QuoteResp struct {
	TransactionRequest TxReq `json:"transactionRequest"`
}

func main() {
	// 1) 调用 /v1/quote（与 Demo A 相同的参数，按需替换）
	quoteURL := "https://li.quest/v1/quote?fromChain=56&toChain=8453&" +
		"fromToken=0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE&toToken=0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE&" +
		"fromAmount=10000000000000&fromAddress=0xYourSender&toAddress=0xRecipient&order=FASTEST&slippage=0.005"

	resp, err := http.Get(quoteURL) // 可在 Header 里加 x-lifi-api-key 提高限额
	if err != nil { panic(err) }
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 { panic(string(body)) }

	var q QuoteResp
	if err := json.Unmarshal(body, &q); err != nil { panic(err) }

	// 2) 解析 transactionRequest
	to := common.HexToAddress(q.TransactionRequest.To)
	data := common.FromHex(q.TransactionRequest.Data)

	val := new(big.Int)
	if q.TransactionRequest.Value != "" {
		// 字段可能是十六进制或十进制，LI.FI 常返回十进制字符串
		_, ok := val.SetString(q.TransactionRequest.Value, 10)
		if !ok { _, _ = val.SetString(q.TransactionRequest.Value, 0) }
	}

	gasLimit := uint64(300000)
	if q.TransactionRequest.GasLimit != "" {
		if gl, ok := new(big.Int).SetString(q.TransactionRequest.GasLimit, 10); ok {
			gasLimit = gl.Uint64()
		}
	}

	gasPrice := new(big.Int)
	if q.TransactionRequest.GasPrice != "" {
		_, _ = gasPrice.SetString(q.TransactionRequest.GasPrice, 10)
	}

	// 3) 连接发送链（这里以 BSC 为例；如果第一步返回的是源链 tx，就连源链）
	client, _ := ethclient.Dial("https://bsc-dataseed.binance.org")
	defer client.Close()

	from := common.HexToAddress("0xYourSender")
	nonce, _ := client.PendingNonceAt(context.Background(), from)
	if gasPrice.Sign() == 0 {
		gp, _ := client.SuggestGasPrice(context.Background())
		gasPrice = gp
	}

	// 4) 构造、签名、发送
	tx := types.NewTransaction(nonce, to, val, gasLimit, gasPrice, data)

	// 私钥
	priv, _ := crypto.HexToECDSA("YOUR_PRIVATE_KEY_HEX") // 切勿明文硬编码，示例而已
	signed, _ := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(56)), priv) // 56=BSC
	if err := client.SendTransaction(context.Background(), signed); err != nil { panic(err) }

	fmt.Println("txHash:", signed.Hash().Hex())
}
```

> 说明：
>
> * 若 `fromToken` 是 **ERC-20**，发送前需要根据 `estimate.approvalAddress` 先做 `approve`（SDK/端到端示例里有完整流程）。([docs.li.fi][6])
> * `x-lifi-api-key` 不是必需，但能提升速率限制（只放在**服务端**）。([docs.li.fi][7])

---

# Demo C：查询跨链状态（cURL）

```bash
# 已知发送链 txHash（或接收链 txHash / transactionId 其一即可）
curl --request GET \
  --url 'https://li.quest/v1/status?txHash=0xYourTxHash'
```

返回会包含 `status` / `substatus`、发送/接收链的 txLink 等信息。([docs.li.fi][3])

---

# 常见坑位（bridge 场景）

* **同链同币**不是 quote 的目标（BSC→BSC 的 BNB→BNB），会被拒；纯转账用 RPC。([docs.li.fi][1])
* 原生币 vs ERC-20 的 **`fromToken`**：原生用 `0xEeeee...`；ERC-20 必须用合约地址，并按 `decimals` 提供 `fromAmount`。([docs.li.fi][1])
* 需要**限制或偏好桥**时，用 `/tools` 拿桥清单，并在请求里用 allow/deny 参数（见 Connections 说明）。([docs.li.fi][4])

如果你告诉我具体的「fromChain / toChain / token / 金额」，我可以把上面的 Demo 参数替换成你要测试的精确组合，顺带把 `approve` 流程也接上。

[1]: https://docs.li.fi/api-reference/get-a-quote-for-a-token-transfer?utm_source=chatgpt.com "Get a quote for a token transfer"
[2]: https://docs.li.fi/introduction/user-flows-and-examples/requesting-route-fetching-quote?utm_source=chatgpt.com "Fetching a Quote/Route"
[3]: https://docs.li.fi/api-reference/check-the-status-of-a-cross-chain-transfer?utm_source=chatgpt.com "Check the status of a cross chain transfer"
[4]: https://docs.li.fi/api-reference/get-available-bridges-and-exchanges?utm_source=chatgpt.com "Get available bridges and exchanges"
[5]: https://docs.li.fi/api-reference/returns-all-possible-connections-based-on-a-from-or-tochain?utm_source=chatgpt.com "Returns all possible connections based on a from"
[6]: https://docs.li.fi/introduction/user-flows-and-examples/end-to-end-example?utm_source=chatgpt.com "End-to-end Transaction Example"
[7]: https://docs.li.fi/api-reference/introduction?utm_source=chatgpt.com "Overview - LI.FI"
