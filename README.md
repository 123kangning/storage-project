# go 分布式存储项目

在项目中构建ApiServer和DataServer两个模块。ApiServer负责处理用户请求。
DataServer负责执行真正意义上的数据读写。ApiServer通过消息队列收集所有可用的DataServer节点地址，
实现了两个模块之间的解耦。引入Elasticsearch存储文件元数据信息以及hash值，
使得在进行大数据量处理时仍然具有较高的搜索性能。

Flink SQL:
```html
CREATE TABLE file (
    `id` INTEGER NOT NULL ,
     `name` VARCHAR(255) NOT NULL,
     `size` INTEGER NOT NULL,
     `hash` CHAR(64) NOT NULL,
     `is_delete` BOOLEAN NOT NULL,
     `update_at` TIMESTAMP_LTZ NOT NULL,
     PRIMARY KEY (`id`) NOT ENFORCED
) WITH (
   'connector' = 'mysql-cdc',
    'hostname' = 'localhost',
    'port' = '3306',
    'username' = 'file',
    'password' = 'file',
    'database-name' = 'file',
    'table-name' = 'file',
    'format' = 'UTF-8'
);

CREATE TABLE esfile (
   id INT NOT NULL,
   name STRING NOT NULL,
   size INT NOT NULL,
   hash STRING NOT NULL,
   is_delete BOOLEAN NOT NULL,
   PRIMARY KEY (id) NOT ENFORCED
) WITH (
   'connector' = 'elasticsearch-7',
   'hosts' = 'http://localhost:9200',
   'index' = 'esfile',
   'document-type' = '_doc',
   'format' = 'json',
   'sink.bulk-flush.interval' = '1000',
   'sink.bulk-flush.max-size' = '5mb',
   'format'='UTF-8'
);
```

