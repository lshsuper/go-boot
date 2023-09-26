package file_store

import (
	"fmt"
	"sync"
)

//fileStoreFactory 文档存储工厂
type fileStoreFactory struct {
	rw       *sync.RWMutex
	adapters map[StoreType]map[string]Adapter
}

//NewFileStoreFactory 构建文档存储工厂
func NewFileStoreFactory() *fileStoreFactory {
	return &fileStoreFactory{adapters: make(map[StoreType]map[string]Adapter), rw: new(sync.RWMutex)}
}

func (factory *fileStoreFactory) UseLocal(tag string) *fileStoreFactory {

	defer factory.rw.Unlock()
	factory.rw.Lock()
	if factory.adapters[LocalStore] == nil {
		factory.adapters[LocalStore] = make(map[string]Adapter)
	}

	if factory.adapters[LocalStore][tag] != nil {
		panic("tag:已存在，不能重复注册")

	}
	factory.adapters[LocalStore][tag] = newLocalAdapter(tag)

	return factory
}

func (factory *fileStoreFactory) UseMinio(tag, endpoint, accessKeyID, secretAccessKey string) *fileStoreFactory {

	defer factory.rw.Unlock()
	factory.rw.Lock()

	if factory.adapters[MinioStore] == nil {
		factory.adapters[MinioStore] = make(map[string]Adapter)
	}

	if factory.adapters[MinioStore][tag] != nil {
		panic("tag:已存在，不能重复注册")
	}

	factory.adapters[MinioStore][tag] = newMinioAdapter(tag, endpoint, accessKeyID, secretAccessKey)

	return factory
}

//UseCos 启动cos
func (factory *fileStoreFactory) UseCos(tag, bucketURL, secretId, secretKey string) *fileStoreFactory {

	defer factory.rw.Unlock()
	factory.rw.Lock()

	if factory.adapters[CosStore] == nil {
		factory.adapters[CosStore] = make(map[string]Adapter)
	}

	if factory.adapters[CosStore][tag] != nil {
		panic("tag:已存在，不能重复注册")
	}

	factory.adapters[CosStore][tag] = newCosAdapter(tag, bucketURL, secretId, secretKey)

	return factory
}

//UseOss oss注册
func (factory *fileStoreFactory) UseOss(tag, endpoint, accessKeyID, secretAccessKey string) *fileStoreFactory {

	defer factory.rw.Unlock()
	factory.rw.Lock()

	if factory.adapters[OssStore] == nil {
		factory.adapters[OssStore] = make(map[string]Adapter)
	}

	if factory.adapters[OssStore][tag] != nil {
		panic("tag:已存在，不能重复注册")
	}

	factory.adapters[OssStore][tag] = newOssAdapter(tag, endpoint, accessKeyID, secretAccessKey)

	return factory
}

//UseBos bos注册
func (factory *fileStoreFactory) UseBos(tag, endpoint, ak, sk string) *fileStoreFactory {

	defer factory.rw.Unlock()
	factory.rw.Lock()

	if factory.adapters[BosStore] == nil {
		factory.adapters[BosStore] = make(map[string]Adapter)
	}

	if factory.adapters[BosStore][tag] != nil {
		panic("tag:已存在，不能重复注册")
	}

	factory.adapters[BosStore][tag] = newBosAdapter(tag, endpoint, ak, sk)

	return factory
}

//Adapter 获取适配器
func (factory *fileStoreFactory) Adapter(st StoreType, tag string) Adapter {
	defer factory.rw.RUnlock()
	factory.rw.RLock()
	if factory.adapters[st] == nil {
		panic(fmt.Sprintf("未注册任何[%s]实例", st.String()))
	}
	adapter := factory.adapters[st][tag]
	if adapter == nil {
		panic(fmt.Sprintf("指定tag:[%s]实例不存在", tag))
	}
	return factory.adapters[st][tag]
}

//LocalAdapter 本地适配器
func (factory *fileStoreFactory) LocalAdapter(tag string) *localAdapter {
	adapter := factory.Adapter(LocalStore, tag)
	return adapter.(*localAdapter)
}

//MinioAdapter 获取指定minio适配器
func (factory *fileStoreFactory) MinioAdapter(tag string) *minioAdapter {
	adapter := factory.Adapter(MinioStore, tag)
	return adapter.(*minioAdapter)
}

//CosAdapter 获取指定Cos适配器
func (factory *fileStoreFactory) CosAdapter(tag string) *cosAdapter {
	adapter := factory.Adapter(CosStore, tag)
	return adapter.(*cosAdapter)
}

//OssAdapter 获取指定Oss适配器
func (factory *fileStoreFactory) OssAdapter(tag string) *ossAdapter {
	adapter := factory.Adapter(OssStore, tag)
	return adapter.(*ossAdapter)
}

//BosAdapter 获取指定Oss适配器
func (factory *fileStoreFactory) BosAdapter(tag string) *bosAdapter {
	adapter := factory.Adapter(BosStore, tag)
	return adapter.(*bosAdapter)
}
