# go 分布式存储项目

## 项目简介
- `go-object-storage`目录下存放所有源代码
- `conf`目录下存放配置信息以及数据库初始化文件
	- `my.cnf` 配置运行canal所需配置
	- `init.sql` 创建file用户和canal用户，file用户用于业务上的操作数据库，canal用户用于canal同步数据订阅binlog时使用
- `dao`目录下存放数据库操作代码
- `final`目录下构建apiServer和dataServer两个模块。apiServer负责处理用户请求。dataServer负责执行真正意义上的数据读写。apiServer通过消息队列收集所有可用dataServer节点地址，实现了两个模块之间的解耦。
- `myes`目录下存放操作es以及canal同步数据的代码
- `src`目录下多种操作整合的代码
	- `objectstream`用于将文件对象转化为stream流`向dataServer中存储或者从dataServer中获取`
	- `rabbitmq`放置消息队列相关操作：
		- apiServer获取活跃dataServer的数量以及它们的信息(监听dataServer心跳)
		- apiServer中通过文件hash获得存放文件分片的节点地址
	- `rs`存放调用rs纠删库实现文件分片存储的代码，还包括了获取文件流和生成到存储节点的文件流的代码。这部分会在底层调用`objectstream`中的代码，通过http将apiServer和dataServer连接起来
- `tools`目录下存放初始化环境、清理测试环境和启动存储节点等操作的快捷脚本
## 项目运行
1. 进入`myes`目录，执行`docker-compose up -d`启动es相关环境
2. 根据`conf`目录下的配置(init.sql和my.cnf)，初始化数据库配置
3. 运行`tools`目录下的`inittestenv.sh`和`starttestenv.sh`脚本，设置网络环境并启动存储节点
4. 进入`final/apiServer/`目录，执行`go run .`启动apiServer,默认在本机的`10.29.2.1:12345`地址上运行
## 接口文档
- 主目录下`storage.openapi.json`文件
- 使用方法：通过`Apifox`等工具将该文件导入，导入时选择`OpenAPI/Swagger`数据格式，构建出项目的`API`文档
## 设计要点
1. 数据同步到ES时，针对MySQL中不同操作对ES进行定制化更新，保证了数据的精炼，但同时损失了一部分灵活性。
2. 后期可以做一个数据池，在同步数据时，进行简单地过滤，将过滤之后的数据全量写入ES，会多占用一些空间，但有原始数据可供使用，为之后的功能修改提供了灵活性。
