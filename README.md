# logagent_go
> 项目简介：
![image](https://user-images.githubusercontent.com/82981823/192181686-91b528ef-1988-4469-a125-4f915be7bae4.png)
动态收集部署在不同服务器的日志；通过将要监听的日志项归档发布到etcd服务器当中，从etcd服务器当中拉取要监听的日志,实现动态追踪日志信息功能

> 使用的技术栈：
 etcd kafka ini tail

> 有待完善和掌握的地方：
将日志文件存储到数据库中，并将其连接到可视化grafana中，实现收集的日志实际数据的可视化


>启动etcd
```bash
brew services start etcd
# don't want/need a background service you can just run:
ETCD_UNSUPPORTED_ARCH="arm64" /opt/homebrew/opt/etcd/bin/etcd
```

启动kafka
```bash
/opt/homebrew/opt/kafka/bin/kafka-server-start /opt/homebrew/etc/kafka/server.properties
```