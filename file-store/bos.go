package file_store

import (
	"github.com/baidubce/bce-sdk-go/bce"
	"github.com/baidubce/bce-sdk-go/services/bos"
	"github.com/baidubce/bce-sdk-go/services/bos/api"
	"io"
	"path"
)

//BOS
type bosAdapter struct {
	endpoint string
	ak       string
	sk       string
	tag      string
	client   *bos.Client
}

func newBosAdapter(tag, endpoint, ak, sk string) *bosAdapter {

	clientConfig := bos.BosClientConfiguration{
		Ak:               ak,
		Sk:               sk,
		Endpoint:         endpoint,
		RedirectDisabled: false,
	}

	// 初始化一个BosClient
	bosClient, err := bos.NewClientWithConfig(&clientConfig)
	if err != nil {
		panic(err)
	}
	return &bosAdapter{
		endpoint: endpoint,
		ak:       ak,
		sk:       sk,
		tag:      tag,
		client:   bosClient,
	}

}

func (o *bosAdapter) Store(rootPath, bucket, key string, size int64, fs io.ReadCloser) error {

	fullPath := path.Join(rootPath, bucket, key)
	stream, err := bce.NewBodyFromSizedReader(fs, size)
	if err != nil {
		return err
	}
	_, err = o.client.PutObject(bucket, fullPath, stream, &api.PutObjectArgs{})

	return err

}
