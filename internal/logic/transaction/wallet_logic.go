package transaction

import (
	"context"
	"crypto/ecdsa"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"

	"demo/internal/constant"
	"demo/internal/model"
	"demo/internal/svc"
	"demo/internal/types"

	"github.com/ethereum/go-ethereum/crypto"
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
	l.Infof("--- 开始处理 /wallet_init 请求, chain: %s, name: %s ---", req.Chain, req.Name)

	// 1. 校验请求的链是否受支持
	l.Infof("步骤 1: 校验链是否受支持...")
	if !constant.IsChainSupported(req.Chain) {
		l.Errorf("校验失败: 不支持的链 %s", req.Chain)
		return nil, errors.New("unsupported chain")
	}
	l.Infof("校验成功.")

	var privateKey *ecdsa.PrivateKey
	var address string

	// 2. 根据常量配置的链类型生成密钥和地址
	l.Infof("步骤 2: 为链 %s 生成密钥和地址...", req.Chain)
	switch constant.Chain(req.Chain) {
	case constant.ChainETH, constant.ChainBSC:
		// 对于 EVM 兼容链 (ETH, BSC等)，密钥生成逻辑相同
		var errGen error
		privateKey, errGen = crypto.GenerateKey()
		if errGen != nil {
			l.Errorf("failed to generate EVM private key for %s: %v", req.Chain, errGen)
			return nil, errors.New("failed to generate wallet key")
		}

		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			return nil, errors.New("failed to cast public key to ECDSA")
		}
		address = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	default:
		// 此处为未来扩展非 EVM 链保留
		l.Errorf("未实现的初始化逻辑 for chain: %s", req.Chain)
		return nil, fmt.Errorf("initialization logic not implemented for chain: %s", req.Chain)
	}
	l.Infof("密钥和地址生成成功, address: %s", address)

	// 3. 将私钥转换为十六进制字符串以便存储
	// !!! 警告: 在生产环境中，私钥在存入数据库前必须经过强加密 !!!
	// !!! 推荐使用 KMS 管理加密密钥 !!!
	privateKeyHex := hex.EncodeToString(crypto.FromECDSA(privateKey))

	// 4. 准备数据并存入数据库
	l.Infof("步骤 3: 准备钱包数据并存入数据库...")
	newWallet := &model.Wallets{
		// TODO: 从认证中间件中获取真实的 UserID
		UserId:              "_test_user_id_",
		Address:             address,
		EncryptedPrivateKey: privateKeyHex, // 字段名提醒我们它应该是加密的
		PhoneNumber:         sql.NullString{String: req.PhoneNumber, Valid: req.PhoneNumber != ""},
		Email:               sql.NullString{String: req.Email, Valid: req.Email != ""},
	}

	err = l.svcCtx.WalletsDao.Insert(l.ctx, newWallet)
	if err != nil {
		l.Errorf("数据库插入失败: %v", err)
		return nil, errors.New("failed to save wallet")
	}
	l.Infof("钱包数据成功存入数据库.")

	// 5. 返回成功响应
	resp = &types.WalletInitResp{
		Address: address,
		Chain:   req.Chain,
	}
	l.Infof("--- /wallet_init 请求处理完成, 返回地址: %s ---", resp.Address)
	return resp, nil
}
