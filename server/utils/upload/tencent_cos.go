package upload

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"github.com/icosmos-space/iadmin/server/global"

	"github.com/tencentyun/cos-go-sdk-v5"
	"go.uber.org/zap"
)

type TencentCOS struct{}

// UploadFile upload file to COS
func (*TencentCOS) UploadFile(file *multipart.FileHeader) (string, string, error) {
	client := NewClient()
	f, openError := file.Open()
	if openError != nil {
		global.IADMIN_LOG.Error("function file.Open() failed", zap.Any("err", openError.Error()))
		return "", "", errors.New("function file.Open() failed, err:" + openError.Error())
	}
	defer f.Close() // 创建文件 defer 关闭
	fileKey := fmt.Sprintf("%d%s", time.Now().Unix(), file.Filename)

	_, err := client.Object.Put(context.Background(), global.IADMIN_CONFIG.TencentCOS.PathPrefix+"/"+fileKey, f, nil)
	if err != nil {
		panic(err)
	}
	return global.IADMIN_CONFIG.TencentCOS.BaseURL + "/" + global.IADMIN_CONFIG.TencentCOS.PathPrefix + "/" + fileKey, fileKey, nil
}

// DeleteFile delete file form COS
func (*TencentCOS) DeleteFile(key string) error {
	client := NewClient()
	name := global.IADMIN_CONFIG.TencentCOS.PathPrefix + "/" + key
	_, err := client.Object.Delete(context.Background(), name)
	if err != nil {
		global.IADMIN_LOG.Error("function bucketManager.Delete() failed", zap.Any("err", err.Error()))
		return errors.New("function bucketManager.Delete() failed, err:" + err.Error())
	}
	return nil
}

// NewClient init COS client
func NewClient() *cos.Client {
	urlStr, _ := url.Parse("https://" + global.IADMIN_CONFIG.TencentCOS.Bucket + ".cos." + global.IADMIN_CONFIG.TencentCOS.Region + ".myqcloud.com")
	baseURL := &cos.BaseURL{BucketURL: urlStr}
	client := cos.NewClient(baseURL, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  global.IADMIN_CONFIG.TencentCOS.SecretID,
			SecretKey: global.IADMIN_CONFIG.TencentCOS.SecretKey,
		},
	})
	return client
}
