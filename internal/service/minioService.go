package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"mime/multipart"
	"net/url"
	"path"
	"strings"
	"time"
	"vibe-music-server/internal/config"
)

type MinioService struct {
	client   *minio.Client
	bucket   string
	endpoint string
	ctx      context.Context
}

func NewMinioService() *MinioService {
	MinioConf := config.Get().Minio
	client, err := minio.New(MinioConf.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(MinioConf.AccessKey, MinioConf.SecretKey, ""),
		Secure: MinioConf.UseSSL,
	})
	if err != nil {
		panic(err)
	}
	return &MinioService{
		client:   client,
		bucket:   MinioConf.Bucket,
		endpoint: MinioConf.Endpoint,
		ctx:      context.Background(),
	}
}

func (m MinioService) UploadFile(file *multipart.FileHeader, folder string) (string, error) {
	// 1. 构造对象名：folder/UUID-原文件名
	objectName := path.Join(folder,
		uuid.NewString()+"-"+file.Filename)

	// 2. 打开文件流
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("open multipart: %w", err)
	}
	defer src.Close()

	// 3. 上传（大小已知，用 file.Size）
	ctx, cancel := context.WithTimeout(m.ctx, time.Second*5)
	defer cancel()
	_, err = m.client.PutObject(ctx,
		m.bucket,
		objectName,
		src,
		file.Size,
		minio.PutObjectOptions{
			ContentType: file.Header.Get("Content-Type"),
		})
	if err != nil {
		return "", fmt.Errorf("put object: %w", err)
	}

	// 4. 返回拼接好的访问 URL
	return fmt.Sprintf("%s/%s/%s", m.endpoint, m.bucket, objectName), nil
}

func (m MinioService) DeleteFile(fileURL string) error {
	// 1. 解析 URL
	u, err := url.Parse(fileURL)
	if err != nil {
		return fmt.Errorf("invalid fileURL: %w", err)
	}

	// 2. 去掉最前面 "/" 得到  bucket+对象 路径
	fullPath := strings.TrimPrefix(u.Path, "/") // vibe-music-data/img/a/b.jpg

	// 3. 去掉 bucket 前缀，拿到纯对象名
	if !strings.HasPrefix(fullPath, m.bucket+"/") {
		return errors.New("url does not contain expected bucket")
	}
	objectName := strings.TrimPrefix(fullPath, m.bucket+"/") // img/a/b.jpg

	// 4. 删除对象
	ctx, cancel := context.WithTimeout(m.ctx, 5*time.Second)
	defer cancel()

	return m.client.RemoveObject(ctx, m.bucket, objectName, minio.RemoveObjectOptions{})
}
