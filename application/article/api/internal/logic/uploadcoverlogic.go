package logic

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"go_code/zhihu/application/article/api/internal/code"
	"go_code/zhihu/pkg/xcode"
	"os"
	"path"
	"strings"
	"time"

	"go_code/zhihu/application/article/api/internal/svc"
	"go_code/zhihu/application/article/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadCoverLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadCoverLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadCoverLogic {
	return &UploadCoverLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadCoverLogic) UploadCover(req *types.UploadCoverRequest) (resp *types.UploadCoverResponse, err error) {
	userId, err := l.ctx.Value("userId").(json.Number).Int64()
	if err != nil {
		logx.Errorf("l.ctx.Value error: %v", err)
		return nil, xcode.NoLogin
	}
	coverURL, err := saveCoverFile(req.Cover, userId)
	if err != nil {
		return nil, code.UploadCoverFailed
	}

	return &types.UploadCoverResponse{CoverUrl: coverURL}, nil
}

func saveCoverFile(base64Str string, userId int64) (string, error) {
	// 生成唯一文件名
	now := time.Now().Unix()
	fileName := fmt.Sprintf("cover_%d_%d.png", userId, now)
	filePath := path.Join("..", "..", "..", "static", "upload", "cover", fileName)

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
	publicPath := fmt.Sprintf("/static/upload/cover/%s", fileName)
	return publicPath, nil
}
