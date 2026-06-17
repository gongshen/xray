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

## SSL 证书

证书通常由部署脚本或 Nginx/Caddy 处理。如果需要手动申请：

```bash
wget -O - https://get.acme.sh | sh -s email=my@example.com
yum install socat -y
acme.sh --issue --standalone -d your-domain.com
```
