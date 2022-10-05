# stat
> 统计用户流量和管理用户的程序

### 编译指令
```shell
# 普通编译
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -o stat *.go
# 交叉编译
GOOS=linux CGO_ENABLED=1 CC=x86_64-linux-musl-gcc go build -a -o stat *.go 
```

base
used

base = 0
used = 100

used = 10

base = 100
used = 10

used = 30
base = 100

used = 10

base = 100 + 30
used = 10

