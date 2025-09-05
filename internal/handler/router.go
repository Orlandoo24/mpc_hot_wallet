package handler

import (
	"demo/internal/mid"
	"demo/internal/svc"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	authMiddleware := mid.NewAuthMiddleware()
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{authMiddleware.Handle},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/wallet_init",
					Handler: WalletInitHandler(serverCtx),
				},
			}...,
		),
		rest.WithPrefix("/api/"),
		rest.WithTimeout(30000*time.Millisecond),
	)

}
