package file_store

import "io"

//Adapter 适配器
type Adapter interface {
	Store(rootPath, bucket, key string, size int64, fs io.ReadCloser) error
}

type StoreType string

const (
	LocalStore StoreType = "local" //本地存储
	MinioStore StoreType = "minio" //minio存储
	CosStore   StoreType = "cos"   //腾讯cos
	OssStore   StoreType = "oss"   //阿里oss
	BosStore   StoreType = "bos"   //百度bos
)

func (e StoreType) String() string {
	return string(e)
}

func (e StoreType) Remark() string {
	switch e {
	case LocalStore:
		return "本地存储"
	case MinioStore:
		return "Minio存储"
	case CosStore:
		return "Cos存储"
	case OssStore:
		return "Oss存储"
	case BosStore:
		return "Bos存储"
	default:
		panic("暂不支持当前类型存储")
	}
}
