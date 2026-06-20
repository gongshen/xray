# Deploy README

`deploy/install.sh` 是项目的一键部署脚本，用于安装和配置 Xray、stat、xray-admin、MySQL、防火墙和证书工具。

## 下载并运行

在目标服务器上执行：

```bash
wget -O install.sh https://raw.githubusercontent.com/gongshen/xray/main/deploy/install.sh
chmod +x install.sh
bash install.sh
```

如果服务器没有 `wget`，可以使用 `curl`：

```bash
curl -L -o install.sh https://raw.githubusercontent.com/gongshen/xray/main/deploy/install.sh
chmod +x install.sh
bash install.sh
```

脚本运行后会进入交互式菜单。安装 xray-admin、stat 时会提示输入 GitHub Release 版本号，例如 `v1.0.0`。脚本会从对应 Release 下载二进制文件：

```text
https://github.com/gongshen/xray/releases/download/<version>/xray-admin
https://github.com/gongshen/xray/releases/download/<version>/stat
```

## 常用安装流程

### 管理端服务器

管理端服务器运行 xray-admin 和数据库。

推荐菜单顺序：

```text
5. 安装 xray-admin 管理端
6. 安装 MySQL
7. 初始化管理端数据库
8. 配置 iptables 防火墙
```

安装完成后，默认访问地址：

```text
http://<server-ip>:8888
```

如果在 `xray-admin/config.yaml` 中配置了证书，并把 `system.addr` 改成 `443`，则访问：

```text
https://<domain>
```

### 代理节点服务器

代理节点服务器运行 Xray 和 stat。

推荐菜单顺序：

```text
1. 安装 Xray
2. 配置 Xray
4. 安装 Stat 服务
11. 配置 Xray 日志轮转（新安装或已有节点都建议执行一次）
8. 配置 iptables 防火墙
```

配置 stat 时，需要填写管理端 IP，脚本会把该 IP 写入 `stat.service` 的 `REMOTE_IP`，用于限制只有管理端可以拉取节点统计数据。

## 菜单说明

| 选项 | 功能 |
| --- | --- |
| 1 | 安装 Xray 及依赖 |
| 2 | 生成或更新 Xray 配置 |
| 3 | 安装 BBR 网络加速 |
| 4 | 安装 stat 流量统计服务 |
| 5 | 安装 xray-admin 管理端 |
| 6 | 安装 MySQL |
| 7 | 初始化 xray-admin 数据库 |
| 8 | 配置 iptables 防火墙 |
| 9 | 安装 acme.sh 证书工具 |
| 10 | 续期 SSL 证书 |
| 11 | 配置 Xray 日志轮转和时区 |
| 12 | 分析用户流量明细 |
| 99 | 退出 |

### 4. 安装 Stat 服务

菜单 `4` 会安装代理节点上的 `stat` 服务：

- 从 GitHub Release 下载 `stat` 二进制到 `/usr/local/bin/stat`。
- 创建 `/etc/systemd/system/stat.service`。
- 提示输入管理端 IP，并写入 `REMOTE_IP`，用于限制只有管理端可以访问 stat 接口。
- 提示输入 stat 监听端口，默认 `56611`。
- 本地 SQLite 默认路径为 `/var/lib/xray-stat/stat.db`。
- 默认本地采集间隔为 `10s`，写入 systemd 参数 `-collect-interval 10s`。
- 提示输入 stat 接口自身请求/响应流量归属的用户 tag，默认 `1`，写入 systemd 参数 `-stat-api-traffic-tag 1`。
- `stat.service` 会写入 `TZ=Asia/Shanghai`；`traffic_event.collected_at` 存储 Unix 秒，菜单分析时会按北京时间转换显示和过滤。
- stat 会连接本机 Xray API 读取用户流量统计，因此代理节点需要先完成菜单 `2` 的 Xray 配置。

安装完成后可以检查：

菜单 `4` 生成的 `/etc/systemd/system/stat.service` 会显式写出 stat 当前支持的主要参数，默认值如下：

```text
-level info
-traffic-db /var/lib/xray-stat/stat.db
-collect-interval 10s
-traffic-retention-months 12
-stat-api-traffic-tag 1
-log-clean-dir /root/log
-log-retention-months 12
-xray-log-dir /var/log/xray
-xray-log-retention-months 12
```

```bash
systemctl status stat
journalctl -u stat -n 100 --no-pager
ss -ltnp | grep ':56611'
```

如果后续只想手动调整采集间隔，可以修改 systemd：

```bash
sed -i 's|-collect-interval [^ ]*|-collect-interval 10s|' /etc/systemd/system/stat.service
systemctl daemon-reload
systemctl restart stat
```

如果后续只想手动调整 stat 接口自身流量归属 tag，可以修改 systemd：

```bash
sed -i 's|-stat-api-traffic-tag [^ ]*|-stat-api-traffic-tag 1|' /etc/systemd/system/stat.service
systemctl daemon-reload
systemctl restart stat
```

### 8. 配置 iptables 防火墙

菜单 `8` 按最小权限原则生成 iptables 规则：

- 自动检测 SSH 端口，检测不到时要求手动输入。
- 询问当前服务器是否存在 xray-admin 管理端；如果存在，会检测并放行管理端端口，默认 `8888`，如果配置为 `443` 则放行 `443`。
- 询问当前服务器是否是代理节点；如果是，会检测并放行 Xray 代理端口，默认 `80`。
- 检测 stat 端口，默认 `56611`，并要求填写允许访问 stat 的管理端 IP。
- 如果 `80` 没有被管理端或代理端口使用，会询问是否额外放行 `80`，用于 HTTP/ACME standalone 等场景。
- 不放行 MySQL、Redis 等数据库端口。
- 禁用 IPv6，避免 IPv6 绕过 IPv4 防火墙规则。
- 应用规则前会备份当前 iptables，并在 120 秒内等待确认；未确认会自动恢复旧规则。

执行菜单 `8` 后，必须新开一个 SSH 窗口确认还能登录，再测试必要端口：

```bash
ssh -p <ssh-port> root@<node-ip>
```

在管理端服务器测试 stat：

```bash
curl --connect-timeout 5 http://<node-ip>:56611/stat/sysinfo
```

测试代理端口：

```bash
nc -vz <node-ip> 80
```

确认 SSH 和必要服务都可访问后，再在原窗口输入 `y` 保存为持久规则。

### 11. 配置 Xray 日志轮转和时区

菜单 `11` 会安装/确认 `logrotate`，并写入：

```text
/etc/logrotate.d/xray
```

同时会写入 Xray 的 systemd drop-in：

```text
/etc/systemd/system/xray.service.d/timezone.conf
```

其中包含 `Environment="TZ=Asia/Shanghai"`，用于让后续 `/var/log/xray/access.log` 按北京时间写入。已有旧日志不会被改写。

轮转对象为：

```text
/var/log/xray/access.log
/var/log/xray/error.log
```

默认策略：

- `daily`：按天轮转。
- `rotate 365` 和 `maxage 365`：最多保留约一年。
- `compress` 和 `delaycompress`：旧日志压缩，刚轮转出来的文件可能到下一次轮转才压缩。
- `copytruncate`：不需要重启 Xray。
- `dateext`：轮转文件在 `/var/log/xray` 同一目录下带日期后缀。

常见轮转文件格式：

```text
access.log-20260617
access.log-20260617.gz
error.log-20260617
error.log-20260617.gz
```

手动触发日志轮转：

```bash
logrotate -f /etc/logrotate.d/xray
```

dry-run 检查配置，不实际轮转：

```bash
logrotate -d /etc/logrotate.d/xray
```

触发后检查文件：

```bash
ls -lh /var/log/xray
```

注意：已有的 `access.log` 不会按内部日期拆分。手动触发时，当前整个 `access.log` 会作为一个轮转文件切出去；后续新日志再按天自然轮转。

### 12. 分析用户流量明细

菜单 `12` 用于在代理节点本机排查某个用户的流量明细：

- 输入分析日期，例如 `20260617`。
- 输入开始时间和结束时间，格式为 `小时:分钟`，例如 `8:10` 到 `9:00`；开始和结束不能超过 2 小时，结束时间会自动包含这一整分钟，即 `9:00:59`。
- 输入用户 `tag/id`，例如 `8`。
- 查询 stat 本地 SQLite，默认 `/var/lib/xray-stat/stat.db`。
- 查询 `traffic_event` 中该用户的采集周期流量，当前默认采集间隔为 `10s`。
- 菜单输入和输出都按北京时间处理：SQLite 的 `collected_at` 查询会把北京时间转换为 Unix 秒比较，展示时再从 Unix 秒转换回北京时间。
- 同时扫描 `/var/log/xray/access.log` 和同目录下的轮转文件。
- access.log 使用后缀匹配，例如 `grep -E 'email: 8$'`，避免把 `18` 匹配成 `8`。
- 输出一张按分钟聚合的表格，每分钟展示该用户的流量采集汇总和这一分钟访问过的目标域名/IP。
- 访问目标会按根域名归一并去重，例如 `rr2---sn-3pm7dne6.googlevideo.com` 显示为 `googlevideo.com`，`android.clients.google.com` 显示为 `google.com`；公网 IP 访问保留原 IP，例如 `8.8.8.8`；本机、内网、链路本地、组播、CGNAT 等内部/保留 IP 会过滤掉。
- 不再直接打印原始采集事件和原始 access.log 明细，避免长时间段输出过多流水日志。

菜单 `12` 会读取这些 access 日志文件：

```text
access.log
access.log-20260617
access.log-20260617.gz
access.log-2026-06-17.gz
```

不会读取 `access.log.backup` 这类手工备份文件。

## 默认路径

| 内容 | 路径 |
| --- | --- |
| xray-admin 二进制 | `/usr/local/bin/xray-admin` |
| xray-admin 配置 | `/usr/local/etc/xray-admin/config.yaml` |
| xray-admin systemd | `/etc/systemd/system/xray_admin.service` |
| stat 二进制 | `/usr/local/bin/stat` |
| stat systemd | `/etc/systemd/system/stat.service` |
| stat 数据库 | `/var/lib/xray-stat/stat.db` |
| Xray 配置目录 | `/usr/local/etc/xray` |
| Xray 日志目录 | `/var/log/xray` |
| Xray 日志轮转配置 | `/etc/logrotate.d/xray` |
| Xray systemd 时区配置 | `/etc/systemd/system/xray.service.d/timezone.conf` |
| iptables 规则文件 | `/usr/local/etc/xray/iptables` |

## 默认端口

| 服务 | 默认端口 | 说明 |
| --- | --- | --- |
| xray-admin | `8888` | 管理后台 Web 服务端口 |
| stat | `56611` | 节点流量统计 HTTP 服务端口 |
| Xray | `80` | 代理服务公网入口端口 |
| Xray API | `11111` | Xray 本机内部 API 端口 |
| MySQL | `3306` | 数据库端口，防火墙脚本默认不放行 |

## 访问目标分类

服务器管理页的“用户流量分析”可以把当前结果中的访问目标域名和公网 IP 提交给硅基流动模型做用途聚合分组。需要在 `/usr/local/etc/xray-admin/config.yaml` 配置：

```yaml
silicon-flow:
  api-key: "你的 SiliconFlow API Key"
  base-url: https://api.siliconflow.cn
  model: deepseek-ai/DeepSeek-V3.2
  timeout: 90s
```

API Key 只保存在 xray-admin 后端配置文件中，前端不会直接访问硅基流动。后端会在调用模型前过滤内网、保留、链路本地、组播、CGNAT 等内部 IP。`timeout` 控制后端访问硅基流动的最长等待时间，小于 `30s` 会按 `30s` 处理，DeepSeek 分类建议先用 `90s`。

分类结果会按单个访问对象缓存在管理端 MySQL 表 `v2ray_traffic_target_classification_cache` 中。再次分类时会先查缓存，只有缓存中不存在的域名/IP 才会发送给硅基流动。

## 常用命令

查看服务状态：

```bash
systemctl status xray_admin
systemctl status stat
systemctl status xray
```

重启服务：

```bash
systemctl restart xray_admin
systemctl restart stat
systemctl restart xray
```

查看日志：

```bash
journalctl -u xray_admin -f
journalctl -u stat -f
journalctl -u xray -f
tail -f /usr/local/etc/xray-admin/log/server.log
```

查看端口监听：

```bash
ss -ltnp | grep -E ':(80|443|8888|56611|11111)\b'
```

## 证书说明

菜单 `9` 会安装 `acme.sh`，用于申请证书：

```bash
~/.acme.sh/acme.sh --issue -d your-domain.com --standalone
```

证书是否让 xray-admin 监听 HTTPS，取决于 `/usr/local/etc/xray-admin/config.yaml`：

```yaml
system:
  addr: 443
  cert-file: /path/to/fullchain.cer
  key-file: /path/to/private.key
```

如果 `addr` 仍是 `8888`，即使配置了证书，xray-admin 也是监听 `https://<domain>:8888`，不会自动变成 `443`。

## 注意事项

- 生产环境建议先确认云厂商安全组已放行需要的端口。
- 配置防火墙前，确认当前 SSH 端口填写正确，避免被锁在服务器外。
- 管理端和代理节点可以部署在同一台服务器，也可以分开部署；分开部署时 stat 需要填写管理端公网 IP。
- 若域名通过 CDN、WAF 或负载均衡接入，443 到后端端口的转发可能不在服务器本机配置中，需要同时检查云平台配置。
