curl https://li.quest/v1/chains

```
{
  "chains": [
    {
      "key": "eth",
      "chainType": "EVM",
      "name": "Ethereum",
      "coin": "ETH",
      "id": 1,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/ethereum.svg",
      "tokenlistUrl": "https://gateway.ipfs.io/ipns/tokens.uniswap.org",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "relayerSupported": false,
      "metamask": {
        "chainId": "0x1",
        "blockExplorerUrls": [
          "https://etherscan.io/"
        ],
        "chainName": "Ethereum Mainnet",
        "nativeCurrency": {
          "name": "ETH",
          "symbol": "ETH",
          "decimals": 18
        },
        "rpcUrls": [
          "https://ethereum-rpc.publicnode.com",
          "https://eth.drpc.org"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 1,
        "symbol": "ETH",
        "decimals": 18,
        "name": "ETH",
        "coinKey": "ETH",
        "logoURI": "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2/logo.png",
        "priceUSD": "4304.895592234823"
      },
      "diamondAddress": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
      "permit2": "0x000000000022D473030F116dDEE9F6B43aC78BA3",
      "permit2Proxy": "0x89c6340B1a1f4b25D36cd8B063D49045caF3f818"
    },
    {
      "key": "bsc",
      "chainType": "EVM",
      "name": "BSC",
      "coin": "BNB",
      "id": 56,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/bsc.svg",
      "tokenlistUrl": "https://tokens.pancakeswap.finance/pancakeswap-extended.json",
      "faucetUrls": [
        "https://stakely.io/faucet/bsc-chain-bnb"
      ],
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "relayerSupported": false,
      "metamask": {
        "chainId": "0x38",
        "blockExplorerUrls": [
          "https://bscscan.com/"
        ],
        "chainName": "BNB Smart Chain Mainnet",
        "nativeCurrency": {
          "name": "BNB",
          "symbol": "BNB",
          "decimals": 18
        },
        "rpcUrls": [
          "https://bsc-dataseed.binance.org",
          "https://bsc-dataseed.bnbchain.org",
          "https://bsc-rpc.publicnode.com",
          "https://bsc-dataseed1.defibit.io",
          "https://bsc-dataseed1.ninicoin.io"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 56,
        "symbol": "BNB",
        "decimals": 18,
        "name": "BNB",
        "coinKey": "BNB",
        "logoURI": "https://assets.coingecko.com/coins/images/825/small/binance-coin-logo.png?1547034615",
        "priceUSD": "854.04553882"
      },
      "diamondAddress": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
      "permit2": "0x000000000022D473030F116dDEE9F6B43aC78BA3",
      "permit2Proxy": "0x89c6340B1a1f4b25D36cd8B063D49045caF3f818"
    },
    {
      "key": "arb",
      "chainType": "EVM",
      "name": "Arbitrum",
      "coin": "ETH",
      "id": 42161,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/arbitrum.svg",
      "tokenlistUrl": "https://raw.githubusercontent.com/sushiswap/default-token-list/master/tokens/arbitrum.json",
      "faucetUrls": [
        "https://bridge.arbitrum.io/"
      ],
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "relayerSupported": true,
      "metamask": {
        "chainId": "0xa4b1",
        "blockExplorerUrls": [
          "https://arbiscan.io/"
        ],
        "chainName": "Arbitrum One",
        "nativeCurrency": {
          "name": "ETH",
          "symbol": "ETH",
          "decimals": 18
        },
        "rpcUrls": [
          "https://arb1.arbitrum.io/rpc",
          "https://arbitrum-one-rpc.publicnode.com",
          "https://arbitrum.drpc.org"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 42161,
        "symbol": "ETH",
        "decimals": 18,
        "name": "ETH",
        "coinKey": "ETH",
        "logoURI": "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2/logo.png",
        "priceUSD": "4304.895592234823"
      },
      "diamondAddress": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
      "permit2": "0x000000000022D473030F116dDEE9F6B43aC78BA3",
      "permit2Proxy": "0x89c6340B1a1f4b25D36cd8B063D49045caF3f818"
    },
    {
      "key": "bas",
      "chainType": "EVM",
      "name": "Base",
      "coin": "ETH",
      "id": 8453,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/base.svg",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "relayerSupported": true,
      "metamask": {
        "chainId": "0x2105",
        "blockExplorerUrls": [
          "https://basescan.org/"
        ],
        "chainName": "Base",
        "nativeCurrency": {
          "name": "Ethereum",
          "symbol": "ETH",
          "decimals": 18
        },
        "rpcUrls": [
          "https://mainnet.base.org",
          "https://base-rpc.publicnode.com",
          "https://base.drpc.org"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 8453,
        "symbol": "ETH",
        "decimals": 18,
        "name": "ETH",
        "coinKey": "ETH",
        "logoURI": "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2/logo.png",
        "priceUSD": "4304.895592234823"
      },
      "diamondAddress": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
      "permit2": "0x000000000022D473030F116dDEE9F6B43aC78BA3",
      "permit2Proxy": "0x89c6340B1a1f4b25D36cd8B063D49045caF3f818"
    },
    {
      "key": "bls",
      "chainType": "EVM",
      "name": "Blast",
      "coin": "ETH",
      "id": 81457,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/blast.svg",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "relayerSupported": true,
      "metamask": {
        "chainId": "0x13e31",
        "blockExplorerUrls": [
          "https://blastscan.io/",
          "https://blastexplorer.io/"
        ],
        "chainName": "Blast",
        "nativeCurrency": {
          "name": "Ethereum",
          "symbol": "ETH",
          "decimals": 18
        },
        "rpcUrls": [
          "https://rpc.blast.io",
          "https://blast-rpc.publicnode.com",
          "https://blast.drpc.org"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 81457,
        "symbol": "ETH",
        "decimals": 18,
        "name": "ETH",
        "coinKey": "ETH",
        "logoURI": "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2/logo.png",
        "priceUSD": "4305.92"
      },
      "diamondAddress": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
      "permit2": "0x000000000022D473030F116dDEE9F6B43aC78BA3",
      "permit2Proxy": "0x89c6340B1a1f4b25D36cd8B063D49045caF3f818"
    },
    {
      "key": "ava",
      "chainType": "EVM",
      "name": "Avalanche",
      "coin": "AVAX",
      "id": 43114,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/avalanche.svg",
      "tokenlistUrl": "https://raw.githubusercontent.com/sushiswap/default-token-list/master/tokens/avalanche.json",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "relayerSupported": true,
      "metamask": {
        "chainId": "0xa86a",
        "blockExplorerUrls": [
          "https://snowtrace.io/"
        ],
        "chainName": "Avalanche C-Chain",
        "nativeCurrency": {
          "name": "AVAX",
          "symbol": "AVAX",
          "decimals": 18
        },
        "rpcUrls": [
          "https://api.avax.network/ext/bc/C/rpc",
          "https://avalanche-c-chain-rpc.publicnode.com",
          "https://avalanche.drpc.org"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 43114,
        "symbol": "AVAX",
        "decimals": 18,
        "name": "AVAX",
        "coinKey": "AVAX",
        "logoURI": "https://static.debank.com/image/avax_token/logo_url/avax/0b9c84359c84d6bdd5bfda9c2d4c4a82.png",
        "priceUSD": "24.25"
      },
      "diamondAddress": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
      "permit2": "0x000000000022D473030F116dDEE9F6B43aC78BA3",
      "permit2Proxy": "0x89c6340B1a1f4b25D36cd8B063D49045caF3f818"
    },
    {
      "key": "pol",
      "chainType": "EVM",
      "name": "Polygon",
      "coin": "POL",
      "id": 137,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/polygon.svg",
      "tokenlistUrl": "https://unpkg.com/quickswap-default-token-list@1.0.71/build/quickswap-default.tokenlist.json",
      "faucetUrls": [
        "https://stakely.io/faucet/polygon-matic"
      ],
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "relayerSupported": true,
      "metamask": {
        "chainId": "0x89",
        "blockExplorerUrls": [
          "https://polygonscan.com/"
        ],
        "chainName": "Polygon Mainnet",
        "nativeCurrency": {
          "name": "Polygon Ecosystem Token",
          "symbol": "POL",
          "decimals": 18
        },
        "rpcUrls": [
          "https://polygon-bor-rpc.publicnode.com",
          "https://polygon.drpc.org",
          "https://polygon-rpc.com"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 137,
        "symbol": "POL",
        "decimals": 18,
        "name": "Polygon Ecosystem Token",
        "coinKey": "POL",
        "logoURI": "https://static.debank.com/image/matic_token/logo_url/matic/6f5a6b6f0732a7a235131bd7804d357c.png",
        "priceUSD": "0.2804382"
      },
      "diamondAddress": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
      "permit2": "0x000000000022D473030F116dDEE9F6B43aC78BA3",
      "permit2Proxy": "0x89c6340B1a1f4b25D36cd8B063D49045caF3f818"
    },
    {
      "key": "scl",
      "chainType": "EVM",
      "name": "Scroll",
      "coin": "ETH",
      "id": 534352,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/scroll.svg",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "metamask": {
        "chainId": "0x82750",
        "blockExplorerUrls": [
          "https://scrollscan.com/"
        ],
        "chainName": "Scroll",
        "nativeCurrency": {
          "name": "ETH",
          "symbol": "ETH",
          "decimals": 18
        },
        "rpcUrls": [
          "https://rpc.scroll.io",
          "https://scroll.drpc.org",
          "https://scroll-mainnet.public.blastapi.io",
          "https://1rpc.io/scroll"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 534352,
        "symbol": "ETH",
        "decimals": 18,
        "name": "ETH",
        "coinKey": "ETH",
        "logoURI": "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2/logo.png",
        "priceUSD": "4305.92"
      },
      "diamondAddress": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
      "permit2Proxy": "0x89c6340B1a1f4b25D36cd8B063D49045caF3f818"
    },
    {
      "key": "opt",
      "chainType": "EVM",
      "name": "Optimism",
      "coin": "ETH",
      "id": 10,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/optimism.svg",
      "tokenlistUrl": "https://static.optimism.io/optimism.tokenlist.json",
      "faucetUrls": [
        "https://gateway.optimism.io/"
      ],
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "relayerSupported": true,
      "metamask": {
        "chainId": "0xa",
        "blockExplorerUrls": [
          "https://optimistic.etherscan.io/"
        ],
        "chainName": "OP Mainnet",
        "nativeCurrency": {
          "name": "ETH",
          "symbol": "ETH",
          "decimals": 18
        },
        "rpcUrls": [
          "https://mainnet.optimism.io/",
          "https://optimism-rpc.publicnode.com",
          "https://optimism.drpc.org"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 10,
        "symbol": "ETH",
        "decimals": 18,
        "name": "ETH",
        "coinKey": "ETH",
        "logoURI": "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2/logo.png",
        "priceUSD": "4304.895592234823"
      },
      "diamondAddress": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
      "permit2": "0x000000000022D473030F116dDEE9F6B43aC78BA3",
      "permit2Proxy": "0x89c6340B1a1f4b25D36cd8B063D49045caF3f818"
    },
    {
      "key": "lna",
      "chainType": "EVM",
      "name": "Linea",
      "coin": "ETH",
      "id": 59144,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/linea.svg",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "metamask": {
        "chainId": "0xe708",
        "blockExplorerUrls": [
          "https://lineascan.build/"
        ],
        "chainName": "Linea",
        "nativeCurrency": {
          "name": "ETH",
          "symbol": "ETH",
          "decimals": 18
        },
        "rpcUrls": [
          "https://rpc.linea.build"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 59144,
        "symbol": "ETH",
        "decimals": 18,
        "name": "ETH",
        "coinKey": "ETH",
        "logoURI": "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2/logo.png",
        "priceUSD": "4305.92"
      },
      "diamondAddress": "0xde1e598b81620773454588b85d6b5d4eec32573e",
      "permit2Proxy": "0xE5A89411c7Ef1502123d84ec1F9Da9D8574F4293"
    },
    {
      "key": "era",
      "chainType": "EVM",
      "name": "zkSync",
      "coin": "ETH",
      "id": 324,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/zksync.svg",
      "multicallAddress": "0xF9cda624FBC7e059355ce98a31693d299FACd963",
      "relayerSupported": true,
      "metamask": {
        "chainId": "0x144",
        "blockExplorerUrls": [
          "https://explorer.zksync.io/"
        ],
        "chainName": "zkSync Mainnet",
        "nativeCurrency": {
          "name": "ETH",
          "symbol": "ETH",
          "decimals": 18
        },
        "rpcUrls": [
          "https://mainnet.era.zksync.io"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 324,
        "symbol": "ETH",
        "decimals": 18,
        "name": "ETH",
        "coinKey": "ETH",
        "logoURI": "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2/logo.png",
        "priceUSD": "4305.92"
      },
      "diamondAddress": "0x341e94069f53234fe6dabef707ad424830525715",
      "permit2": "0x0000000000225e31D15943971F47aD3022F714Fa",
      "permit2Proxy": "0x2ec2c4C3dc212C990d1bc2b48b0392A3951D926E"
    },
    {
      "key": "pze",
      "chainType": "EVM",
      "name": "Polygon zkEVM",
      "coin": "ETH",
      "id": 1101,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/zkevm.svg",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "metamask": {
        "chainId": "0x44d",
        "blockExplorerUrls": [
          "https://zkevm.polygonscan.com/"
        ],
        "chainName": "Polygon zkEVM",
        "nativeCurrency": {
          "name": "ETH",
          "symbol": "ETH",
          "decimals": 18
        },
        "rpcUrls": [
          "https://zkevm-rpc.com",
          "https://polygon-zkevm.drpc.org"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 1101,
        "symbol": "ETH",
        "decimals": 18,
        "name": "ETH",
        "coinKey": "ETH",
        "logoURI": "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2/logo.png",
        "priceUSD": "4305.92"
      },
      "diamondAddress": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
      "permit2Proxy": "0x89c6340B1a1f4b25D36cd8B063D49045caF3f818"
    },
    {
      "key": "dai",
      "chainType": "EVM",
      "name": "Gnosis",
      "coin": "DAI",
      "id": 100,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/gnosis.svg",
      "tokenlistUrl": "https://tokens.honeyswap.org/",
      "faucetUrls": [
        "https://stakely.io/faucet/xdai-chain"
      ],
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "metamask": {
        "chainId": "0x64",
        "blockExplorerUrls": [
          "https://gnosis.blockscout.com/"
        ],
        "chainName": "Gnosis",
        "nativeCurrency": {
          "name": "XDAI",
          "symbol": "XDAI",
          "decimals": 18
        },
        "rpcUrls": [
          "https://rpc.gnosischain.com",
          "https://gnosis-rpc.publicnode.com",
          "https://gnosis.drpc.org"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 100,
        "symbol": "xDAI",
        "decimals": 18,
        "name": "xDAI Native Token",
        "coinKey": "DAI",
        "logoURI": "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0x6B175474E89094C44Da98b954EedeAC495271d0F/logo.png",
        "priceUSD": "0.99997053"
      },
      "diamondAddress": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
      "permit2Proxy": "0x89c6340B1a1f4b25D36cd8B063D49045caF3f818"
    },
    {
      "key": "ftm",
      "chainType": "EVM",
      "name": "Fantom",
      "coin": "FTM",
      "id": 250,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/fantom.svg",
      "tokenlistUrl": "https://raw.githubusercontent.com/SpookySwap/spooky-info/master/src/constants/token/spookyswap.json",
      "faucetUrls": [
        "https://stakely.io/faucet/fantom-blockchain-ftm",
        "https://docs.spookyswap.finance/getting-started/how-to-get-fantom-gas"
      ],
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "metamask": {
        "chainId": "0xfa",
        "blockExplorerUrls": [
          "https://ftmscan.com/"
        ],
        "chainName": "Fantom Opera",
        "nativeCurrency": {
          "name": "FTM",
          "symbol": "FTM",
          "decimals": 18
        },
        "rpcUrls": [
          "https://rpcapi.fantom.network",
          "https://rpc.fantom.network",
          "https://fantom-rpc.publicnode.com",
          "https://fantom.drpc.org"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 250,
        "symbol": "FTM",
        "decimals": 18,
        "name": "FTM",
        "coinKey": "FTM",
        "logoURI": "https://static.debank.com/image/ftm_token/logo_url/ftm/33fdb9c5067e94f3a1b9e78f6fa86984.png",
        "priceUSD": "0.30115322"
      },
      "diamondAddress": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
      "permit2Proxy": "0x89c6340B1a1f4b25D36cd8B063D49045caF3f818"
    },
    {
      "key": "mor",
      "chainType": "EVM",
      "name": "Moonriver",
      "coin": "MOVR",
      "id": 1285,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/moonriver.svg",
      "tokenlistUrl": "https://raw.githubusercontent.com/sushiswap/default-token-list/master/tokens/moonriver.json",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "metamask": {
        "chainId": "0x505",
        "blockExplorerUrls": [
          "https://moonriver.moonscan.io/"
        ],
        "chainName": "Moonriver",
        "nativeCurrency": {
          "name": "Moonriver",
          "symbol": "MOVR",
          "decimals": 18
        },
        "rpcUrls": [
          "https://rpc.api.moonriver.moonbeam.network",
          "https://moonriver-rpc.publicnode.com",
          "https://moonriver.drpc.org"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 1285,
        "symbol": "MOVR",
        "decimals": 18,
        "name": "MOVR",
        "coinKey": "MOVR",
        "logoURI": "https://assets.coingecko.com/coins/images/17984/small/9285.png",
        "priceUSD": "5.747"
      },
      "diamondAddress": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
      "permit2Proxy": "0x89c6340B1a1f4b25D36cd8B063D49045caF3f818"
    },
    {
      "key": "moo",
      "chainType": "EVM",
      "name": "Moonbeam",
      "coin": "GLMR",
      "id": 1284,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/moonbeam.svg",
      "tokenlistUrl": "https://raw.githubusercontent.com/BeamSwap/beamswap-tokenlist/main/tokenlist.json",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "metamask": {
        "chainId": "0x504",
        "blockExplorerUrls": [
          "https://moonbeam.moonscan.io/"
        ],
        "chainName": "Moonbeam",
        "nativeCurrency": {
          "name": "GLMR",
          "symbol": "GLMR",
          "decimals": 18
        },
        "rpcUrls": [
          "https://rpc.api.moonbeam.network",
          "https://moonbeam-rpc.publicnode.com",
          "https://moonbeam.drpc.org"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 1284,
        "symbol": "GLMR",
        "decimals": 18,
        "name": "GLMR",
        "coinKey": "GLMR",
        "logoURI": "https://s2.coinmarketcap.com/static/img/coins/64x64/6836.png",
        "priceUSD": "0.06703538"
      },
      "diamondAddress": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
      "permit2Proxy": "0x89c6340B1a1f4b25D36cd8B063D49045caF3f818"
    },
    {
      "key": "fus",
      "chainType": "EVM",
      "name": "FUSE",
      "coin": "FUSE",
      "id": 122,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/fuse.svg",
      "tokenlistUrl": "https://raw.githubusercontent.com/sushiswap/default-token-list/master/tokens/fuse.json",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "metamask": {
        "chainId": "0x7a",
        "blockExplorerUrls": [
          "https://explorer.fuse.io/"
        ],
        "chainName": "Fuse Mainnet",
        "nativeCurrency": {
          "name": "FUSE",
          "symbol": "FUSE",
          "decimals": 18
        },
        "rpcUrls": [
          "https://rpc.fuse.io",
          "https://fuse.drpc.org"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 122,
        "symbol": "FUSE",
        "decimals": 18,
        "name": "FUSE",
        "coinKey": "FUSE",
        "logoURI": "https://s2.coinmarketcap.com/static/img/coins/64x64/5634.png",
        "priceUSD": "0.010725"
      },
      "diamondAddress": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
      "permit2Proxy": "0x89c6340B1a1f4b25D36cd8B063D49045caF3f818"
    },
    {
      "key": "bob",
      "chainType": "EVM",
      "name": "Boba",
      "coin": "ETH",
      "id": 288,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/boba.svg",
      "tokenlistUrl": "https://raw.githubusercontent.com/OolongSwap/boba-community-token-list/main/build/boba.tokenlist.json",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "metamask": {
        "chainId": "0x120",
        "blockExplorerUrls": [
          "https://bobascan.com/"
        ],
        "chainName": "Boba Network",
        "nativeCurrency": {
          "name": "ETH",
          "symbol": "ETH",
          "decimals": 18
        },
        "rpcUrls": [
          "https://mainnet.boba.network",
          "https://replica.boba.network"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 288,
        "symbol": "ETH",
        "decimals": 18,
        "name": "ETH",
        "coinKey": "ETH",
        "logoURI": "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2/logo.png",
        "priceUSD": "4304.69"
      },
      "diamondAddress": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
      "permit2Proxy": "0x89c6340B1a1f4b25D36cd8B063D49045caF3f818"
    },
    {
      "key": "mod",
      "chainType": "EVM",
      "name": "Mode",
      "coin": "ETH",
      "id": 34443,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/mode.svg",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "metamask": {
        "chainId": "0x868b",
        "blockExplorerUrls": [
          "https://explorer.mode.network/",
          "https://modescan.io/"
        ],
        "chainName": "Mode",
        "nativeCurrency": {
          "name": "ETH",
          "symbol": "ETH",
          "decimals": 18
        },
        "rpcUrls": [
          "https://mainnet.mode.network",
          "https://mode.drpc.org"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 34443,
        "symbol": "ETH",
        "decimals": 18,
        "name": "ETH",
        "coinKey": "ETH",
        "logoURI": "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2/logo.png",
        "priceUSD": "4305.92"
      },
      "diamondAddress": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
      "permit2Proxy": "0xBDfF0c1C8B0b779581C4aC3bA1F29667C366C56e"
    },
    {
      "key": "mam",
      "chainType": "EVM",
      "name": "Metis",
      "coin": "METIS",
      "id": 1088,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/metis.svg",
      "tokenlistUrl": "https://raw.githubusercontent.com/digitalnativeinc/default-token-list/dnf/tokens/metis.json",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "metamask": {
        "chainId": "0x440",
        "blockExplorerUrls": [
          "https://andromeda-explorer.metis.io/"
        ],
        "chainName": "Metis Andromeda Mainnet",
        "nativeCurrency": {
          "name": "METIS",
          "symbol": "METIS",
          "decimals": 18
        },
        "rpcUrls": [
          "https://andromeda.metis.io",
          "https://metis.drpc.org"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 1088,
        "symbol": "METIS",
        "decimals": 18,
        "name": "METIS",
        "coinKey": "METIS",
        "logoURI": "https://s2.coinmarketcap.com/static/img/coins/64x64/9640.png",
        "priceUSD": "15.33"
      },
      "diamondAddress": "0x24ca98fB6972F5eE05f0dB00595c7f68D9FaFd68",
      "permit2Proxy": "0x82D7dD86603A6344b82a3914Fe02fEF7b6CDe2E2"
    },
    {
      "key": "lsk",
      "chainType": "EVM",
      "name": "Lisk",
      "coin": "ETH",
      "id": 1135,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/lisk.svg",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "metamask": {
        "chainId": "0x46f",
        "blockExplorerUrls": [
          "https://blockscout.lisk.com/"
        ],
        "chainName": "Lisk",
        "nativeCurrency": {
          "name": "Ethereum",
          "symbol": "ETH",
          "decimals": 18
        },
        "rpcUrls": [
          "https://rpc.api.lisk.com",
          "https://lisk.drpc.org"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 1135,
        "symbol": "ETH",
        "decimals": 18,
        "name": "ETH",
        "coinKey": "ETH",
        "logoURI": "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2/logo.png",
        "priceUSD": "4304.7"
      },
      "diamondAddress": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
      "permit2Proxy": "0x89c6340B1a1f4b25D36cd8B063D49045caF3f818"
    },
    {
      "key": "uni",
      "chainType": "EVM",
      "name": "Unichain",
      "coin": "ETH",
      "id": 130,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/unichain.svg",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "relayerSupported": true,
      "metamask": {
        "chainId": "0x82",
        "blockExplorerUrls": [
          "https://unichain.blockscout.com/",
          "https://uniscan.xyz/"
        ],
        "chainName": "Unichain",
        "nativeCurrency": {
          "name": "Ether",
          "symbol": "ETH",
          "decimals": 18
        },
        "rpcUrls": [
          "https://mainnet.unichain.org"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 130,
        "symbol": "ETH",
        "decimals": 18,
        "name": "ETH",
        "coinKey": "ETH",
        "logoURI": "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2/logo.png",
        "priceUSD": "4305.92"
      },
      "diamondAddress": "0x864b314D4C5a0399368609581d3E8933a63b9232",
      "permit2Proxy": "0x1BCd304fDAd1d1D66529159b1Bc8D13C9158D586"
    },
    {
      "key": "aur",
      "chainType": "EVM",
      "name": "Aurora",
      "coin": "ETH",
      "id": 1313161554,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/aurora.svg",
      "tokenlistUrl": "https://aurora.dev/tokens.json",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "metamask": {
        "chainId": "0x4e454152",
        "blockExplorerUrls": [
          "https://explorer.aurora.dev/"
        ],
        "chainName": "Aurora Mainnet",
        "nativeCurrency": {
          "name": "ETH",
          "symbol": "ETH",
          "decimals": 18
        },
        "rpcUrls": [
          "https://mainnet.aurora.dev"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 1313161554,
        "symbol": "AETH",
        "decimals": 18,
        "name": "AETH",
        "coinKey": "AETH",
        "logoURI": "https://static.debank.com/image/aurora_token/logo_url/aurora/d61441782d4a08a7479d54aea211679e.png",
        "priceUSD": "4295.654344114324"
      },
      "diamondAddress": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
      "permit2Proxy": "0x89c6340B1a1f4b25D36cd8B063D49045caF3f818"
    },
    {
      "key": "sei",
      "chainType": "EVM",
      "name": "Sei",
      "coin": "SEI",
      "id": 1329,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/sei.svg",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "metamask": {
        "chainId": "0x531",
        "blockExplorerUrls": [
          "https://seitrace.com/",
          "https://seistream.app/"
        ],
        "chainName": "Sei Network",
        "nativeCurrency": {
          "name": "SEI",
          "symbol": "SEI",
          "decimals": 18
        },
        "rpcUrls": [
          "https://evm-rpc.sei-apis.com"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 1329,
        "symbol": "SEI",
        "decimals": 18,
        "name": "SEI",
        "coinKey": "SEI",
        "logoURI": "https://cdn.sei.io/sei-app/sei-icon.png",
        "priceUSD": "0.2882"
      },
      "diamondAddress": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
      "permit2Proxy": "0x89c6340B1a1f4b25D36cd8B063D49045caF3f818"
    },
    {
      "key": "hpl",
      "chainType": "EVM",
      "name": "Hyperliquid",
      "coin": "USDC",
      "id": 1337,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/hyperliquid.svg",
      "relayerSupported": false,
      "metamask": {
        "chainId": "0xa4b1",
        "blockExplorerUrls": [
          "https://app.hyperliquid.xyz/explorer/"
        ],
        "chainName": "Hyperliquid",
        "nativeCurrency": {
          "name": "USDC",
          "symbol": "USDC",
          "decimals": 6
        },
        "rpcUrls": [
          "https://li.quest/v1/rpc/1337"
        ]
      },
      "nativeToken": {
        "address": "0xaf88d065e77c8cC2239327C5EDb3A432268e5831",
        "chainId": 1337,
        "symbol": "USDC",
        "decimals": 6,
        "name": "USD Coin",
        "coinKey": "USDC",
        "logoURI": "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48/logo.png",
        "priceUSD": "0.99984316"
      }
    },
    {
      "key": "imx",
      "chainType": "EVM",
      "name": "Immutable zkEVM",
      "coin": "IMX",
      "id": 13371,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/imx.svg",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "metamask": {
        "chainId": "0x343b",
        "blockExplorerUrls": [
          "https://explorer.immutable.com/"
        ],
        "chainName": "Immutable",
        "nativeCurrency": {
          "name": "ImmutableX",
          "symbol": "IMX",
          "decimals": 18
        },
        "rpcUrls": [
          "https://rpc.immutable.com/",
          "https://immutable-zkevm.drpc.org"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 13371,
        "symbol": "IMX",
        "decimals": 18,
        "name": "IMX",
        "coinKey": "IMX",
        "logoURI": "https://static.debank.com/image/eth_token/logo_url/0xf57e7e7c23978c3caec3c3548e3d615c346e79ff/7a875818146ec0508d3e5f5b14f1b4eb.png",
        "priceUSD": "0.5148552"
      },
      "diamondAddress": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
      "permit2Proxy": "0x89c6340B1a1f4b25D36cd8B063D49045caF3f818"
    },
    {
      "key": "flr",
      "chainType": "EVM",
      "name": "Flare",
      "coin": "FLR",
      "id": 14,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/flare.svg",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "relayerSupported": false,
      "metamask": {
        "chainId": "0xe",
        "blockExplorerUrls": [
          "https://flarescan.com/"
        ],
        "chainName": "Flare",
        "nativeCurrency": {
          "name": "Flare",
          "symbol": "FLR",
          "decimals": 18
        },
        "rpcUrls": [
          "https://flare-api.flare.network/ext/C/rpc"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 14,
        "symbol": "FLR",
        "decimals": 18,
        "name": "Flare",
        "coinKey": "FLR",
        "logoURI": "https://static.debank.com/image/flr_token/logo_url/flr/c7d8087092d5d7b80794630612afb32e.png",
        "priceUSD": "0.02061"
      },
      "diamondAddress": "0x198FC70Dfe05E755C81e54bd67Bff3F729344B9b",
      "permit2Proxy": "0xac08C219Fd4856A33a090503BCbc14733d588588"
    },
    {
      "key": "son",
      "chainType": "EVM",
      "name": "Sonic",
      "coin": "S",
      "id": 146,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/sonic.svg",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "metamask": {
        "chainId": "0x92",
        "blockExplorerUrls": [
          "https://explorer.soniclabs.com/",
          "https://sonicscan.org/"
        ],
        "chainName": "Sonic Network",
        "nativeCurrency": {
          "name": "Sonic",
          "symbol": "S",
          "decimals": 18
        },
        "rpcUrls": [
          "https://rpc.soniclabs.com",
          "https://sonic.drpc.org"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 146,
        "symbol": "S",
        "decimals": 18,
        "name": "S",
        "coinKey": "S",
        "logoURI": "https://static.debank.com/image/sonic_token/logo_url/0x039e2fb66102314ce7b64ce5ce3e5183bc94ad38/b4cc70d040518a88adac18d906fcbfff.png",
        "priceUSD": "0.30115322"
      },
      "diamondAddress": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
      "permit2Proxy": "0x89c6340B1a1f4b25D36cd8B063D49045caF3f818"
    },
    {
      "key": "van",
      "chainType": "EVM",
      "name": "Vana",
      "coin": "VAN",
      "id": 1480,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/vana.svg",
      "multicallAddress": "0xD8d2dFca27E8797fd779F8547166A2d3B29d360E",
      "relayerSupported": false,
      "metamask": {
        "chainId": "0x5c8",
        "blockExplorerUrls": [
          "https://vanascan.io/"
        ],
        "chainName": "Vana",
        "nativeCurrency": {
          "name": "VAN",
          "symbol": "VAN",
          "decimals": 18
        },
        "rpcUrls": [
          "https://rpc.vana.org"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 1480,
        "symbol": "VAN",
        "decimals": 18,
        "name": "Vana",
        "coinKey": "VAN",
        "logoURI": "https://static.debank.com/image/vana_token/logo_url/vana/4e286c230ccb1e83967dc85edefbc513.png",
        "priceUSD": "4.066"
      },
      "diamondAddress": "0x198FC70Dfe05E755C81e54bd67Bff3F729344B9b",
      "permit2Proxy": "0xac08C219Fd4856A33a090503BCbc14733d588588"
    },
    {
      "key": "gra",
      "chainType": "EVM",
      "name": "Gravity",
      "coin": "G",
      "id": 1625,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/gravity.svg",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "metamask": {
        "chainId": "0x659",
        "blockExplorerUrls": [
          "https://explorer.gravity.xyz/"
        ],
        "chainName": "Gravity Alpha Mainnet",
        "nativeCurrency": {
          "name": "G",
          "symbol": "G",
          "decimals": 18
        },
        "rpcUrls": [
          "https://rpc.gravity.xyz/"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 1625,
        "symbol": "G",
        "decimals": 18,
        "name": "G",
        "coinKey": "G",
        "logoURI": "https://assets.gravity.xyz/token_logo.png",
        "priceUSD": "0.01107"
      },
      "diamondAddress": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
      "permit2Proxy": "0x89c6340B1a1f4b25D36cd8B063D49045caF3f818"
    },
    {
      "key": "tai",
      "chainType": "EVM",
      "name": "Taiko",
      "coin": "ETH",
      "id": 167000,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/taiko.svg",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "metamask": {
        "chainId": "0x28c58",
        "blockExplorerUrls": [
          "https://taikoscan.io/"
        ],
        "chainName": "Taiko Mainnet",
        "nativeCurrency": {
          "name": "Ethereum",
          "symbol": "ETH",
          "decimals": 18
        },
        "rpcUrls": [
          "https://rpc.mainnet.taiko.xyz",
          "https://rpc.taiko.xyz",
          "https://rpc.taiko.tools"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 167000,
        "symbol": "ETH",
        "decimals": 18,
        "name": "ETH",
        "coinKey": "ETH",
        "logoURI": "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2/logo.png",
        "priceUSD": "4305.92"
      },
      "diamondAddress": "0x3A9A5dBa8FE1C4Da98187cE4755701BCA182f63b",
      "permit2Proxy": "0xE5A89411c7Ef1502123d84ec1F9Da9D8574F4293"
    },
    {
      "key": "soe",
      "chainType": "EVM",
      "name": "Soneium",
      "coin": "ETH",
      "id": 1868,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/soneium.svg",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "metamask": {
        "chainId": "0x74c",
        "blockExplorerUrls": [
          "https://soneium.blockscout.com/"
        ],
        "chainName": "Soneium",
        "nativeCurrency": {
          "name": "Ether",
          "symbol": "ETH",
          "decimals": 18
        },
        "rpcUrls": [
          "https://rpc.soneium.org/"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 1868,
        "symbol": "ETH",
        "decimals": 18,
        "name": "ETH",
        "coinKey": "ETH",
        "logoURI": "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2/logo.png",
        "priceUSD": "4304.69"
      },
      "diamondAddress": "0x864b314D4C5a0399368609581d3E8933a63b9232",
      "permit2Proxy": "0x1BCd304fDAd1d1D66529159b1Bc8D13C9158D586"
    },
    {
      "key": "swl",
      "chainType": "EVM",
      "name": "Swellchain",
      "coin": "ETH",
      "id": 1923,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/swell.svg",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "relayerSupported": false,
      "metamask": {
        "chainId": "0x783",
        "blockExplorerUrls": [
          "https://swellchainscan.io/"
        ],
        "chainName": "Swellchain",
        "nativeCurrency": {
          "name": "Ether",
          "symbol": "ETH",
          "decimals": 18
        },
        "rpcUrls": [
          "https://swell-mainnet.alt.technology"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 1923,
        "symbol": "ETH",
        "decimals": 18,
        "name": "ETH",
        "coinKey": "ETH",
        "logoURI": "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2/logo.png",
        "priceUSD": "4304.69"
      },
      "diamondAddress": "0x76F6937a41910F075024138066708B36139AC104",
      "permit2Proxy": "0xD2494E9C32b643120d5c8457A9cd0505701B8054"
    },
    {
      "key": "opb",
      "chainType": "EVM",
      "name": "opBNB",
      "coin": "BNB",
      "id": 204,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/opbnb.svg",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "metamask": {
        "chainId": "0xcc",
        "blockExplorerUrls": [
          "https://opbnb.bscscan.com/",
          "https://opbnbscan.com/"
        ],
        "chainName": "opBNB Mainnet",
        "nativeCurrency": {
          "name": "BNB",
          "symbol": "BNB",
          "decimals": 18
        },
        "rpcUrls": [
          "https://opbnb-mainnet-rpc.bnbchain.org",
          "https://1rpc.io/opbnb"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 204,
        "symbol": "BNB",
        "decimals": 18,
        "name": "BNB",
        "coinKey": "BNB",
        "logoURI": "https://assets.coingecko.com/coins/images/825/small/binance-coin-logo.png?1547034615",
        "priceUSD": "854.04553882"
      },
      "diamondAddress": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
      "permit2Proxy": "0x89c6340B1a1f4b25D36cd8B063D49045caF3f818"
    },
    {
      "key": "crn",
      "chainType": "EVM",
      "name": "Corn",
      "coin": "BTCN",
      "id": 21000000,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/corn.svg",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "relayerSupported": false,
      "metamask": {
        "chainId": "0x1406f40",
        "blockExplorerUrls": [
          "https://cornscan.io/"
        ],
        "chainName": "Corn",
        "nativeCurrency": {
          "name": "Bitcorn",
          "symbol": "BTCN",
          "decimals": 18
        },
        "rpcUrls": [
          "https://mainnet.corn-rpc.com"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 21000000,
        "symbol": "BTCN",
        "decimals": 18,
        "name": "Bitcorn",
        "coinKey": "BTCN",
        "logoURI": "https://static.debank.com/image/corn_token/logo_url/corn/fe230e468272b84aba78d08bb4140456.png",
        "priceUSD": "110816.99"
      },
      "diamondAddress": "0x11d8E4207a976B8Bb33bD3b494d80f8a6854F06b",
      "permit2Proxy": "0x9C1D1db1CD3EbE4727c2a4208Ca5D65624cC9b4D"
    },
    {
      "key": "lns",
      "chainType": "EVM",
      "name": "Lens",
      "coin": "GHO",
      "id": 232,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/lens.svg",
      "multicallAddress": "0xeee5a340Cdc9c179Db25dea45AcfD5FE8d4d3eB8",
      "relayerSupported": false,
      "metamask": {
        "chainId": "0xe8",
        "blockExplorerUrls": [
          "https://explorer.lens.xyz/"
        ],
        "chainName": "Lens",
        "nativeCurrency": {
          "name": "GHO",
          "symbol": "GHO",
          "decimals": 18
        },
        "rpcUrls": [
          "https://api.lens.matterhosted.dev"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 232,
        "symbol": "GHO",
        "decimals": 18,
        "name": "GHO",
        "coinKey": "GHO",
        "logoURI": "https://s2.coinmarketcap.com/static/img/coins/64x64/23508.png",
        "priceUSD": "0.9509550281762301"
      },
      "diamondAddress": "0xF3B20515d9B193531c48E47c18aF16d1e5d28f9a",
      "permit2Proxy": "0xEC32b79Ce6c4F80160029BDb26CB4a13bB03999A"
    },
    {
      "key": "cro",
      "chainType": "EVM",
      "name": "Cronos",
      "coin": "CRO",
      "id": 25,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/cronos.svg",
      "tokenlistUrl": "https://raw.githubusercontent.com/cronaswap/default-token-list/refs/heads/main/assets/tokens/cronos.json",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "metamask": {
        "chainId": "0x19",
        "blockExplorerUrls": [
          "https://explorer.cronos.org/",
          "https://cronoscan.com/"
        ],
        "chainName": "Cronos Mainnet",
        "nativeCurrency": {
          "name": "Cronos",
          "symbol": "CRO",
          "decimals": 18
        },
        "rpcUrls": [
          "https://evm.cronos.org",
          "https://cronos.drpc.org"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 25,
        "symbol": "CRO",
        "decimals": 18,
        "name": "CRO",
        "coinKey": "CRO",
        "logoURI": "https://s2.coinmarketcap.com/static/img/coins/64x64/3635.png",
        "priceUSD": "0.27311"
      },
      "diamondAddress": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
      "permit2Proxy": "0x89c6340B1a1f4b25D36cd8B063D49045caF3f818"
    },
    {
      "key": "fra",
      "chainType": "EVM",
      "name": "Fraxtal",
      "coin": "FRAX",
      "id": 252,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/fraxtal.svg",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "metamask": {
        "chainId": "0xfc",
        "blockExplorerUrls": [
          "https://fraxscan.com/"
        ],
        "chainName": "Fraxtal",
        "nativeCurrency": {
          "name": "FRAX",
          "symbol": "FRAX",
          "decimals": 18
        },
        "rpcUrls": [
          "https://rpc.frax.com"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 252,
        "symbol": "FRAX",
        "decimals": 18,
        "name": "FRAX",
        "coinKey": "FRAX",
        "logoURI": "https://static.debank.com/image/movr_token/logo_url/0x965f84d915a9efa2dd81b653e3ae736555d945f4/1f2c42cba1add081f314ee899ab67816.png",
        "priceUSD": "2.69150665"
      },
      "diamondAddress": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
      "permit2Proxy": "0x89c6340B1a1f4b25D36cd8B063D49045caF3f818"
    },
    {
      "key": "abs",
      "chainType": "EVM",
      "name": "Abstract",
      "coin": "ETH",
      "id": 2741,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/abstract.svg",
      "multicallAddress": "0xAa4De41dba0Ca5dCBb288b7cC6b708F3aaC759E7",
      "metamask": {
        "chainId": "0xab5",
        "blockExplorerUrls": [
          "https://abscan.org/"
        ],
        "chainName": "Abstract",
        "nativeCurrency": {
          "name": "Ethereum",
          "symbol": "ETH",
          "decimals": 18
        },
        "rpcUrls": [
          "https://api.mainnet.abs.xyz"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 2741,
        "symbol": "ETH",
        "decimals": 18,
        "name": "ETH",
        "coinKey": "ETH",
        "logoURI": "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2/logo.png",
        "priceUSD": "4305.92"
      },
      "diamondAddress": "0x4f8C9056bb8A3616693a76922FA35d53C056E5b3",
      "permit2Proxy": "0xbb7c30Dc5AcE1A51940b202182c07D45e3e2Cba4"
    },
    {
      "key": "rsk",
      "chainType": "EVM",
      "name": "Rootstock",
      "coin": "RBTC",
      "id": 30,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/rootstock.svg",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "metamask": {
        "chainId": "0x1e",
        "blockExplorerUrls": [
          "https://explorer.rootstock.io/"
        ],
        "chainName": "Rootstock Mainnet",
        "nativeCurrency": {
          "name": "Rootstock Smart Bitcoin",
          "symbol": "RBTC",
          "decimals": 18
        },
        "rpcUrls": [
          "https://public-node.rsk.co",
          "https://mycrypto.rsk.co"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 30,
        "symbol": "RBTC",
        "decimals": 18,
        "name": "Rootstock Smart Bitcoin",
        "coinKey": "RBTC",
        "logoURI": "https://static.debank.com/image/rsk_token/logo_url/0x542fda317318ebf1d3deaf76e0b632741a7e677d/4785a26ef5bb5df987e67ad49fc62137.png",
        "priceUSD": "110832.18983228"
      },
      "diamondAddress": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
      "permit2Proxy": "0x89c6340B1a1f4b25D36cd8B063D49045caF3f818"
    },
    {
      "key": "ape",
      "chainType": "EVM",
      "name": "Apechain",
      "coin": "APE",
      "id": 33139,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/apechain.svg",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "metamask": {
        "chainId": "0x8173",
        "blockExplorerUrls": [
          "https://apescan.io/"
        ],
        "chainName": "ApeChain",
        "nativeCurrency": {
          "name": "ApeCoin",
          "symbol": "APE",
          "decimals": 18
        },
        "rpcUrls": [
          "https://rpc.apechain.com"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 33139,
        "symbol": "APE",
        "decimals": 18,
        "name": "APE",
        "coinKey": "APE",
        "logoURI": "https://static.debank.com/image/ape_token/logo_url/ape/2357165eac1453c46f526704b51a801b.png",
        "priceUSD": "0.58038215"
      },
      "diamondAddress": "0x2dea447e7dc6cd2f10b31bF10dCB30F87E838417",
      "permit2Proxy": "0x89c6340B1a1f4b25D36cd8B063D49045caF3f818"
    },
    {
      "key": "cel",
      "chainType": "EVM",
      "name": "Celo",
      "coin": "CELO",
      "id": 42220,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/celo.svg",
      "tokenlistUrl": "https://raw.githubusercontent.com/sushiswap/default-token-list/master/tokens/celo.json",
      "faucetUrls": [
        "https://stakely.io/faucet/celo-platform",
        "https://free-online-app.com/faucet-for-eth-evm-chains/"
      ],
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "metamask": {
        "chainId": "0xa4ec",
        "blockExplorerUrls": [
          "https://celoscan.io/"
        ],
        "chainName": "Celo Mainnet",
        "nativeCurrency": {
          "name": "CELO",
          "symbol": "CELO",
          "decimals": 18
        },
        "rpcUrls": [
          "https://forno.celo.org",
          "https://celo.drpc.org"
        ]
      },
      "nativeToken": {
        "address": "0x471EcE3750Da237f93B8E339c536989b8978a438",
        "chainId": 42220,
        "symbol": "CELO",
        "decimals": 18,
        "name": "Celo native asset",
        "coinKey": "CELO",
        "logoURI": "https://s2.coinmarketcap.com/static/img/coins/64x64/5567.png",
        "priceUSD": "0.29027295"
      },
      "diamondAddress": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
      "permit2": "0x000000000022D473030F116dDEE9F6B43aC78BA3",
      "permit2Proxy": "0x89c6340B1a1f4b25D36cd8B063D49045caF3f818"
    },
    {
      "key": "etl",
      "chainType": "EVM",
      "name": "Etherlink",
      "coin": "XTZ",
      "id": 42793,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/etherlink.svg",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "relayerSupported": false,
      "metamask": {
        "chainId": "0xa729",
        "blockExplorerUrls": [
          "https://explorer.etherlink.com/"
        ],
        "chainName": "Etherlink",
        "nativeCurrency": {
          "name": "Tezos",
          "symbol": "XTZ",
          "decimals": 18
        },
        "rpcUrls": [
          "https://node.mainnet.etherlink.com"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 42793,
        "symbol": "XTZ",
        "decimals": 18,
        "name": "Tezos",
        "coinKey": "XTZ",
        "logoURI": "https://assets.coingecko.com/coins/images/976/standard/Tezos-logo.png?1696502091",
        "priceUSD": "0.7135"
      },
      "diamondAddress": "0x977474593c982cFa8b197cAE302e6d01f789435b",
      "permit2Proxy": "0x799356A518cFe5BDf874A593dF41cA2A5e1575b7"
    },
    {
      "key": "wcc",
      "chainType": "EVM",
      "name": "World Chain",
      "coin": "ETH",
      "id": 480,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/world.svg",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "metamask": {
        "chainId": "0x1e0",
        "blockExplorerUrls": [
          "https://worldscan.org/"
        ],
        "chainName": "Worldchain Mainnet",
        "nativeCurrency": {
          "name": "Ethereum",
          "symbol": "ETH",
          "decimals": 18
        },
        "rpcUrls": [
          "https://worldchain-mainnet.g.alchemy.com/public",
          "https://worldchain-mainnet.gateway.tenderly.co"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 480,
        "symbol": "ETH",
        "decimals": 18,
        "name": "ETH",
        "coinKey": "ETH",
        "logoURI": "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2/logo.png",
        "priceUSD": "4304.895592234823"
      },
      "diamondAddress": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
      "permit2": "0x000000000022D473030F116dDEE9F6B43aC78BA3",
      "permit2Proxy": "0x89c6340B1a1f4b25D36cd8B063D49045caF3f818"
    },
    {
      "key": "xdc",
      "chainType": "EVM",
      "name": "XDC",
      "coin": "XDC",
      "id": 50,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/xdc.svg",
      "multicallAddress": "0x0b1795cca8e4ec4df02346a082df54d437f8d9af",
      "relayerSupported": false,
      "metamask": {
        "chainId": "0x32",
        "blockExplorerUrls": [
          "https://xdcscan.com/"
        ],
        "chainName": "XDC",
        "nativeCurrency": {
          "name": "XDC",
          "symbol": "XDC",
          "decimals": 18
        },
        "rpcUrls": [
          "https://rpc.xdcrpc.com"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 50,
        "symbol": "XDC",
        "decimals": 18,
        "name": "XDC",
        "coinKey": "XDC",
        "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/xdc.svg",
        "priceUSD": "0.08055361981928222"
      },
      "diamondAddress": "0x055d4612Ec74aD799C6cB4dF72C0Ab8dbDBCBAfa",
      "permit2Proxy": "0x767cEE938DF7d88D3e880Cbfc35C9cAfee882f4d"
    },
    {
      "key": "mnt",
      "chainType": "EVM",
      "name": "Mantle",
      "coin": "MNT",
      "id": 5000,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/mantle.svg",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "metamask": {
        "chainId": "0x1388",
        "blockExplorerUrls": [
          "https://explorer.mantle.xyz/",
          "https://mantlescan.xyz/"
        ],
        "chainName": "Mantle",
        "nativeCurrency": {
          "name": "Mantle",
          "symbol": "MNT",
          "decimals": 18
        },
        "rpcUrls": [
          "https://rpc.mantle.xyz",
          "https://mantle-rpc.publicnode.com",
          "https://mantle.drpc.org",
          "https://mantle.public-rpc.com"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 5000,
        "symbol": "MNT",
        "decimals": 18,
        "name": "MNT",
        "coinKey": "MNT",
        "logoURI": "https://static.debank.com/image/mnt_token/logo_url/0x78c1b0c915c4faa5fffa6cabf0219da63d7f4cb8/a443c78c33704d48f06e5686bb87f85e.png",
        "priceUSD": "1.15756286"
      },
      "diamondAddress": "0x1231DEB6f5749EF6cE6943a275A1D3E7486F4EaE",
      "permit2Proxy": "0xBDfF0c1C8B0b779581C4aC3bA1F29667C366C56e"
    },
    {
      "key": "sup",
      "chainType": "EVM",
      "name": "Superposition",
      "coin": "ETH",
      "id": 55244,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/superposition.svg",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "relayerSupported": false,
      "metamask": {
        "chainId": "0xd7cc",
        "blockExplorerUrls": [
          "https://explorer.superposition.so/"
        ],
        "chainName": "Superposition",
        "nativeCurrency": {
          "name": "Ether",
          "symbol": "ETH",
          "decimals": 18
        },
        "rpcUrls": [
          "https://rpc.superposition.so"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 55244,
        "symbol": "ETH",
        "decimals": 18,
        "name": "ETH",
        "coinKey": "ETH",
        "logoURI": "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2/logo.png",
        "priceUSD": "4297.99266621"
      },
      "diamondAddress": "0x03d55A7896097801B1dE90b4E3E0392CE279180A",
      "permit2": "0x000000000022D473030F116dDEE9F6B43aC78BA3",
      "permit2Proxy": "0x1BCd304fDAd1d1D66529159b1Bc8D13C9158D586"
    },
    {
      "key": "ink",
      "chainType": "EVM",
      "name": "Ink",
      "coin": "ETH",
      "id": 57073,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/ink.svg",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "relayerSupported": true,
      "metamask": {
        "chainId": "0xdef1",
        "blockExplorerUrls": [
          "https://explorer.inkonchain.com/"
        ],
        "chainName": "Ink",
        "nativeCurrency": {
          "name": "Ether",
          "symbol": "ETH",
          "decimals": 18
        },
        "rpcUrls": [
          "https://rpc-gel.inkonchain.com"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 57073,
        "symbol": "ETH",
        "decimals": 18,
        "name": "ETH",
        "coinKey": "ETH",
        "logoURI": "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2/logo.png",
        "priceUSD": "4305.92"
      },
      "diamondAddress": "0x864b314D4C5a0399368609581d3E8933a63b9232",
      "permit2Proxy": "0x1BCd304fDAd1d1D66529159b1Bc8D13C9158D586"
    },
    {
      "key": "boc",
      "chainType": "EVM",
      "name": "BOB",
      "coin": "ETH",
      "id": 60808,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/bob.svg",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "relayerSupported": false,
      "metamask": {
        "chainId": "0xed88",
        "blockExplorerUrls": [
          "https://explorer.gobob.xyz/"
        ],
        "chainName": "BOB",
        "nativeCurrency": {
          "name": "ETH",
          "symbol": "ETH",
          "decimals": 18
        },
        "rpcUrls": [
          "https://rpc.gobob.xyz"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 60808,
        "symbol": "ETH",
        "decimals": 18,
        "name": "ETH",
        "coinKey": "ETH",
        "logoURI": "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2/logo.png",
        "priceUSD": "4305.92"
      },
      "diamondAddress": "0x452Cf1B8597E6319Cd21abd847312bF17E26d8d1",
      "permit2Proxy": "0x1BCd304fDAd1d1D66529159b1Bc8D13C9158D586"
    },
    {
      "key": "kat",
      "chainType": "EVM",
      "name": "Katana",
      "coin": "ETH",
      "id": 747474,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/katana.svg",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "relayerSupported": false,
      "metamask": {
        "chainId": "0xb67d2",
        "blockExplorerUrls": [
          "https://katanascan.com/",
          "https://explorer.katana.network/"
        ],
        "chainName": "Katana",
        "nativeCurrency": {
          "name": "ETH",
          "symbol": "ETH",
          "decimals": 18
        },
        "rpcUrls": [
          "https://rpc.katana.network"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 747474,
        "symbol": "ETH",
        "decimals": 18,
        "name": "ETH",
        "coinKey": "ETH",
        "logoURI": "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2/logo.png",
        "priceUSD": "4304.7"
      },
      "diamondAddress": "0xC59fe32C9549e3E8B5dCcdAbC45BD287Bd5bA2bc",
      "permit2Proxy": "0x628d684D57c73A5D8cA77F455Fdf2CC8Bd503c16"
    },
    {
      "key": "ber",
      "chainType": "EVM",
      "name": "Berachain",
      "coin": "BERA",
      "id": 80094,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/bera.svg",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "relayerSupported": true,
      "metamask": {
        "chainId": "0x138de",
        "blockExplorerUrls": [
          "https://berascan.com/",
          "https://beratrail.io/"
        ],
        "chainName": "Berachain",
        "nativeCurrency": {
          "name": "Bera",
          "symbol": "BERA",
          "decimals": 18
        },
        "rpcUrls": [
          "https://rpc.berachain.com"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 80094,
        "symbol": "BERA",
        "decimals": 18,
        "name": "BERA",
        "coinKey": "BERA",
        "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/tokens/bera.svg",
        "priceUSD": "2.29884412"
      },
      "diamondAddress": "0xf909c4Ae16622898b885B89d7F839E0244851c66",
      "permit2": "0x000000000022D473030F116dDEE9F6B43aC78BA3",
      "permit2Proxy": "0x66C8F3F645Bc61ed1740067D381830f7c7925d98"
    },
    {
      "key": "kai",
      "chainType": "EVM",
      "name": "Kaia",
      "coin": "KAIA",
      "id": 8217,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/kaia.svg",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "metamask": {
        "chainId": "0x2019",
        "blockExplorerUrls": [
          "https://kaiascan.io/",
          "https://kaiascope.com/"
        ],
        "chainName": "Kaia Mainnet",
        "nativeCurrency": {
          "name": "KAIA",
          "symbol": "KAIA",
          "decimals": 18
        },
        "rpcUrls": [
          "https://public-en.node.kaia.io",
          "https://klaytn.drpc.org"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 8217,
        "symbol": "KAIA",
        "decimals": 18,
        "name": "KAIA",
        "coinKey": "KAIA",
        "logoURI": "https://static.debank.com/image/chain/logo_url/klay/4182ee077031d843a57e42746c30c072.png",
        "priceUSD": "0.1451"
      },
      "diamondAddress": "0x1255d17c1bc2f764d087536410879f2d0d8772fd",
      "permit2Proxy": "0x56E3B9Cf704991cc554902b16610255568A58B12"
    },
    {
      "key": "plu",
      "chainType": "EVM",
      "name": "Plume",
      "coin": "PLUME",
      "id": 98866,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/plume.svg",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "relayerSupported": false,
      "metamask": {
        "chainId": "0x18232",
        "blockExplorerUrls": [
          "https://explorer.plume.org/"
        ],
        "chainName": "Plume",
        "nativeCurrency": {
          "name": "PLUME",
          "symbol": "PLUME",
          "decimals": 18
        },
        "rpcUrls": [
          "https://rpc.plume.org"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 98866,
        "symbol": "PLU",
        "decimals": 18,
        "name": "Plume",
        "coinKey": "PLU",
        "logoURI": "https://static.debank.com/image/plume_token/logo_url/plume/20ecb01558edaac59ad4cd20b9ccc51d.png",
        "priceUSD": "0.0811891"
      },
      "diamondAddress": "0x6f5C8Bb0C5Fe4ECeAC40EE1C238EaB6bbb29761c",
      "permit2Proxy": "0x9C1D1db1CD3EbE4727c2a4208Ca5D65624cC9b4D"
    },
    {
      "key": "hyp",
      "chainType": "EVM",
      "name": "HyperEVM",
      "coin": "HYPE",
      "id": 999,
      "mainnet": true,
      "logoURI": "https://raw.githubusercontent.com/lifinance/types/main/src/assets/icons/chains/hyperevm.svg",
      "multicallAddress": "0xcA11bde05977b3631167028862bE2a173976CA11",
      "relayerSupported": true,
      "metamask": {
        "chainId": "0x3e7",
        "blockExplorerUrls": [
          "https://hyperevmscan.io/"
        ],
        "chainName": "HyperEVM",
        "nativeCurrency": {
          "name": "HYPE",
          "symbol": "HYPE",
          "decimals": 18
        },
        "rpcUrls": [
          "https://rpc.hyperliquid.xyz/evm",
          "https://rpc.hyperlend.finance",
          "https://hyperliquid-json-rpc.stakely.io",
          "https://hyperliquid.drpc.org",
          "https://rpc.hypurrscan.io"
        ]
      },
      "nativeToken": {
        "address": "0x0000000000000000000000000000000000000000",
        "chainId": 999,
        "symbol": "HYPE",
        "decimals": 18,
        "name": "HYPE",
        "coinKey": "HYPE",
        "logoURI": "https://static.debank.com/image/hyper_token/logo_url/hyper/0b3e288cfe418e9ce69eef4c96374583.png",
        "priceUSD": "46.989"
      },
      "diamondAddress": "0x0a0758d937d1059c356D4714e57F5df0239bce1A",
      "permit2Proxy": "0x3e851976DCA1bc02A2F35Ce926dCAFD1dEB3359b"
    }
  ]
}
```