package handler

import (
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go_code/zhihu/application/chat/api/internal/logic"
	"go_code/zhihu/application/chat/api/internal/svc"
	"go_code/zhihu/application/chat/api/internal/types"
)

func createSessionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateSessionReq
		fmt.Println("before parse")
		if err := httpx.Parse(r, &req); err != nil {
			fmt.Println(err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		fmt.Println("after parse")
		l := logic.NewCreateSessionLogic(r.Context(), svcCtx)
		resp, err := l.CreateSession(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
