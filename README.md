# postgresql-cdc

postgreSQL CDC

## 示例`sample`

```shell
# 启动环境
make up 

# 启动服务
make sample

# 写入数据
echo "insert into sample.user(name) values('jesonouyang');" | docker-compose exec -T postgres psql -U postgres

# 查看topic列表
docker-compose exec -T broker /bin/kafka-topics --bootstrap-server=localhost:9092 --list

# 读取sample topic中的数据
docker-compose exec -T broker /bin/kafka-console-consumer --bootstrap-server=localhost:9092 --topic sample --from-beginning
```

## 参考

https://github.com/chobostar/pg_listener/
