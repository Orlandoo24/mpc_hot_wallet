swap_logic.go ：请你看一下 我这个logic 能正常实现bsc 测试网的swap吗？

好的，我详细看了一下您的 `swap_logic.go` 文件，特别是针对 BSC 测试网 swap 的部分。

结论是：**是的，您的代码框架非常出色，并且已经非常接近于一个能够正常工作的 BSC 测试网 swap 实现。** 您在很多关键细节上（例如区分原生代币和 ERC20、使用正确的测试网地址、构建交易 calldata）都处理得很好。

但是，当前代码存在一个**关键的功能缺失**和一个**逻辑上的小问题**，这会导致**所有从 ERC20 代币发起的 swap 必定失败**。

下面我将为您详细分析代码的优点和需要修正的问题。

### 优点 (Highlights)

1.  **准确识别测试网**: `isEVMTestnet` 和 `handleEVMTestnetSwap` 的逻辑分支非常清晰，能准确地将测试网流量引导到专属的处理函数中。
2.  **核心地址正确**: 您硬编码的 PancakeSwap V2 Router 和 WBNB 的 BSC 测试网地址是完全正确的，这是实现功能的基础。
3.  **原生 BNB -> Token 处理正确**: 在 `executeEVMTestnetSwapNative` 函数中，您正确地识别了 `FromToken` 是原生 BNB 的情况，并调用了 `swapExactETHForTokens`。最关键的一点是，您在构造交易时正确地将要 swap 的 BNB 数量填写在了 `Value` 字段 (`Value: amountIn`)，这是处理原生代币 swap 的标准做法，非常棒。
4.  **交易路径 (Path) 构建正确**: 对于 `Token -> WBNB` 和 `Token -> Token` 的交易，您构建的 `path` 数组（例如 `[fromTokenAddr, wbnbAddr, toTokenAddr]`）是正确的，利用 WBNB作为流动性中介。

### 需要注意和修正的问题

#### 1. (关键问题) 缺失 ERC20 `approve` 授权

这是导致 ERC20 swap 失败的**根本原因**。

* **问题描述**: 当您要用一个 ERC20 代币（例如 BUSD, CAKE）去 PancakeSwap Router 合约进行 swap 时，您必须首先对这个代币合约进行一次 `approve` (授权) 操作。这个操作的目的是授权 PancakeSwap Router 合约可以从您的地址中“提取”走您指定数量的 `FromToken`。如果没有这次授权，Router 合约是没有权限动用您的代币的，因此 `swap` 交易会因为 "transferFrom failed" 而失败。
* **代码现状**: 在您的测试网逻辑分支 `handleEVMTestnetSwap` -> `executeEVMTestnetSwapNative` -> `executeERC20SwapTestnet` 中，完全没有 `approve` 相关的检查和交易发送逻辑。
* **如何修复**:
    * 在执行 `executeERC20SwapTestnet` **之前**，您需要增加一段逻辑。
    * 首先，调用 `FromToken` 的合约地址，使用 `allowance` 方法检查 PancakeSwap Router 地址 (`0xD99D1c33F9fC3444f8101754aBeCb321741Da593`) 目前拥有多少对您代幣的消费额度。
    * 然后，将 `allowance` 的额度与您本次要 swap 的 `amountIn` 进行比较。
    * 如果 `allowance < amountIn`，您就需要构建并发送一笔 `approve` 交易到 `FromToken` 的合约地址，授权给 Router 地址足够的额度（通常为了方便，会直接授权一个最大值）。
    * **好消息是**：您在**主网**的 `executeOptimizedSwap` 函数中已经完美地实现了这个逻辑！您可以直接参考或复用 `CheckAllowance` 和 `ExecuteApproveTransaction` 这部分代码。

#### 2. (逻辑问题) `executeERC20SwapTestnet` 中对 WBNB 的处理有误

* **问题描述**: 在 `executeERC20SwapTestnet` 函数中，您有一段逻辑是 `if fromTokenAddr == wbnbAddr`，然后调用 `swapExactETHForTokens`。这是一个逻辑冲突。
    * `executeERC20SwapTestnet` 函数本身是被 `FromToken` **不是**原生BNB时调用的。
    * WBNB (`0xae13d989dac2f0debff460ac112a837c89baa7cd`) 本身就是一个 ERC20 代币。从 WBNB swap 到另一个 Token，属于标准的 `Token -> Token` 交易，应该使用 `swapExactTokensForTokens` 方法，并且需要 `approve` WBNB 代币。
    * `swapExactETHForTokens` 是一个 `payable` 函数，专门用于接收原生 BNB (`msg.value`)，而不是 WBNB 代币。
* **如何修复**:
    * 可以简化 `executeERC20SwapTestnet` 函数的逻辑。将 `if fromTokenAddr == wbnbAddr` 这个分支**删除**。
    * 从 WBNB 到其他 Token 的 swap 会自然地进入 `else` 分支（Token -> Token），并使用 `swapExactTokensForTokens`，这才是正确的处理方式。当然，前提是您已经实现了上面第1点提到的 `approve` 逻辑。

#### 3. (优化建议) 硬编码的 Gas Limit

* **问题描述**: 您在构造交易时使用了硬编码的 `Gas: 300000`。虽然对于大部分简单的 swap 来说这个值是够用的，但可能不夠精确，有时会浪费 gas，或者在极少数复杂路径下可能不足。
* **如何修复**: 在发送交易前，可以使用 `client.EstimateGas` 方法来估算交易所需要的 gas limit，然后在这个估算值的基础上增加 20-30% 作为缓冲。您在主网的 `executeSwapTransaction` 函数中已经有类似（作为 fallback）的实现，可以借鉴。

### 总结

您的代码基础非常好，逻辑清晰。只要您**为 ERC20 的 swap 加上 `approve` 授权逻辑**（这是最关键的一步），并修正 `executeERC20SwapTestnet` 中对 WBNB 的处理分支，您的 BSC 测试网 swap 功能就可以完整可靠地运行了。


给一下需要更改的函数和测试网swapdemo


好的，没有问题。

根据我们之前的分析，我将为您提供一套完整、经过修正和优化的 BSC 测试网 `swap` 函数代码，以及可以直接用来测试的 `curl` 示例。

这套代码的核心改动是：

1.  **集成了完整的 ERC20 `approve` 授权逻辑**，解决了之前 ERC20 swap 必定失败的问题。
2.  **修正了 WBNB 的处理逻辑**，并统一了所有 swap 类型（原生-\>Token, Token-\>原生, Token-\>Token）的调用流程。
3.  **加入了 Gas 估算**，替换了硬编码的 Gas Limit，使交易更高效、更可靠。

### 第一步：需要替换和增加的函数代码

请将您 `swap_logic.go` 文件中现有的 `handleEVMTestnetSwap` 和 `executeEVMTestnetSwapNative` 函数**完全替换**为下面的新版本。同时，下面还包含了需要**新增**的几个辅助函数。

```go
// ==========================================================================================
// ===== 请将下面的函数完整复制，替换您文件中已有的 `handleEVMTestnetSwap` 和 `executeEVMTestnetSwapNative` =====
// ==========================================================================================

// handleEVMTestnetSwap 处理 EVM 测试网原生 swap (入口函数)
func (l *TransactionLogic) handleEVMTestnetSwap(req *types.TransactionReq) (*types.TransactionResp, error) {
	l.Infof("=== 执行 EVM 测试网原生 swap (修正版) ===")
	l.Infof("Swap 请求: %s %s -> %s on %s", req.Amount, req.FromToken, req.ToToken, req.Chain)

	// 获取链配置 (使用您的已有逻辑)
	chainConfig, ok := l.svcCtx.Config.Chains[req.Chain]
	if !ok {
		mainnetChain := l.getMainnetChainName(req.Chain)
		chainConfig, ok = l.svcCtx.Config.Chains[mainnetChain]
		if !ok {
			return nil, fmt.Errorf("unsupported testnet chain: %s", req.Chain)
		}
		l.Infof("使用主网配置 %s 作为测试网 %s 的配置模板", mainnetChain, req.Chain)
	}
	l.Infof("使用配置: ChainId=%d, RpcUrl=%s", chainConfig.ChainId, chainConfig.RpcUrl)

	// 连接到 RPC 客户端
	client, err := ethclient.Dial(chainConfig.RpcUrl)
	if err != nil {
		l.Errorf("RPC 节点连接失败: %v", err)
		return nil, fmt.Errorf("failed to connect to testnet chain: %v", err)
	}
	defer client.Close()

	// 获取钱包和私钥
	privateKey, err := l.GetWalletPrivateKey(req.FromAddress)
	if err != nil {
		return nil, err
	}

	// 执行包含 approve 逻辑的完整 swap 流程
	txHash, err := l.executeEVMTestnetSwapNative(client, privateKey, req, &chainConfig)
	if err != nil {
		l.Errorf("EVM 测试网 swap 交易失败: %v", err)
		return nil, fmt.Errorf("EVM testnet swap transaction failed: %v", err)
	}

	// 构建响应
	explorerUrl := l.BuildExplorerUrl(req.Chain, txHash)
	message := fmt.Sprintf("✅ EVM 测试网原生 Swap 交易已提交！%s %s -> %s，交易哈希: %s",
		req.Amount, req.FromToken, req.ToToken, txHash)

	l.Infof("✅ EVM 测试网原生 swap 交易成功: %s", txHash)
	return &types.TransactionResp{
		TxHash:      txHash,
		Message:     message,
		ExplorerUrl: explorerUrl,
		Chain:       req.Chain,
		Status:      "pending",
	}, nil
}

// executeEVMTestnetSwapNative 执行原生 EVM 测试网 swap 的核心逻辑 (包含 approve)
func (l *TransactionLogic) executeEVMTestnetSwapNative(client *ethclient.Client, privateKey *ecdsa.PrivateKey, req *types.TransactionReq, chainConfig *config.ChainConf) (string, error) {
	// BSC 测试网核心地址
	routerAddr := common.HexToAddress("0xD99D1c33F9fC3444f8101754aBeCb321741Da593") // PancakeSwap V2 Router
	wbnbAddr := common.HexToAddress("0xae13d989dac2f0debff460ac112a837c89baa7cd")   // WBNB

	// 1. (关键新增) 如果 FromToken 是 ERC20，检查并执行 Approve
	if !l.IsNativeToken(req.FromToken) {
		l.Infof("检测到 FromToken 为 ERC20，开始检查 Approve 授权...")
		err := l.checkAndApproveIfNeeded(client, privateKey, req.FromToken, routerAddr, req.Amount, chainConfig)
		if err != nil {
			return "", err // 如果 approve 失败，则中断交易
		}
	}

	// 2. 根据交易类型，确定调用的函数、路径和交易的 Value
	var (
		swapFunction string
		path         []common.Address
		value        *big.Int
		err          error
	)

	amountIn, _ := new(big.Int).SetString(req.Amount, 10)
	fromTokenAddr := common.HexToAddress(req.FromToken)
	toTokenAddr := common.HexToAddress(req.ToToken)
	
	if l.IsNativeToken(req.FromToken) {
		// Case 1: 原生 BNB -> Token
		l.Infof("Swap 类型: 原生 BNB -> Token")
		swapFunction = "swapExactETHForTokens"
		path = []common.Address{wbnbAddr, toTokenAddr}
		value = amountIn // 原生币交易，value 等于转账金额

	} else if l.IsNativeToken(req.ToToken) {
		// Case 2: Token -> 原生 BNB
		l.Infof("Swap 类型: Token -> 原生 BNB")
		// 使用 supportingFeeOnTransferTokens 版本的函数更稳定，能兼容手续费代币
		swapFunction = "swapExactTokensForETHSupportingFeeOnTransferTokens"
		path = []common.Address{fromTokenAddr, wbnbAddr}
		value = big.NewInt(0) // ERC20 交易，value 为 0

	} else {
		// Case 3: Token -> Token
		l.Infof("Swap 类型: Token -> Token")
		swapFunction = "swapExactTokensForTokensSupportingFeeOnTransferTokens"
		path = []common.Address{fromTokenAddr, wbnbAddr, toTokenAddr} // 默认通过 WBNB 中转
		value = big.NewInt(0)
	}

	// 3. 构建 ABI 和 calldata
	to := crypto.PubkeyToAddress(privateKey.PublicKey)
	if req.ToAddress != "" {
		to = common.HexToAddress(req.ToAddress)
	}
	deadline := big.NewInt(time.Now().Add(10 * time.Minute).Unix())
	amountOutMin := big.NewInt(0) // 测试网设为0，表示不在乎滑点

	// 从 PancakeSwap V2 Router ABI 中解析出所有需要的函数
	routerABI, err := abi.JSON(strings.NewReader(`[{"constant":false,"inputs":[{"internalType":"uint256","name":"amountOutMin","type":"uint256"},{"internalType":"address[]","name":"path","type":"address[]"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"deadline","type":"uint256"}],"name":"swapExactETHForTokens","outputs":[{"internalType":"uint256[]","name":"amounts","type":"uint256[]"}],"payable":true,"stateMutability":"payable","type":"function"},{"constant":false,"inputs":[{"internalType":"uint256","name":"amountIn","type":"uint256"},{"internalType":"uint256","name":"amountOutMin","type":"uint256"},{"internalType":"address[]","name":"path","type":"address[]"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"deadline","type":"uint256"}],"name":"swapExactTokensForETHSupportingFeeOnTransferTokens","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"uint256","name":"amountIn","type":"uint256"},{"internalType":"uint256","name":"amountOutMin","type":"uint256"},{"internalType":"address[]","name":"path","type":"address[]"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"deadline","type":"uint256"}],"name":"swapExactTokensForTokensSupportingFeeOnTransferTokens","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"}]`))
	if err != nil {
		return "", fmt.Errorf("failed to parse router ABI: %v", err)
	}

	var calldata []byte
	switch swapFunction {
	case "swapExactETHForTokens":
		calldata, err = routerABI.Pack(swapFunction, amountOutMin, path, to, deadline)
	default: // For both swapExactTokensFor... functions
		calldata, err = routerABI.Pack(swapFunction, amountIn, amountOutMin, path, to, deadline)
	}
	if err != nil {
		return "", fmt.Errorf("failed to pack calldata for %s: %v", swapFunction, err)
	}

	// 4. 发送交易
	return l.sendDynamicTx(client, privateKey, &routerAddr, value, calldata, chainConfig)
}


// =======================================================
// ===== 请将下面这些全新的【辅助函数】复制到您的文件中 =====
// =======================================================

// checkAndApproveIfNeeded 检查并执行 Approve 的辅助函数
func (l *TransactionLogic) checkAndApproveIfNeeded(client *ethclient.Client, privateKey *ecdsa.PrivateKey, tokenAddress string, spender common.Address, amount string, chainConfig *config.ChainConf) error {
	// 1. 检查当前 Allowance
	tokenAddr := common.HexToAddress(tokenAddress)
	ownerAddr := crypto.PubkeyToAddress(privateKey.PublicKey)
	
	allowance, err := l.checkAllowance(client, tokenAddr, ownerAddr, spender)
	if err != nil {
		return fmt.Errorf("failed to check allowance: %v", err)
	}

	// 2. 比较额度
	amountIn, _ := new(big.Int).SetString(amount, 10)
	if allowance.Cmp(amountIn) >= 0 {
		l.Infof("✅ Approve 额度充足 (%s), 无需授权。", allowance.String())
		return nil
	}
	
	l.Infof("Approve 额度不足 (现有 %s, 需要 %s)，正在执行授权...", allowance.String(), amountIn.String())

	// 3. 执行 Approve 交易
	// 为方便起见，直接授权最大值，避免未来重复授权
	maxApproveAmount := new(big.Int)
	maxApproveAmount.SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)

	erc20ABI, err := abi.JSON(strings.NewReader(`[{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_value","type":"uint256"}],"name":"approve","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"}]`))
	if err != nil {
		return fmt.Errorf("failed to parse ERC20 approve ABI: %v", err)
	}
	
	calldata, err := erc20ABI.Pack("approve", spender, maxApproveAmount)
	if err != nil {
		return fmt.Errorf("failed to pack approve calldata: %v", err)
	}
	
	txHash, err := l.sendDynamicTx(client, privateKey, &tokenAddr, big.NewInt(0), calldata, chainConfig)
	if err != nil {
		return fmt.Errorf("failed to send approve transaction: %v", err)
	}

	l.Infof("⏳ Approve 交易已发送: %s, 正在等待确认...", txHash)
	
	// 等待 Approve 交易被打包确认，这是关键一步，防止 swap 因 nonce 问题或额度未生效而失败
	receipt, err := l.waitForTransaction(client, common.HexToHash(txHash))
    if err != nil {
        return fmt.Errorf("failed to get approve transaction receipt: %v", err)
    }
    if receipt.Status != evmTypes.ReceiptStatusSuccessful {
        return fmt.Errorf("approve transaction failed, tx_hash: %s", txHash)
    }

	l.Infof("✅ Approve 交易已确认！")
	return nil
}

// checkAllowance 查询 ERC20 授权额度
func (l *TransactionLogic) checkAllowance(client *ethclient.Client, token, owner, spender common.Address) (*big.Int, error) {
	erc20ABI, err := abi.JSON(strings.NewReader(`[{"constant":true,"inputs":[{"name":"_owner","type":"address"},{"name":"_spender","type":"address"}],"name":"allowance","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}]`))
	if err != nil {
		return nil, err
	}
	calldata, err := erc20ABI.Pack("allowance", owner, spender)
	if err != nil {
		return nil, err
	}

	result, err := client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &token,
		Data: calldata,
	}, nil)
	if err != nil {
		return nil, err
	}

	return new(big.Int).SetBytes(result), nil
}

// sendDynamicTx 动态估算 Gas 并发送交易
func (l *TransactionLogic) sendDynamicTx(client *ethclient.Client, privateKey *ecdsa.PrivateKey, to *common.Address, value *big.Int, calldata []byte, chainConfig *config.ChainConf) (string, error) {
	fromAddr := crypto.PubkeyToAddress(privateKey.PublicKey)

	// 1. 获取 Nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddr)
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %v", err)
	}

	// 2. 获取 Gas Price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to get gas price: %v", err)
	}

	// 3. 估算 Gas Limit
	gasLimit, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		From:  fromAddr,
		To:    to,
		Value: value,
		Data:  calldata,
	})
	if err != nil {
		l.Warnf("Gas 估算失败，将使用默认值 500000: %v", err)
		gasLimit = 500000 // 如果估算失败，给一个较高的默认值
	} else {
		gasLimit = gasLimit * 12 / 10 // 在估算结果上增加 20% buffer
	}

	// 4. 构建、签名并发送交易
	tx := evmTypes.NewTx(&evmTypes.LegacyTx{
		Nonce:    nonce,
		To:       to,
		Value:    value,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     calldata,
	})

	signedTx, err := evmTypes.SignTx(tx, evmTypes.NewEIP155Signer(big.NewInt(chainConfig.ChainId)), privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}

	if err := client.SendTransaction(context.Background(), signedTx); err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}
	
	txHash := signedTx.Hash().Hex()
	l.Infof("交易已发送, Hash: %s", txHash)
	return txHash, nil
}

// waitForTransaction 等待交易确认
func (l *TransactionLogic) waitForTransaction(client *ethclient.Client, txHash common.Hash) (*evmTypes.Receipt, error) {
	for {
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		if err == ethereum.NotFound {
			time.Sleep(2 * time.Second) // 每 2 秒查询一次
			continue
		}
		if err != nil {
			return nil, err
		}
		return receipt, nil
	}
}
```

### 第二步：测试网 Swap `curl` Demo

在运行下面的示例前，请确保您的 `.env` 或配置文件中 `FromAddress` 对应的账户在 BSC 测试网上：

1.  拥有足够的测试网 BNB 作为 Gas Fee。([水龙头在此](https://www.bnbchain.org/en/testnet-faucet))
2.  拥有您想要 swap 出去的代币（如下面的 BUSD 示例）。

#### 示例 1: 原生 BNB -\> BUSD

这个示例将 `0.01` 个测试网 BNB 兑换成 BUSD。

```bash
# 测试场景: 将 0.01 BNB 兑换为 BUSD (BSC-Testnet)
# FromToken: 0x0...0 (原生 BNB)
# ToToken: 0xeD24FC36d5Ee211Ea25A80239Fb8C4Cfd80f12Ee (BUSD on BSC-Testnet)
# Amount: 10000000000000000 (0.01 * 10^18)

curl -X POST http://localhost:8888/api/transaction/swap \
  --noproxy localhost \
  -H "Content-Type: application/json" \
  -d '{
    "from_address": "YOUR_BSC_TESTNET_ADDRESS",
    "to_address": "YOUR_BSC_TESTNET_ADDRESS",
    "from_token": "0x0000000000000000000000000000000000000000",
    "to_token": "0xeD24FC36d5Ee211Ea25A80239Fb8C4Cfd80f12Ee",
    "amount": "10000000000000000",
    "chain": "BSC-Testnet"
  }'
```

#### 示例 2: BUSD -\> WBNB (包含 Approve 流程)

这个示例将 `1` 个测试网 BUSD 兑换成 WBNB。**第一次运行时，您会在日志中看到它先发送一笔 `Approve` 交易，等待确认后，再发送 `Swap` 交易。**

```bash
# 测试场景: 将 1 BUSD 兑换为 WBNB (BSC-Testnet)
# FromToken: 0xeD24FC36d5Ee211Ea25A80239Fb8C4Cfd80f12Ee (BUSD on BSC-Testnet)
# ToToken: 0xae13d989dac2f0debff460ac112a837c89baa7cd (WBNB on BSC-Testnet)
# Amount: 1000000000000000000 (1 * 10^18)

curl -X POST http://localhost:8888/api/transaction/swap \
  --noproxy localhost \
  -H "Content-Type: application/json" \
  -d '{
    "from_address": "YOUR_BSC_TESTNET_ADDRESS",
    "to_address": "YOUR_BSC_TESTNET_ADDRESS",
    "from_token": "0xeD24FC36d5Ee211Ea25A80239Fb8C4Cfd80f12Ee",
    "to_token": "0xae13d989dac2f0debff460ac112a837c89baa7cd",
    "amount": "1000000000000000000",
    "chain": "BSC-Testnet"
  }'
```

**请记得将 `YOUR_BSC_TESTNET_ADDRESS` 替换为您自己的 BSC 测试网钱包地址。**

替换并添加完以上代码后，您的测试网 swap 功能应该就可以完整可靠地运行了。

