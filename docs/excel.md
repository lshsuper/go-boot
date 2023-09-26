### excel provider提供器

**基于excelize做二次封装并读写分离操作excel**

#### 一、读提供器的使用

##### Step1：构建

```

    reader:=NewReaderProvider("路径")
    defer reader.Close()

```



##### Step2：具体使用

```

	type TestModel struct {
		Name   string  `x-col:"姓名"`
		Phone  string  `x-col:"手机号"`
		Hobby  string  `x-col:"兴趣"`
		ID     int     `x-col:"ID"`
		Credit float64 `x-col:"学分"`
	}
	arr := make([]TestModel, 0)
    //映射指定模型切片
	reader.ReadModel(&arr)

```



#### 二、写提供器的使用

##### Step1：构建

```

	writer:=NewWriterProvider()
	defer writer.Close()

```



##### Step2：具体使用

```
	//header []Header --表头定义
	//body [][]string --具体数据
	//rules []MergeRule --合并单元格策略
	writer.Writer(header,body,rules)
    writer.Save("具体本地路径")
    	or
	writer.Buffer("你定义外的io.Reader")

```

