# base
> 安装xray需要的杂七杂八集合

### xray版本：v1.4.2

### 额外信息
[增加pprof](https://github.com/XTLS/Xray-core/pull/1000)
[流量统计官网配置介绍](https://xtls.github.io/config/stats.html#statsobject)
[流量统计踩坑](https://bytemeta.vip/repo/XTLS/Xray-core/issues/687)

### xray api使用
```shell
# 命令行获取流量
xray api stats --name="inbound>>>api>>>traffic>>>downlink" --server=127.0.0.1:port
xray api stats --name="user>>>love@v2fly.org>>>traffic>>>uplink" --server=127.0.0.1:port
```

### xray安装脚本
wget -N --no-check-certificate -q -O install.sh "https://raw.githubusercontent.com/gongshen/xray/main/base/install.sh" && chmod +x install.sh && bash install.sh

### 注意事项：
1. VMESS没有fallback的功能
2. 使用Xray1.4.2版本