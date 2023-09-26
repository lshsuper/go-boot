package file_store

import (
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"net/http"
	"net/url"
	"path"
)

//COS
type cosAdapter struct {
	bucketURL string
	secretId  string
	secretKey string
	tag       string
}

func newCosAdapter(tag, bucketURL, secretId, secretKey string) *cosAdapter {

	return &cosAdapter{
		bucketURL: bucketURL,
		secretId:  secretId,
		secretKey: secretKey,
		tag:       tag,
	}
}

func (o *cosAdapter) Store(rootPath, bucket, key string, size int64, fs io.ReadCloser) error {

	client, err := o.getClient()
	if err != nil {
		return err
	}
	fullPath := path.Join(rootPath, bucket, key)
	_, err = client.Object.Put(context.Background(), fullPath, fs, &cos.ObjectPutOptions{})
	return err

}

func (o *cosAdapter) getClient() (*cos.Client, error) {

	u, err := url.Parse(o.bucketURL)
	if err != nil {
		return nil, err
	}

	client := cos.NewClient(&cos.BaseURL{BucketURL: u}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  o.secretId,  // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考 https://cloud.tencent.com/document/product/598/37140
			SecretKey: o.secretKey, // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参考 https://cloud.tencent.com/document/product/598/37140
		},
	})
	return client, nil

}
