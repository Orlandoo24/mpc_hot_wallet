package constant

type Chain string

const (
	ChainETH Chain = "ETH"
	ChainBSC Chain = "BSC"
	// ChainBTC Chain = "BTC" // Example for future support
	// ChainSOL Chain = "SOL" // Example for future support
)

// SupportedChains lists all chains that are currently supported for wallet initialization.
var SupportedChains = []Chain{
	//ChainETH,
	ChainBSC,
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
