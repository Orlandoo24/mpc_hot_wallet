package main

import (
"context"
"crypto/ecdsa"
"encoding/json"
"fmt"
"log"
"math/big"
"net/http"
"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
rpcURL        = "https://bsc-dataseed.binance.org/" // BSC RPC
privateKeyHex = "YOUR_PRIVATE_KEY"                  // ⚠️ 替换为你的私钥
liFiAPI       = "https://li.quest/v1/quote"        // LI.FI Quote API
tokenAddress  = "0x55d398326f99059fF775485246999027B3197955" // USDT on BSC
spender       = "0x1111111254EEB25477B68fb85Ed929f73A960582" // 目标合约地址 (比如 1inch/LI.FI Router)
)

func main() {
// 1. 连接 BSC 节点
client, err := ethclient.Dial(rpcURL)
if err != nil {
log.Fatal(err)
}
defer client.Close()

	// 2. 导入私钥
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKey)
	fmt.Println("钱包地址:", fromAddress.Hex())

	// 3. 构建 approve 交易
	erc20ABI, _ := abi.JSON(strings.NewReader(`[{"constant":false,"inputs":[{"name":"spender","type":"address"},{"name":"amount","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"type":"function"}]`))

	amount := new(big.Int).Mul(big.NewInt(100), big.NewInt(1e18)) // 100 USDT
	data, _ := erc20ABI.Pack("approve", common.HexToAddress(spender), amount)

	nonce, _ := client.PendingNonceAt(context.Background(), fromAddress)
	gasPrice, _ := client.SuggestGasPrice(context.Background())

	tx := types.NewTransaction(nonce, common.HexToAddress(tokenAddress), big.NewInt(0), 100000, gasPrice, data)

	chainID, _ := client.NetworkID(context.Background())
	signedTx, _ := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Approve 交易已发送:", signedTx.Hash().Hex())

	// 4. 调用 LI.FI API 获取 send 的交易数据
	resp, err := http.Get(fmt.Sprintf("%s?fromChain=56&toChain=56&fromToken=%s&toToken=%s&fromAddress=%s&toAddress=%s&amount=1000000000000000000",
		liFiAPI,
		tokenAddress,
		"0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE", // 原生 BNB
		fromAddress.Hex(),
		"0x目标钱包地址", // ⚠️ 替换为接收方
	))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	// 从 result 里拿到 transactionRequest
	txReq := result["transactionRequest"].(map[string]interface{})
	to := common.HexToAddress(txReq["to"].(string))
	value := new(big.Int)
	value.SetString(txReq["value"].(string), 10)
	input := common.FromHex(txReq["data"].(string))

	// 5. 构建 send 交易
	nonce, _ = client.PendingNonceAt(context.Background(), fromAddress)
	gasLimit := uint64(500000) // 估算可替换

	sendTx := types.NewTransaction(nonce, to, value, gasLimit, gasPrice, input)
	signedSendTx, _ := types.SignTx(sendTx, types.NewEIP155Signer(chainID), privateKey)

	err = client.SendTransaction(context.Background(), signedSendTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Send 交易已发送:", signedSendTx.Hash().Hex())
}
