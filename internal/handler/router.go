package handler

import (
	"demo/internal/svc"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	// authMiddleware := mid.NewAuthMiddleware() // Temporarily commented out for debugging
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/wallet_init",
				Handler: WalletInitHandler(serverCtx),
			},
			// --- Transaction Routes ---
			{
				Method:  http.MethodPost,
				Path:    "/transaction/hello",
				Handler: Hello(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/transaction/send",
				Handler: SendHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/transaction/approve",
				Handler: NotImplementedHandler(),
			},
			{
				Method:  http.MethodPost,
				Path:    "/transaction/revoke",
				Handler: NotImplementedHandler(),
			},
			{
				Method:  http.MethodGet, // Receive is typically a GET request to fetch address/info
				Path:    "/transaction/receive",
				Handler: NotImplementedHandler(),
			},
			{
				Method:  http.MethodPost,
				Path:    "/transaction/swap",
				Handler: SwapHandler(serverCtx),
			},
			// --- Bridge Routes ---
			{
				Method:  http.MethodPost,
				Path:    "/bridge/quote",
				Handler: BridgeQuoteHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/bridge/execute",
				Handler: BridgeHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/bridge/status",
				Handler: BridgeStatusHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/"),
		rest.WithTimeout(30000*time.Millisecond),
	)

}

// NotImplementedHandler returns a handler that responds with a "Not Implemented" error.
func NotImplementedHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httpx.Error(w, http.ErrAbortHandler)
	}
}
