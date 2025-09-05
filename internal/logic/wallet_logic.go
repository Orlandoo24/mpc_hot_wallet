package logic

import (
	"context"

	"demo/internal/svc"

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
