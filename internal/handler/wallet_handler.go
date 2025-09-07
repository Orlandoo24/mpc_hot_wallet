package handler

import (
	"demo/internal/logic"
	"demo/internal/svc"
	"demo/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// WalletInitHandler 创建 eth 钱包
func WalletInitHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WalletInitReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewWalletLogic(r.Context(), svcCtx)
		resp, err := l.WalletInit(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

// Hello 端点健康校验
func Hello(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := "hello world"
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
