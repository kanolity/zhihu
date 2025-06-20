package logic

import (
	"context"
	"encoding/json"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
	"go_code/zhihu/application/like/mq/internal/svc"

	"github.com/zeromicro/go-queue/kq"
)

type ThumbupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewThumbupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThumbupLogic {
	return &ThumbupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ThumbupLogic) Consume(ctx context.Context, key, val string) error {
	var msg struct {
		BizId    string `json:"bizId"`
		TargetId int64  `json:"targetId"`
		UserId   int64  `json:"userId"`
		Type     int32  `json:"type"`
		Time     int64  `json:"time"` // optional
	}
	if err := json.Unmarshal([]byte(val), &msg); err != nil {
		l.Errorf("[MQ][Thumbup] 解码失败 key=%s val=%s err=%v", key, val, err)
		return nil
	}

	l.Infof("[MQ][Thumbup] 收到消息 → BizId=%s TargetId=%d UserId=%d Type=%d", msg.BizId, msg.TargetId, msg.UserId, msg.Type)

	return nil
}

func Consumers(ctx context.Context, svcCtx *svc.ServiceContext) []service.Service {
	return []service.Service{
		kq.MustNewQueue(svcCtx.Config.KqConsumerConf, NewThumbupLogic(ctx, svcCtx)),
	}
}
