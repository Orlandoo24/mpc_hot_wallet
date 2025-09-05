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
				Handler: TransactionSendHandler(serverCtx),
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
				Handler: NotImplementedHandler(),
			},
			{
				Method:  http.MethodPost,
				Path:    "/transaction/bridge",
				Handler: NotImplementedHandler(),
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
