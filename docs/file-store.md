### 对象存储动态适配器

>概述：通过工厂管理实例，通过适配器适配不同存储机制

#### 一、如何使用

##### Step1：关于工厂实例的定义
1. 可以使用包内默认实例：DefaultFactory

2.自行构建
```
   factory:=file_store.NewFileStoreFactory()
```
##### Step2：初始化项目需要使用的文件存储适配器（Adapter）

```
    //初始化本地存储
    factory.UseLocal("local")
    
    //初始化minio(tag, endpoint, accessKeyID, secretAccessKey string)
    //tag:唯一标识一个实例
    factory.UseMinio("minio","地址","秘钥key","秘钥secret")
   
```

##### Step3：通过配置去自动适配使用哪种存储方式

```
   adapter:=file_store.StoreType("从配置中读取的值")   //转换成适配器可用类型
   //rootPath, key, bucket string, size int64, fs io.ReadCloser
   factory.Adapter(adapter).Store("根目录","文件key","存储桶名","文件流")
```

#### 二、关于适配器扩展

##### Step1：/file-store/adapter.go 扩展StoreType类型枚举值

##### Step2：/file-store/ 定义对应的适配器方法（需要实现/file-store/adapter.go的Adapter适配器接口，如minio.go）

##### Step3：/file-store/factory 定义对应的注册实例的方法名如：UseXXX()（eg UseMinio）


