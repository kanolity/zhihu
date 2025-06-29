// Code generated by goctl. DO NOT EDIT.
// goctl 1.8.3

package handler

import (
	"net/http"

	"go_code/zhihu/application/user/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/login",
				Handler: LoginHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/register",
				Handler: RegisterHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/verification",
				Handler: VerificationHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/user"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPut,
				Path:    "/change_avatar",
				Handler: ChangeAvatarHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/change_password",
				Handler: ChangePasswordHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/info",
				Handler: UserInfoHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/otherinfo",
				Handler: GetOtherInfoHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
		rest.WithSignature(serverCtx.Config.Signature),
		rest.WithPrefix("/api/user"),
	)
}
