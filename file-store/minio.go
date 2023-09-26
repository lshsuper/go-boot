package file_store

import (
	"github.com/minio/minio-go/v6"
	"io"
	"path"
)

type minioAdapter struct {
	endpoint        string
	accessKeyID     string
	secretAccessKey string
	client          *minio.Client
	tag             string
}

func newMinioAdapter(tag, endpoint, accessKeyID, secretAccessKey string) *minioAdapter {
	client, err := minio.New(endpoint, accessKeyID, secretAccessKey, true)
	if err != nil {
		panic(err.Error())
	}
	return &minioAdapter{
		endpoint:        endpoint,
		accessKeyID:     accessKeyID,
		secretAccessKey: secretAccessKey,
		client:          client,
		tag:             tag,
	}
}

func (o *minioAdapter) Store(rootPath, bucket, key string, size int64, fs io.ReadCloser) error {

	fullPath := path.Join(rootPath, bucket, key)
	_, err := o.client.PutObject(bucket, fullPath, fs, size, minio.PutObjectOptions{})
	return err
}
