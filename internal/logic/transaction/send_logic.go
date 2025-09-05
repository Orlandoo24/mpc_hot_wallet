package transaction

import (
	"context"
	"demo/internal/svc"
	"demo/internal/types"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	evmTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/zeromicro/go-zero/core/logx"
)

type TransactionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTransactionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransactionLogic {
	return &TransactionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// WrapSend 封装了授权和转账操作的 集合 send ， 测试 rpc ：https://bsc-testnet-rpc.publicnode.com
func (l *TransactionLogic) WrapSend(req *types.TransactionSendReq) (resp *types.TransactionSendResp, err error) {
	l.Infof("--- 开始处理 /transaction/send 请求 for address %s ---", req.FromAddress)

	// 1. Get Chain Configuration
	l.Infof("步骤 1: 获取链配置 for chain: %s", req.FromChain)
	chainConfig, ok := l.svcCtx.Config.Chains[req.FromChain]
	if !ok {
		l.Errorf("获取链配置失败: 不支持的链 %s", req.FromChain)
		return nil, fmt.Errorf("unsupported chain: %s", req.FromChain)
	}
	l.Infof("链配置获取成功: ChainId=%d, RpcUrl=%s", chainConfig.ChainId, chainConfig.RpcUrl)

	// 2. Connect to RPC client
	l.Infof("步骤 2: 连接到 RPC 节点...")
	client, err := ethclient.Dial(chainConfig.RpcUrl)
	if err != nil {
		l.Errorf("RPC 节点连接失败: %v", err)
		return nil, errors.New("failed to connect to chain")
	}
	defer client.Close()
	l.Infof("RPC 节点连接成功.")

	// 3. Fetch wallet from DB and get private key
	l.Infof("步骤 3: 从数据库查询钱包地址: %s", req.FromAddress)
	wallet, err := l.svcCtx.WalletsDao.FindOneByAddress(l.ctx, req.FromAddress)
	if err != nil {
		l.Errorf("数据库查询失败 for address %s: %v", req.FromAddress, err)
		return nil, errors.New("wallet not found")
	}
	l.Infof("钱包查询成功.")

	// !!! WARNING: The private key should be decrypted here.
	// For this demo, we assume it's stored as a hex string.
	privateKey, err := crypto.HexToECDSA(wallet.EncryptedPrivateKey)
	if err != nil {
		l.Errorf("私钥解析失败: %v", err)
		return nil, errors.New("invalid private key")
	}
	l.Infof("私钥解析成功.")

	// 4. Call LI.FI API to get transaction data
	// Use the same chain ID as our RPC connection to avoid contract mismatch
	lifiFromChain := int64(chainConfig.ChainId) // Use actual chain ID from config
	lifiToChain := int64(chainConfig.ChainId)   // Use actual chain ID from config

	// For testing, let's do a BNB to USDT swap instead of same-token transfer
	fromToken := req.FromToken
	toToken := req.ToToken

	// If both tokens are native (0x0000...), change toToken to USDT for a real swap
	if req.FromToken == "0x0000000000000000000000000000000000000000" && req.ToToken == "0x0000000000000000000000000000000000000000" {
		toToken = "0x55d398326f99059fF775485246999027B3197955" // USDT on BSC mainnet
		l.Infof("Converting native-to-native transfer to BNB-to-USDT swap for LI.FI compatibility")
	}

	quoteURL := fmt.Sprintf(
		"%s/quote?fromChain=%d&toChain=%d&fromToken=%s&toToken=%s&fromAddress=%s&toAddress=%s&fromAmount=%s",
		l.svcCtx.Config.Lifi.ApiUrl,
		lifiFromChain,
		lifiToChain,
		fromToken,
		toToken,
		req.FromAddress,
		req.ToAddress,
		req.Amount,
	)
	l.Infof("步骤 4: 调用 LI.FI API 获取报价, URL: %s", quoteURL)

	httpResp, err := http.Get(quoteURL)
	if err != nil {
		l.Errorf("LI.FI API 调用失败: %v", err)
		return nil, errors.New("failed to get quote")
	}
	defer httpResp.Body.Close()
	l.Infof("LI.FI API 调用成功, 状态码: %d", httpResp.StatusCode)

	body, _ := ioutil.ReadAll(httpResp.Body)
	l.Infof("LI.FI API 原始响应: %s", string(body))

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		l.Errorf("LI.FI 响应解析失败: %v", err)
		return nil, errors.New("failed to parse quote response")
	}
	l.Infof("LI.FI API 解析后的响应结构: %+v", result)

	txReq, ok := result["transactionRequest"].(map[string]interface{})
	if !ok {
		l.Errorf("LI.FI 响应格式错误，找不到 transactionRequest 字段: %s", string(body))
		return nil, fmt.Errorf("invalid quote response")
	}
	l.Infof("LI.FI 响应解析成功.")

	// 5. Build, sign, and send the transaction
	l.Infof("步骤 5: 构建、签名并发送交易...")
	to := common.HexToAddress(txReq["to"].(string))
	value := new(big.Int)
	value.SetString(txReq["value"].(string), 0) // 0 means auto-detect hex/decimal
	data := common.FromHex(txReq["data"].(string))

	// Parse gasLimit from hex string
	gasLimitBig := new(big.Int)
	gasLimitBig.SetString(txReq["gasLimit"].(string), 0) // 0 means auto-detect hex/decimal
	gasLimit := gasLimitBig.Uint64()

	// Parse gasPrice from hex string
	gasPriceBig := new(big.Int)
	gasPriceBig.SetString(txReq["gasPrice"].(string), 0) // 0 means auto-detect hex/decimal

	nonce, _ := client.PendingNonceAt(l.ctx, common.HexToAddress(req.FromAddress))

	tx := evmTypes.NewTransaction(nonce, to, value, gasLimit, gasPriceBig, data)
	signedTx, err := evmTypes.SignTx(tx, evmTypes.NewEIP155Signer(big.NewInt(chainConfig.ChainId)), privateKey)
	if err != nil {
		l.Errorf("交易签名失败: %v", err)
		return nil, errors.New("failed to sign transaction")
	}
	l.Infof("交易签名成功, TxHash: %s", signedTx.Hash().Hex())

	err = client.SendTransaction(l.ctx, signedTx)
	if err != nil {
		l.Errorf("交易发送失败: %v", err)
		return nil, errors.New("failed to send transaction")
	}
	l.Infof("交易发送成功!")

	resp = &types.TransactionSendResp{
		TxHash: signedTx.Hash().Hex(),
	}
	l.Infof("--- /transaction/send 请求处理完成, 返回 TxHash: %s ---", resp.TxHash)
	return resp, nil
}
