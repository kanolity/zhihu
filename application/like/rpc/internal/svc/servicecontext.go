package svc

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go_code/zhihu/application/like/rpc/internal/config"
	"go_code/zhihu/application/like/rpc/internal/model"
)

type ServiceContext struct {
	Config         config.Config
	KqPusherClient *kq.Pusher
	BizRedis       *redis.Redis
	LikeModel      model.LikeModel
	LikeCountModel model.LikeCountModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config:         c,
		BizRedis:       redis.New(c.BizRedis.Host, redis.WithPass(c.BizRedis.Pass)),
		KqPusherClient: kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.Topic),
		LikeModel:      model.NewLikeModel(conn, c.CacheRedis),
		LikeCountModel: model.NewLikeCountModel(conn, c.CacheRedis),
	}
}
