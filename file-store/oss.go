package file_store

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"path"
)

//OSS

//COS
type ossAdapter struct {
	endpoint        string
	accessKeyID     string
	accessKeySecret string
	tag             string
	client          *oss.Client
}

func newOssAdapter(tag, endpoint, accessKeyID, accessKeySecret string) *ossAdapter {

	//endpoint, accessKeyID, accessKeySecret string
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		panic(err.Error())
	}

	return &ossAdapter{
		endpoint:        endpoint,
		accessKeyID:     accessKeyID,
		accessKeySecret: accessKeySecret,
		tag:             tag,
		client:          client,
	}

}

func (o *ossAdapter) Store(rootPath, bucket, key string, size int64, fs io.ReadCloser) error {

	bucketObj, err := o.client.Bucket(bucket)
	if err != nil {
		return err
	}
	fullPath := path.Join(rootPath, bucket, key)
	err = bucketObj.PutObject(fullPath, fs)
	return err

}
