# xray
版本：v1.4.2

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

安装脚本
wget -N --no-check-certificate -q -O install.sh "https://raw.githubusercontent.com/gongshen/xray/main/base/install.sh" && chmod +x install.sh && bash install.sh
x-ui地址: http://74.211.97.251:64001/gongshen/

注意：
1. VMESS没有fallback的功能
2. 使用Xray1.4.2版本


TODO
1. 打印内容需要加上AlterID
2. 打印内容中文乱码
3. 支持新增后返回链接或者二维码