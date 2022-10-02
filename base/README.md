# xray

[增加pprof](https://github.com/XTLS/Xray-core/pull/1000)

[流量统计官网配置介绍](https://xtls.github.io/config/stats.html#statsobject)

[流量统计踩坑](https://bytemeta.vip/repo/XTLS/Xray-core/issues/687)

```shell
# 命令行获取流量
xray api stats --name="inbound>>>api>>>traffic>>>downlink" --server=127.0.0.1:port
xray api stats --name="user>>>love@v2fly.org>>>traffic>>>uplink" --server=127.0.0.1:port
```

安装x-ui
```shell
bash <(curl -Ls https://raw.githubusercontent.com/vaxilu/x-ui/master/install.sh)
```

x-ui地址: http://74.211.97.251:64001/gongshen/

编译：GOOS=linux GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-linux-musl-gcc CGO_LDFLAGS="-static" go build -a -o stat *.go