package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	BaseURL = "http://localhost:8888/api"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "create":
		handleCreateWallet(os.Args[2:])
	case "send":
		handleSendTransaction(os.Args[2:])
	case "swap":
		handleSwapTransaction(os.Args[2:])
	case "bridge":
		handleBridgeTransaction(os.Args[2:])
	case "approve":
		handleApproveToken(os.Args[2:])
	case "revoke":
		handleRevokeApproval(os.Args[2:])
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Printf("未知命令: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("MPC Demo CLI - 多链钱包管理工具")
	fmt.Println("")
	fmt.Println("用法:")
	fmt.Println("  cli <command> [options]")
	fmt.Println("")
	fmt.Println("支持的命令:")
	fmt.Println("  create   - 创建钱包")
	fmt.Println("  send     - 发送/转账")
	fmt.Println("  swap     - 代币交换")
	fmt.Println("  bridge   - 跨链转账")
	fmt.Println("  approve  - 授权代币")
	fmt.Println("  revoke   - 取消授权")
	fmt.Println("  help     - 显示帮助信息")
	fmt.Println("")
	fmt.Println("示例:")
	fmt.Println("  cli create --name MyWallet --phone 18888888888")
	fmt.Println("  cli send --from 0x123... --to 0x456... --chain BSC --token BNB --amount 0.1")
	fmt.Println("  cli swap --from 0x123... --chain BSC --from-token BNB --to-token USDT --amount 0.1")
	fmt.Println("  cli bridge --from 0x123... --to 0x456... --from-chain BSC --to-chain ETH --amount 100")
	fmt.Println("  cli approve --owner 0x123... --spender 0x456... --token 0x789... --chain BSC")
	fmt.Println("  cli revoke --owner 0x123... --spender 0x456... --token 0x789... --chain BSC")
}

// 处理创建钱包命令
func handleCreateWallet(args []string) {
	fs := flag.NewFlagSet("create", flag.ExitOnError)
	name := fs.String("name", "My-CLI-Wallet", "钱包名称")
	phone := fs.String("phone", "18888888888", "手机号")
	email := fs.String("email", "test@example.com", "邮箱地址")

	fs.Parse(args)

	requestData := map[string]interface{}{
		"name":         *name,
		"phone_number": *phone,
		"email":        *email,
	}

	fmt.Println("=== 创建多链钱包 ===")
	makeAPIRequest("POST", "/wallet_init", requestData)
}

// 处理发送交易命令
func handleSendTransaction(args []string) {
	fs := flag.NewFlagSet("send", flag.ExitOnError)
	from := fs.String("from", "", "发送地址 (必填)")
	to := fs.String("to", "", "接收地址 (必填)")
	chain := fs.String("chain", "BSC", "区块链网络")
	token := fs.String("token", "BNB", "代币符号 (BNB/USDT/ETH等)")
	amount := fs.String("amount", "", "转账金额 (必填)")

	fs.Parse(args)

	if *from == "" || *to == "" || *amount == "" {
		fmt.Println("错误: --from, --to, --amount 参数为必填项")
		fmt.Println("\n用法: cli send --from 0x123... --to 0x456... --chain BSC --token BNB --amount 0.1")
		os.Exit(1)
	}

	// 标准化代币地址
	fromToken := normalizeTokenAddress(*token, *chain)
	toToken := fromToken // 发送交易中 from_token 和 to_token 相同

	requestData := map[string]interface{}{
		"from_address": *from,
		"to_address":   *to,
		"chain":        *chain,
		"from_token":   fromToken,
		"to_token":     toToken,
		"amount":       convertAmountToWei(*amount, *token),
	}

	fmt.Printf("=== 发送 %s %s (链: %s) ===\n", *amount, *token, *chain)
	makeAPIRequest("POST", "/transaction/send", requestData)
}

// 处理交换交易命令
func handleSwapTransaction(args []string) {
	fs := flag.NewFlagSet("swap", flag.ExitOnError)
	from := fs.String("from", "", "发送地址 (必填)")
	to := fs.String("to", "", "接收地址 (可选，默认为发送地址)")
	chain := fs.String("chain", "BSC", "区块链网络")
	fromToken := fs.String("from-token", "", "源代币 (必填)")
	toToken := fs.String("to-token", "", "目标代币 (必填)")
	amount := fs.String("amount", "", "交换金额 (必填)")

	fs.Parse(args)

	if *from == "" || *fromToken == "" || *toToken == "" || *amount == "" {
		fmt.Println("错误: --from, --from-token, --to-token, --amount 参数为必填项")
		fmt.Println("\n用法: cli swap --from 0x123... --chain BSC --from-token BNB --to-token USDT --amount 0.1")
		os.Exit(1)
	}

	// 如果没有指定接收地址，使用发送地址
	toAddress := *to
	if toAddress == "" {
		toAddress = *from
	}

	// 标准化代币地址
	fromTokenAddr := normalizeTokenAddress(*fromToken, *chain)
	toTokenAddr := normalizeTokenAddress(*toToken, *chain)

	requestData := map[string]interface{}{
		"from_address": *from,
		"to_address":   toAddress,
		"chain":        *chain,
		"from_token":   fromTokenAddr,
		"to_token":     toTokenAddr,
		"amount":       convertAmountToWei(*amount, *fromToken),
	}

	fmt.Printf("=== 交换 %s %s -> %s (链: %s) ===\n", *amount, *fromToken, *toToken, *chain)
	makeAPIRequest("POST", "/transaction/swap", requestData)
}

// 处理跨链交易命令
func handleBridgeTransaction(args []string) {
	fs := flag.NewFlagSet("bridge", flag.ExitOnError)
	from := fs.String("from", "", "发送地址 (必填)")
	to := fs.String("to", "", "接收地址 (必填)")
	fromChain := fs.String("from-chain", "BSC", "源链")
	toChain := fs.String("to-chain", "ETH", "目标链")
	token := fs.String("token", "USDT", "代币符号")
	amount := fs.String("amount", "", "跨链金额 (必填)")

	fs.Parse(args)

	if *from == "" || *to == "" || *amount == "" {
		fmt.Println("错误: --from, --to, --amount 参数为必填项")
		fmt.Println("\n用法: cli bridge --from 0x123... --to 0x456... --from-chain BSC --to-chain ETH --token USDT --amount 100")
		os.Exit(1)
	}

	// 获取链ID
	fromChainId := getChainId(*fromChain)
	toChainId := getChainId(*toChain)

	// 标准化代币地址
	fromToken := normalizeTokenAddress(*token, *fromChain)
	toToken := normalizeTokenAddress(*token, *toChain)

	requestData := map[string]interface{}{
		"from_address": *from,
		"to_address":   *to,
		"from_chain":   fromChainId,
		"to_chain":     toChainId,
		"from_token":   fromToken,
		"to_token":     toToken,
		"amount":       convertAmountToWei(*amount, *token),
	}

	fmt.Printf("=== 跨链转账 %s %s: %s -> %s ===\n", *amount, *token, *fromChain, *toChain)
	makeAPIRequest("POST", "/transaction/bridge/execute", requestData)
}

// 处理授权代币命令
func handleApproveToken(args []string) {
	fs := flag.NewFlagSet("approve", flag.ExitOnError)
	owner := fs.String("owner", "", "代币持有者地址 (必填)")
	spender := fs.String("spender", "", "被授权地址 (必填)")
	token := fs.String("token", "", "代币地址 (必填)")
	chain := fs.String("chain", "BSC", "区块链网络")
	amount := fs.String("amount", "max", "授权金额 (默认为无限授权)")

	fs.Parse(args)

	if *owner == "" || *spender == "" || *token == "" {
		fmt.Println("错误: --owner, --spender, --token 参数为必填项")
		fmt.Println("\n用法: cli approve --owner 0x123... --spender 0x456... --token 0x789... --chain BSC")
		os.Exit(1)
	}

	requestData := map[string]interface{}{
		"token_address":   *token,
		"spender_address": *spender,
		"owner_address":   *owner,
		"chain":           *chain,
		"amount":          *amount,
	}

	fmt.Printf("=== 授权代币 %s (链: %s) ===\n", *token, *chain)
	makeAPIRequest("POST", "/transaction/approve", requestData)
}

// 处理取消授权命令
func handleRevokeApproval(args []string) {
	fs := flag.NewFlagSet("revoke", flag.ExitOnError)
	owner := fs.String("owner", "", "代币持有者地址 (必填)")
	spender := fs.String("spender", "", "被授权地址 (必填)")
	token := fs.String("token", "", "代币地址 (必填)")
	chain := fs.String("chain", "BSC", "区块链网络")

	fs.Parse(args)

	if *owner == "" || *spender == "" || *token == "" {
		fmt.Println("错误: --owner, --spender, --token 参数为必填项")
		fmt.Println("\n用法: cli revoke --owner 0x123... --spender 0x456... --token 0x789... --chain BSC")
		os.Exit(1)
	}

	requestData := map[string]interface{}{
		"token_address":   *token,
		"spender_address": *spender,
		"owner_address":   *owner,
		"chain":           *chain,
	}

	fmt.Printf("=== 取消授权 %s (链: %s) ===\n", *token, *chain)
	makeAPIRequest("POST", "/transaction/revoke", requestData)
}

// 发送API请求的通用函数
func makeAPIRequest(method, endpoint string, data interface{}) {
	url := BaseURL + endpoint

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("错误: 无法序列化请求数据: %v", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("错误: 无法创建请求: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	fmt.Printf("发送请求到: %s\n", url)
	fmt.Printf("请求数据: %s\n", string(jsonData))
	fmt.Println()

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("错误: 发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("错误: 读取响应失败: %v", err)
	}

	fmt.Println("=== 响应结果 ===")
	fmt.Printf("状态码: %d\n", resp.StatusCode)

	// 尝试格式化JSON输出
	var result interface{}
	if err := json.Unmarshal(body, &result); err == nil {
		prettyJSON, _ := json.MarshalIndent(result, "", "  ")
		fmt.Printf("响应内容:\n%s\n", string(prettyJSON))
	} else {
		fmt.Printf("响应内容: %s\n", string(body))
	}

	if resp.StatusCode >= 400 {
		os.Exit(1)
	}
}

// 标准化代币地址
func normalizeTokenAddress(token, chain string) string {
	token = strings.ToUpper(token)

	switch chain {
	case "BSC":
		switch token {
		case "BNB":
			return "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE" // BSC原生代币
		case "USDT":
			return "0x55d398326f99059fF775485246999027B3197955" // BSC USDT
		case "USDC":
			return "0x8AC76a51cc950d9822D68b83fE1Ad97B32Cd580d" // BSC USDC
		case "BUSD":
			return "0xe9e7CEA3DedcA5984780Bafc599bD69ADd087D56" // BSC BUSD
		}
	case "ETH":
		switch token {
		case "ETH":
			return "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE" // ETH原生代币
		case "USDT":
			return "0xdAC17F958D2ee523a2206206994597C13D831ec7" // ETH USDT
		case "USDC":
			return "0xA0b86a33E6441b8044EB1E5CfE9F4C4e3A3A6e1D" // ETH USDC
		}
	}

	// 如果不是已知的代币符号，假设它已经是地址格式
	return token
}

// 获取链ID
func getChainId(chain string) int {
	switch strings.ToUpper(chain) {
	case "BSC":
		return 56
	case "ETH", "ETHEREUM":
		return 1
	case "POLYGON":
		return 137
	case "ARBITRUM":
		return 42161
	case "OPTIMISM":
		return 10
	case "BASE":
		return 8453
	default:
		return 56 // 默认BSC
	}
}

// 将金额转换为最小单位 (wei)
func convertAmountToWei(amount, token string) string {
	// 这里简化处理，实际应该根据代币精度进行转换
	// 对于演示，直接返回输入的金额
	// 在生产环境中，需要：
	// 1. 查询代币的 decimals
	// 2. 将 human-readable 金额转换为 wei
	return amount
}
