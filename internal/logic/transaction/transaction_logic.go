package transaction

import (
	"context"
	"crypto/ecdsa"
	"demo/internal/svc"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
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

// IsNativeToken 判断是否为原生代币
func (l *TransactionLogic) IsNativeToken(token string) bool {
	nativeTokens := []string{
		"0x0000000000000000000000000000000000000000",
		"0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
	}

	for _, native := range nativeTokens {
		if token == native {
			return true
		}
	}
	return false
}

// GetWalletPrivateKey 从数据库获取钱包私钥
func (l *TransactionLogic) GetWalletPrivateKey(fromAddress string) (*ecdsa.PrivateKey, error) {
	wallet, err := l.svcCtx.WalletsDao.FindOneByAddress(l.ctx, fromAddress)
	if err != nil {
		l.Errorf("查询钱包失败 for address %s: %v", fromAddress, err)
		return nil, errors.New("wallet not found")
	}

	privateKey, err := crypto.HexToECDSA(wallet.EncryptedPrivateKey)
	if err != nil {
		l.Errorf("私钥解析失败: %v", err)
		return nil, errors.New("invalid private key")
	}

	return privateKey, nil
}

// BuildERC20ApproveData 构建 ERC20 approve 函数的调用数据
func (l *TransactionLogic) BuildERC20ApproveData(spenderAddress string, amount *big.Int) []byte {
	// approve(address spender, uint256 amount)
	// 方法签名: 0x095ea7b3
	approveMethodId := []byte{0x09, 0x5e, 0xa7, 0xb3}

	spender := common.HexToAddress(spenderAddress)
	paddedSpender := common.LeftPadBytes(spender.Bytes(), 32)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	data := append(approveMethodId, paddedSpender...)
	data = append(data, paddedAmount...)

	return data
}

// ExecuteApproveTransaction 执行 ERC20 approve 交易
func (l *TransactionLogic) ExecuteApproveTransaction(client *ethclient.Client, privateKey *ecdsa.PrivateKey, tokenAddress, spenderAddress string, amount *big.Int, chainId int64) (string, error) {
	l.Infof("执行 ERC20 approve 操作，spender: %s", spenderAddress)

	// 构建 approve 调用数据
	data := l.BuildERC20ApproveData(spenderAddress, amount)

	// 获取交易参数
	fromAddr := crypto.PubkeyToAddress(privateKey.PublicKey)
	gasPrice, err := client.SuggestGasPrice(l.ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get gas price: %v", err)
	}

	// 估算 gas limit
	tokenAddr := common.HexToAddress(tokenAddress)
	gasLimit, err := client.EstimateGas(l.ctx, ethereum.CallMsg{
		From: fromAddr,
		To:   &tokenAddr,
		Data: data,
	})
	if err != nil {
		l.Infof("Gas 估算失败，使用默认值: %v", err)
		gasLimit = 100000
	}

	// 增加 gas limit 缓冲
	gasLimit = gasLimit * 120 / 100

	// 构建并发送交易
	return l.BuildAndSendTransaction(client, privateKey, tokenAddr, big.NewInt(0), data, gasLimit, gasPrice, chainId)
}

// BuildAndSendTransaction 构建并发送交易
func (l *TransactionLogic) BuildAndSendTransaction(client *ethclient.Client, privateKey *ecdsa.PrivateKey, to common.Address, value *big.Int, data []byte, gasLimit uint64, gasPrice *big.Int, chainId int64) (string, error) {
	// 获取 nonce
	fromAddr := crypto.PubkeyToAddress(privateKey.PublicKey)
	nonce, err := client.PendingNonceAt(l.ctx, fromAddr)
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %v", err)
	}

	// 构建交易（使用新的 NewTx 方法替代已弃用的 NewTransaction）
	tx := evmTypes.NewTx(&evmTypes.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    value,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	})

	// 签名交易
	signedTx, err := evmTypes.SignTx(tx, evmTypes.NewEIP155Signer(big.NewInt(chainId)), privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}

	// 发送交易
	err = client.SendTransaction(l.ctx, signedTx)
	if err != nil {
		// 检查错误信息中是否包含交易哈希（有些 RPC 节点会在错误信息中返回成功的交易哈希）
		if strings.Contains(err.Error(), "result") && strings.Contains(err.Error(), "0x") {
			l.Infof("⚠️ RPC 返回误导性错误，但交易可能已成功发送: %v", err)
			l.Infof("使用本地计算的交易哈希继续流程: %s", signedTx.Hash().Hex())
		} else {
			return "", fmt.Errorf("failed to send transaction: %v", err)
		}
	}

	return signedTx.Hash().Hex(), nil
}

// BuildExplorerUrl 根据链类型构建区块浏览器链接
func (l *TransactionLogic) BuildExplorerUrl(chain, txHash string) string {
	switch chain {
	case "BSC":
		return fmt.Sprintf("https://bscscan.com/tx/%s", txHash)
	case "BSC-TestNet":
		return fmt.Sprintf("https://testnet.bscscan.com/tx/%s", txHash)
	case "ETH":
		return fmt.Sprintf("https://etherscan.io/tx/%s", txHash)
	case "ETH-Goerli":
		return fmt.Sprintf("https://goerli.etherscan.io/tx/%s", txHash)
	case "ETH-Sepolia":
		return fmt.Sprintf("https://sepolia.etherscan.io/tx/%s", txHash)
	case "Polygon":
		return fmt.Sprintf("https://polygonscan.com/tx/%s", txHash)
	case "Polygon-Mumbai":
		return fmt.Sprintf("https://mumbai.polygonscan.com/tx/%s", txHash)
	default:
		l.Infof("未知链类型 %s，返回通用浏览器链接", chain)
		return fmt.Sprintf("https://explorer.example.com/tx/%s", txHash)
	}
}

// BuildERC20TransferData 构建 ERC20 transfer 调用数据
func (l *TransactionLogic) BuildERC20TransferData(toAddress string, amount *big.Int) ([]byte, error) {
	// transfer(address to, uint256 amount)
	// 方法签名: 0xa9059cbb
	transferMethodId := []byte{0xa9, 0x05, 0x9c, 0xbb}

	// 编码接收地址 (32 bytes)
	toAddr := common.HexToAddress(toAddress)
	paddedToAddress := common.LeftPadBytes(toAddr.Bytes(), 32)

	// 编码转账金额 (32 bytes)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)

	// 组合调用数据
	data := append(transferMethodId, paddedToAddress...)
	data = append(data, paddedAmount...)

	l.Infof("ERC20 transfer 数据构建成功，长度: %d bytes", len(data))
	return data, nil
}

// CheckAllowance 检查 ERC20 代币的 allowance
func (l *TransactionLogic) CheckAllowance(client *ethclient.Client, tokenAddress, owner, spender string) (*big.Int, error) {
	l.Infof("检查 allowance: token=%s, owner=%s, spender=%s", tokenAddress, owner, spender)

	// 构建 allowance 调用数据
	// allowance(address owner, address spender) returns (uint256)
	allowanceMethodId := []byte{0xdd, 0x62, 0xed, 0x3e}

	ownerAddr := common.HexToAddress(owner)
	spenderAddr := common.HexToAddress(spender)

	paddedOwner := common.LeftPadBytes(ownerAddr.Bytes(), 32)
	paddedSpender := common.LeftPadBytes(spenderAddr.Bytes(), 32)

	data := append(allowanceMethodId, paddedOwner...)
	data = append(data, paddedSpender...)

	// 调用合约（带重试机制）
	tokenAddr := common.HexToAddress(tokenAddress)

	var result []byte
	var err error
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		result, err = client.CallContract(l.ctx, ethereum.CallMsg{
			To:   &tokenAddr,
			Data: data,
		}, nil)

		if err == nil {
			break
		}

		l.Infof("Allowance 查询失败 (尝试 %d/%d): %v", i+1, maxRetries, err)
		if i < maxRetries-1 {
			time.Sleep(time.Duration(i+1) * time.Second) // 递增延迟
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to call allowance after %d retries: %v", maxRetries, err)
	}

	// 解析结果
	allowance := new(big.Int).SetBytes(result)
	l.Infof("当前 allowance: %s", allowance.String())
	return allowance, nil
}

// WaitForTransactionReceipt 等待交易确认
func (l *TransactionLogic) WaitForTransactionReceipt(client *ethclient.Client, txHash common.Hash, timeout time.Duration) (*evmTypes.Receipt, error) {
	ctx, cancel := context.WithTimeout(l.ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			receipt, err := client.TransactionReceipt(ctx, txHash)
			if err != nil {
				if err == ethereum.NotFound {
					l.Infof("交易尚未确认，继续等待...")
					continue
				}
				return nil, err
			}
			return receipt, nil
		}
	}
}

// EstimateNativeTransferGas 估算原生代币转账的 gas
func (l *TransactionLogic) EstimateNativeTransferGas(client *ethclient.Client, fromAddress, toAddress common.Address, value *big.Int) (uint64, *big.Int, error) {
	// 获取 gas price
	gasPrice, err := client.SuggestGasPrice(l.ctx)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to get gas price: %v", err)
	}

	// 使用 EstimateGas 进行精确估算
	gasLimit, err := client.EstimateGas(l.ctx, ethereum.CallMsg{
		From:  fromAddress,
		To:    &toAddress,
		Value: value,
	})

	if err != nil {
		l.Infof("Gas 估算失败，使用默认值: %v", err)
		gasLimit = 21000
	}

	// 确保 gas limit 不低于最小值
	if gasLimit < 21000 {
		gasLimit = 21000
	}

	// 增加缓冲
	gasLimit = gasLimit * 110 / 100

	return gasLimit, gasPrice, nil
}

// EstimateERC20TransferGas 估算 ERC20 转账的 gas
func (l *TransactionLogic) EstimateERC20TransferGas(client *ethclient.Client, fromAddress, tokenAddress common.Address, data []byte) (uint64, *big.Int, error) {
	// 获取 gas price
	gasPrice, err := client.SuggestGasPrice(l.ctx)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to get gas price: %v", err)
	}

	gasLimit, err := client.EstimateGas(l.ctx, ethereum.CallMsg{
		From: fromAddress,
		To:   &tokenAddress,
		Data: data,
	})

	if err != nil {
		l.Infof("ERC20 Gas 估算失败，使用默认值: %v", err)
		gasLimit = 100000
	}

	// 确保 gas limit 不低于最小值
	if gasLimit < 60000 {
		gasLimit = 60000
	}

	// 增加缓冲
	gasLimit = gasLimit * 120 / 100

	return gasLimit, gasPrice, nil
}

// GetChainDisplayName 获取链的显示名称
func (l *TransactionLogic) GetChainDisplayName(chain string) string {
	switch chain {
	case "BSC":
		return "BSC 主网"
	case "BSC-TestNet":
		return "BSC 测试网"
	case "ETH":
		return "以太坊主网"
	case "ETH-Goerli":
		return "以太坊 Goerli 测试网"
	case "ETH-Sepolia":
		return "以太坊 Sepolia 测试网"
	case "Polygon":
		return "Polygon 主网"
	case "Polygon-Mumbai":
		return "Polygon Mumbai 测试网"
	case "Solana", "SOL", "solana", "sol":
		return "Solana 主网"
	default:
		return chain
	}
}

// SendTransactionWithRetry 带重试机制的交易发送
func (l *TransactionLogic) SendTransactionWithRetry(client *ethclient.Client, privateKey *ecdsa.PrivateKey, to common.Address, value *big.Int, data []byte, gasLimit uint64, gasPrice *big.Int, chainId int64, maxRetries int) (string, error) {
	for i := 0; i < maxRetries; i++ {
		txHash, err := l.BuildAndSendTransaction(client, privateKey, to, value, data, gasLimit, gasPrice, chainId)
		if err != nil {
			l.Errorf("发送交易失败 (尝试 %d/%d): %v", i+1, maxRetries, err)
			if i == maxRetries-1 {
				return "", err
			}
			time.Sleep(time.Duration(i+1) * time.Second)
			continue
		}
		return txHash, nil
	}
	return "", errors.New("发送交易最终失败")
}
