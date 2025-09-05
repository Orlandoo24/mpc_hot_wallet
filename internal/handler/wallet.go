package handler

import (
	"demo/internal/svc"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 创建 eth 钱包

func WalletInitHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		httpx.Ok(w)
	}
}
