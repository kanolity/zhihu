package logic

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"go_code/zhihu/application/user/api/internal/code"
	"go_code/zhihu/application/user/rpc/userclient"
	"os"
	"path"
	"strings"
	"time"

	"go_code/zhihu/application/user/api/internal/svc"
	"go_code/zhihu/application/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeAvatarLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangeAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeAvatarLogic {
	return &ChangeAvatarLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangeAvatarLogic) ChangeAvatar(req *types.ChangeAvatarRequest) (resp *types.ChangeAvatarResponse, err error) {
	userId, err := l.ctx.Value(types.UserIdKey).(json.Number).Int64()
	if err != nil {
		return nil, err
	}

	// 判断 Avatar 是否为上传路径
	if strings.HasPrefix(req.Avatar, "data:image") {
		filePath, err := saveAvatarFile(req.Avatar, userId)
		if err != nil {
			return nil, code.AvatarUploadFailed
		}
		// 替换成保存后的访问路径
		req.Avatar = filePath
	}
	_, err = l.svcCtx.UserRpc.ChangeAvatar(l.ctx, &userclient.ChangeAvatarRequest{
		Avatar: req.Avatar,
		UserId: userId,
	})
	return &types.ChangeAvatarResponse{}, nil
}

func saveAvatarFile(base64Str string, userId int64) (string, error) {
	// 生成唯一文件名
	now := time.Now().Unix()
	fileName := fmt.Sprintf("avatar_%d_%d.png", userId, now)
	filePath := path.Join("..", "..", "..", "static", "upload", "avatar", fileName)

	// 去除前缀并解码 base64
	idx := strings.Index(base64Str, "base64,")
	if idx == -1 {
		return "", errors.New("base64 编码格式非法")
	}
	rawData := base64Str[idx+7:]
	data, err := base64.StdEncoding.DecodeString(rawData)
	if err != nil {
		return "", err
	}

	// 确保目录存在
	err = os.MkdirAll(path.Dir(filePath), os.ModePerm)
	if err != nil {
		return "", err
	}

	// 写入文件
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return "", err
	}

	// 构造对外可访问路径
	publicPath := fmt.Sprintf("/static/upload/avatar/%s", fileName)
	return publicPath, nil
}
