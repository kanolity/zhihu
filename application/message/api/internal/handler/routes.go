// Code generated by goctl. DO NOT EDIT.
// goctl 1.8.3

package handler

import (
	"net/http"

	"go_code/zhihu/application/message/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/list",
				Handler: getMessagesHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/read",
				Handler: markAsReadHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/send",
				Handler: sendMessageHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/message"),
	)
}
