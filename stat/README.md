# stat

`stat` 是部署在代理节点上的流量采集服务。它连接本机 Xray API，定时把用户流量增量写入本地 SQLite，再提供接口给 `xray-admin` 拉取。

## 编译方式

### 在 Linux 服务器本机编译

如果是在目标 Linux 服务器上直接编译，使用普通编译即可：

```bash
cd stat
go mod tidy
go build -ldflags="-s -w" -o stat .
```

### 在 Windows / macOS 上编译 Linux 版本

如果是在本机电脑编译后上传到 Linux 服务器，使用交叉编译：

```bash
cd stat
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ../dist/stat .
```

当前代码使用 `modernc.org/sqlite`，是纯 Go SQLite 驱动，推荐 `CGO_ENABLED=0`。不需要 `x86_64-linux-musl-gcc`，也不建议再使用旧的 `CGO_ENABLED=1` 编译方式。

## 启动参数

```bash
./stat -port 56611 -level info -traffic-db /var/lib/xray-stat/stat.db -collect-interval 10s -log-clean-dir /root/log -log-retention-months 12 -xray-log-dir /var/log/xray -xray-log-retention-months 12
```

参数说明：

- `-port`：stat HTTP 服务端口，默认 `56611`
- `-level`：日志级别，默认 `info`
- `-traffic-db`：本地 SQLite 文件路径，默认 `/var/lib/xray-stat/stat.db`
- `-collect-interval`：本地采集间隔，默认 `10s`
- `-log-clean-dir`：xray-admin 按日期目录清理的根目录，默认 `/root/log`
- `-log-retention-months`：xray-admin 日期目录保留月数，默认 `12`
- `-xray-log-dir`：Xray 日志目录，默认 `/var/log/xray`
- `-xray-log-retention-months`：Xray 轮转日志文件保留月数，默认 `12`

`stat` 启动后会每天清理一次 `-log-clean-dir` 下的 xray-admin 一级日期目录。只会删除目录名严格为 `YYYY-MM-DD` 且早于保留期限的目录，普通文件和非日期目录不会删除。

同时会每天清理一次 `-xray-log-dir` 下的 Xray 轮转日志文件，例如 `access.log-20250616`、`access.log-20250616.gz`、`error.log.1.gz`。活跃日志文件 `access.log` 和 `error.log` 不会被改写。

## 本地 SQLite

服务启动时会自动创建 SQLite 文件和表：

- `traffic_snapshot`：保存 Xray 当前计数快照，用于计算下次增量
- `traffic_event`：保存每次采集到的用户流量增量

第一次采集只建立基线，不会把 Xray 里已有累计值直接算成新增流量。

## 提供给管理端的接口

- `GET /stat/traffic`：兼容旧采集方式，返回 Xray 当前 stats
- `POST /stat/traffic/collect`：立即采集一次并写入本地 SQLite
- `GET /stat/traffic/sync?after_id=0&limit=1000`：给 `xray-admin` 拉取本地增量事件
- `GET /stat/sysinfo`：返回节点系统信息
- `POST /conf/update`：更新 Xray 配置
