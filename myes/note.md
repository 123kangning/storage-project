### canal参考文献：
- [go-canal 订阅和消费 binlog](https://mrsnake.top/articles/3/)
- [Golang 处理 MySQL 的 binlog/go语言中文网](https://studygolang.com/articles/21373)
- [基于go实现mysql数据订阅](https://vlambda.com/wz_7i6hhIyR6OM.html)

### 遇到的问题：
- canal 的 ServerId 要与 MySQL配置中的 ServerId 相同，
- canal使用的MySQL用户要授予`REPLICATION SLAVE`权限，该权限允许用户进行主从复制操作
