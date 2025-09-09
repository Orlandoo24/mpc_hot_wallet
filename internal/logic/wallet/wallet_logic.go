package wallet

import (
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"

	"demo/internal/constant"
	"demo/internal/model"
	"demo/internal/svc"
	"demo/internal/types"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mr-tron/base58"
	"github.com/zeromicro/go-zero/core/logx"
)

type WalletLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WalletLogic {
	return &WalletLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *WalletLogic) WalletInit(req *types.WalletInitReq) (resp *types.WalletInitResp, err error) {
	l.Infof("--- 开始处理 /wallet_init 请求, name: %s ---", req.Name)
	l.Infof("将为以下链创建钱包: %v", constant.WalletInitConfig.EnabledChains)

	var wallets []types.WalletAddress
	var failedChains []string
	successCount := 0

	// 为每个配置的链创建钱包
	for _, chain := range constant.WalletInitConfig.EnabledChains {
		l.Infof("步骤 %d: 为链 %s 生成钱包...", len(wallets)+1, chain)

		walletAddr, createErr := l.createWalletForChain(string(chain), req)
		if createErr != nil {
			l.Errorf("为链 %s 创建钱包失败: %v", chain, createErr)
			failedChains = append(failedChains, string(chain))
			continue
		}

		wallets = append(wallets, *walletAddr)
		successCount++
		l.Infof("✅ 链 %s 钱包创建成功: %s", chain, walletAddr.Address)
	}

	// 检查是否有成功创建的钱包
	if successCount == 0 {
		l.Errorf("所有链的钱包创建都失败了")
		return nil, errors.New("failed to create wallets for all chains")
	}

	// 返回成功响应
	resp = &types.WalletInitResp{
		Wallets:      wallets,
		TotalCount:   len(constant.WalletInitConfig.EnabledChains),
		SuccessCount: successCount,
		FailedChains: failedChains,
	}

	l.Infof("--- /wallet_init 请求处理完成, 成功创建 %d/%d 个钱包 ---", successCount, len(constant.WalletInitConfig.EnabledChains))
	return resp, nil
}

// createWalletForChain 为指定链创建单个钱包
func (l *WalletLogic) createWalletForChain(chain string, req *types.WalletInitReq) (*types.WalletAddress, error) {
	// 1. 校验请求的链是否受支持
	if !constant.IsChainSupported(chain) {
		return nil, fmt.Errorf("unsupported chain: %s", chain)
	}

	var privateKey *ecdsa.PrivateKey
	var address string
	var privateKeyHex string

	// 2. 根据链类型生成密钥和地址
	switch constant.Chain(chain) {
	case constant.ChainEVM, constant.ChainETH, constant.ChainBSC:
		// 对于 EVM 兼容链 (ETH, BSC等)，密钥生成逻辑相同
		var errGen error
		privateKey, errGen = crypto.GenerateKey()
		if errGen != nil {
			return nil, fmt.Errorf("failed to generate EVM private key for %s: %v", chain, errGen)
		}

		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			return nil, errors.New("failed to cast public key to ECDSA")
		}
		address = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
		privateKeyHex = hex.EncodeToString(crypto.FromECDSA(privateKey))

	case constant.ChainSOLANA:
		// Solana 使用 Ed25519 加密算法
		publicKeyEd25519, privateKeyEd25519, errGen := ed25519.GenerateKey(rand.Reader)
		if errGen != nil {
			return nil, fmt.Errorf("failed to generate Solana private key: %v", errGen)
		}

		// Solana 地址就是公钥的 Base58 编码
		address = base58.Encode(publicKeyEd25519)
		privateKeyHex = hex.EncodeToString(privateKeyEd25519)
		l.Infof("Solana 钱包生成成功, 公钥长度: %d, 私钥长度: %d", len(publicKeyEd25519), len(privateKeyEd25519))

	case constant.ChainBTC:
		testMode := true
		if testMode {
			// btc 测试网钱包生成逻辑
			btcPrivateKey, errGen := btcec.NewPrivateKey()
			if errGen != nil {
				return nil, fmt.Errorf("failed to generate Bitcoin testnet private key: %v", errGen)
			}

			// 生成测试网络的 P2PKH (Pay-to-Public-Key-Hash) 地址
			pubKeyHash := btcutil.Hash160(btcPrivateKey.PubKey().SerializeCompressed())
			btcAddress, errAddr := btcutil.NewAddressPubKeyHash(pubKeyHash, &chaincfg.TestNet3Params) // 使用测试网络参数
			if errAddr != nil {
				return nil, fmt.Errorf("failed to generate Bitcoin testnet address: %v", errAddr)
			}

			address = btcAddress.EncodeAddress()
			privateKeyHex = hex.EncodeToString(btcPrivateKey.Serialize())
			l.Infof("Bitcoin 测试网钱包生成成功, 使用 P2PKH 地址格式")

		} else {
			// Bitcoin 使用 secp256k1 椭圆曲线
			btcPrivateKey, errGen := btcec.NewPrivateKey()
			if errGen != nil {
				return nil, fmt.Errorf("failed to generate Bitcoin private key: %v", errGen)
			}

			// 生成 P2PKH (Pay-to-Public-Key-Hash) 地址
			pubKeyHash := btcutil.Hash160(btcPrivateKey.PubKey().SerializeCompressed())
			btcAddress, errAddr := btcutil.NewAddressPubKeyHash(pubKeyHash, &chaincfg.MainNetParams)
			if errAddr != nil {
				return nil, fmt.Errorf("failed to generate Bitcoin address: %v", errAddr)
			}

			address = btcAddress.EncodeAddress()
			privateKeyHex = hex.EncodeToString(btcPrivateKey.Serialize())
			l.Infof("Bitcoin 钱包生成成功, 使用 P2PKH 地址格式")
		}

	default:
		return nil, fmt.Errorf("initialization logic not implemented for chain: %s", chain)
	}

	// 3. 准备数据并存入数据库
	// !!! 警告: 在生产环境中，私钥在存入数据库前必须经过强加密 !!!
	// !!! 推荐使用 KMS 管理加密密钥 !!!
	newWallet := &model.Wallets{
		// TODO: 从认证中间件中获取真实的 UserID
		UserId:              "_test_user_id_",
		Address:             address,
		EncryptedPrivateKey: privateKeyHex, // 字段名提醒我们它应该是加密的
		PhoneNumber:         sql.NullString{String: req.PhoneNumber, Valid: req.PhoneNumber != ""},
		Email:               sql.NullString{String: req.Email, Valid: req.Email != ""},
		ChainType:           sql.NullString{String: chain, Valid: true},
	}

	err := l.svcCtx.WalletsDao.Insert(l.ctx, newWallet)
	if err != nil {
		return nil, fmt.Errorf("failed to save wallet to database: %v", err)
	}

	return &types.WalletAddress{
		Chain:   chain,
		Address: address,
	}, nil
}
