package constant

type Chain string

const (
	ChainEVM    Chain = "EVM"    // EVM兼容链统一标识（ETH、BSC、Polygon等）
	ChainBTC    Chain = "BTC"    // Bitcoin链
	ChainSOLANA Chain = "SOLANA" // Solana链

	// 以下为具体链标识，用于交易处理时区分网络
	ChainETH Chain = "ETH"
	ChainBSC Chain = "BSC"
)

// SupportedChains lists all chains that are currently supported for wallet initialization.
var SupportedChains = []Chain{
	ChainEVM,    // EVM兼容链
	ChainBTC,    // Bitcoin链
	ChainSOLANA, // Solana链

	// 具体链标识，用于交易处理
	ChainETH,
	ChainBSC,
}

// WalletInitConfig 定义钱包初始化时需要创建的链类型
// 在实际部署时，可以通过配置文件或环境变量来控制
var WalletInitConfig = struct {
	EnabledChains []Chain `json:"enabled_chains"`
}{
	EnabledChains: []Chain{
		ChainEVM,    // EVM兼容链（统一生成，支持ETH、BSC、Polygon等）
		ChainSOLANA, // Solana链
		ChainBTC,    // Bitcoin链
	},
}

// IsChainSupported checks if a given chain is in the list of supported chains.
func IsChainSupported(chain string) bool {
	for _, supportedChain := range SupportedChains {
		if string(supportedChain) == chain {
			return true
		}
	}
	return false
}
