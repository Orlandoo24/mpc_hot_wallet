package handler

import (
	"demo/internal/logic"
	"demo/internal/logic/transaction"
	"demo/internal/svc"
	"demo/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

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

func TransactionSendHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		logx.WithContext(r.Context()).Infof("TransactionSendHandler")
		var req types.TransactionSendReq
		if err := httpx.Parse(r, &req); err != nil {
			logx.WithContext(r.Context()).Errorf("failed to parse request body: %v", err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		logx.WithContext(r.Context()).Infof("Request body parsed successfully: %+v", req)

		l := transaction.NewTransactionLogic(r.Context(), svcCtx)
		resp, err := l.WrapSend(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func Hello(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		resp := "hello world"
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
