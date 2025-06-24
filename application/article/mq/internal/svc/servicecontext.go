package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
	"go_code/zhihu/application/article/mq/internal/config"
	"go_code/zhihu/application/article/mq/internal/model"
	"go_code/zhihu/application/user/rpc/userclient"
)

type ServiceContext struct {
	Config       config.Config
	ArticleModel model.ArticleModel
	BizRedis     *redis.Redis
	UserRPC      userclient.User
	//Es           *es.Es
}

func NewServiceContext(c config.Config) *ServiceContext {
	rds, err := redis.NewRedis(redis.RedisConf{
		Host: c.BizRedis.Host,
		Pass: c.BizRedis.Pass,
		Type: c.BizRedis.Type,
	})
	if err != nil {
		panic(err)
	}

	conn := sqlx.NewMysql(c.Datasource)
	return &ServiceContext{
		Config:       c,
		ArticleModel: model.NewArticleModel(conn),
		BizRedis:     rds,
		UserRPC:      userclient.NewUser(zrpc.MustNewClient(c.UserRPC)),
		//Es: es.MustNewEs(&es.Config{
		//	Addresses: c.Es.Addresses,
		//	Username:  c.Es.Username,
		//	Password:  c.Es.Password,
		//}),
	}
}
