# server

`server` 是 `xray-admin` 后端服务。生产编译前需要先准备好前端静态资源，否则 `resource.GetPageFS()` 没有页面文件可嵌入。

## 编译前准备前端资源

从项目根目录执行：

```bash
cd gin-vue-admin/web
npm install --legacy-peer-deps
npm run build

rm -rf ../server/resource/page
mkdir -p ../server/resource/page
cp -r dist/* ../server/resource/page/
```

确认文件存在：

```bash
test -f ../server/resource/page/index.html
```

注意：不要复制成 `gin-vue-admin/server/dist`，当前后端不会读取这个目录。

## 编译 xray-admin

### 在 Linux 服务器本机编译

```bash
cd gin-vue-admin/server
go mod tidy
go build -ldflags="-s -w" -o xray-admin .
```

### 在 Windows / macOS 上编译 Linux 版本

```bash
cd gin-vue-admin/server
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ../../dist/xray-admin .
```

在 Windows PowerShell 中可以这样写：

```powershell
cd D:\go\pkg\xray\gin-vue-admin\server
$env:CGO_ENABLED="0"
$env:GOOS="linux"
$env:GOARCH="amd64"
go build -ldflags="-s -w" -o ../../dist/xray-admin .
```

## 运行

```bash
./xray-admin -c /usr/local/etc/xray-admin/config.yaml
```

如果不传 `-c`，程序会按 gin-vue-admin 默认规则读取当前目录下的配置文件。

## xray-admin 配置项

`/usr/local/etc/xray-admin/config.yaml` 中和节点采集相关的配置：

```yaml
stat_port: 56611
traffic_collect_interval: 1h
sysinfo_collect_interval: 5m
traffic-meter:
  enable: true
  stat-url: http://127.0.0.1:56611
  tag: "1"
  flush-interval: 10s
silicon-flow:
  api-key: ""
  base-url: https://api.siliconflow.cn
  model: deepseek-ai/DeepSeek-V3.2
  timeout: 30s
```

- `stat_port`：节点 `stat` 服务默认端口，单个服务器记录未单独配置端口时使用。
- `traffic_collect_interval`：`xray-admin` 拉取各节点本地 SQLite 流量事件的间隔，支持 Go duration 格式，例如 `10m`、`30m`、`1h`；未配置或配置非法时默认 `1h`。
- `sysinfo_collect_interval`：`xray-admin` 刷新节点在线状态、CPU、内存、磁盘信息的间隔，支持 Go duration 格式；未配置或配置非法时默认 `5m`。前端超过 10 分钟未收到更新会显示离线。

## 静态资源说明

后端路由中通过 `/fe` 挂载前端页面：

```go
Router.StaticFS("/fe", resource.GetPageFS())
```

因此发布时需要保证：

- `server/resource/resource.go` 存在并参与编译
- `server/resource/page/index.html` 存在
- `server/resource/page/assets/` 存在

## 打包上传示例

```bash
mkdir -p ../../dist
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ../../dist/xray-admin .
tar -czvf ../../dist/xray-admin.tar.gz -C ../../dist xray-admin
```

如果使用外部 `config.yaml`，建议单独上传并放到 `/usr/local/etc/xray-admin/config.yaml`。

## traffic-meter

`traffic-meter` counts xray-admin HTTP request/response bytes and batches them to stat:

```yaml
traffic-meter:
  enable: true
  stat-url: http://127.0.0.1:56611
  tag: "1"
  flush-interval: 10s
```

The target stat service writes these batches into local SQLite `traffic_event`; default tag is `1`.

`silicon-flow` is used by the server traffic analysis page to group access target domains and public IPs by actual usage. Put the SiliconFlow API key in `api-key`; the browser never receives this key. Internal, private, link-local, multicast, reserved, and CGNAT IPs are filtered before calling the model.

## SSL 证书

证书通常由部署脚本或 Nginx/Caddy 处理。如果需要手动申请：

```bash
wget -O - https://get.acme.sh | sh -s email=my@example.com
yum install socat -y
acme.sh --issue --standalone -d your-domain.com
```
