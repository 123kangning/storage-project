### 遇到的问题：
- canal 的 ServerId 要与 MySQL配置中的 ServerId 相同，
- canal使用的MySQL用户要授予`REPLICATION SLAVE`权限，该权限允许用户进行主从复制操作
