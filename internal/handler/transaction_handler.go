package handler

import (
	"demo/internal/logic"
	"demo/internal/svc"
	"demo/internal/types"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func SendHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		logx.WithContext(r.Context()).Infof("SendHandler")
		var req types.TransactionReq
		if err := httpx.Parse(r, &req); err != nil {
			logx.WithContext(r.Context()).Errorf("failed to parse request body: %v", err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		logx.WithContext(r.Context()).Infof("Request body parsed successfully: %+v", req)

		l := logic.NewTransactionLogic(r.Context(), svcCtx)
		resp, err := l.WrapSend(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func SwapHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.WithContext(r.Context()).Infof("SwapHandler")
		var req types.TransactionReq
		if err := httpx.Parse(r, &req); err != nil {
			logx.WithContext(r.Context()).Errorf("failed to parse request body: %v", err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		logx.WithContext(r.Context()).Infof("Request body parsed successfully: %+v", req)

		l := logic.NewTransactionLogic(r.Context(), svcCtx)
		resp, err := l.WrapSwap(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

// BridgeQuoteHandler 获取跨链报价
func BridgeQuoteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.WithContext(r.Context()).Infof("BridgeQuoteHandler")
		var req types.BridgeQuoteReq
		if err := httpx.Parse(r, &req); err != nil {
			logx.WithContext(r.Context()).Errorf("failed to parse request body: %v", err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		logx.WithContext(r.Context()).Infof("Request body parsed successfully: %+v", req)

		l := logic.NewBridgeLogic(r.Context(), svcCtx)
		resp, err := l.GetBridgeQuote(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

// BridgeHandler 执行跨链转账
func BridgeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.WithContext(r.Context()).Infof("BridgeHandler")
		var req types.BridgeExecuteReq
		if err := httpx.Parse(r, &req); err != nil {
			logx.WithContext(r.Context()).Errorf("failed to parse request body: %v", err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		logx.WithContext(r.Context()).Infof("Request body parsed successfully: %+v", req)

		l := logic.NewBridgeLogic(r.Context(), svcCtx)
		resp, err := l.ExecuteBridge(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

// BridgeStatusHandler 查询跨链状态
func BridgeStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.WithContext(r.Context()).Infof("BridgeStatusHandler")
		var req types.BridgeStatusReq
		if err := httpx.Parse(r, &req); err != nil {
			logx.WithContext(r.Context()).Errorf("failed to parse request body: %v", err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		logx.WithContext(r.Context()).Infof("Request body parsed successfully: %+v", req)

		l := logic.NewBridgeLogic(r.Context(), svcCtx)
		resp, err := l.GetBridgeStatus(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
