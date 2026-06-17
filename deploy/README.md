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
| 11 | 配置 Xray 日志轮转 |
| 0 | 退出 |

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
| iptables 规则文件 | `/usr/local/etc/xray/iptables` |

## 默认端口

| 服务 | 默认端口 | 说明 |
| --- | --- | --- |
| xray-admin | `8888` | 管理后台 Web 服务端口 |
| stat | `56611` | 节点流量统计 HTTP 服务端口 |
| Xray | `80` | 代理服务公网入口端口 |
| Xray API | `11111` | Xray 本机内部 API 端口 |
| MySQL | `3306` | 数据库端口，防火墙脚本默认不放行 |

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

## 防火墙说明

菜单 `8` 会生成最小放行规则：

- SSH 端口会放行。
- xray-admin 端口按配置放行，默认 `8888`。
- 如果 xray-admin 配置为监听 `443`，则 `443` 会作为 xray-admin 管理端口放行。
- Xray 代理端口按配置放行，默认 `80`。
- stat 端口默认 `56611`，并限制只允许管理端 IP 访问。
- MySQL、Redis 等数据库端口默认不放行。
- 脚本会先备份当前 iptables，再应用新规则；确认 SSH 和服务可访问后再保存为持久规则。

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
