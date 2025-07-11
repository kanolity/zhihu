package handler

import (
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go_code/zhihu/application/qa/api/internal/logic"
	"go_code/zhihu/application/qa/api/internal/svc"
	"go_code/zhihu/application/qa/api/internal/types"
)

func getQuestionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetQuestionsReq
		fmt.Println("before parse")
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		fmt.Println("after parse")

		l := logic.NewGetQuestionsLogic(r.Context(), svcCtx)
		resp, err := l.GetQuestions(&req)
		fmt.Printf("resp:%+v, err:%+v", resp, err)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
