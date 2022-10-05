# stat
> 统计用户流量和管理用户的程序

### 访问地址
> 【服务器ip地址】:8080

### 接口
- /user/stat
> 查看用户流量使用信息（如果服务重启，所有用户流量重头开始统计）
- /user/new?name=【用户名】
> 新增用户，name后面填写新增的用户名（新增用户后，所有用户的流量从头开始统计）
- /user/del?name=【用户名】
> 删除用户，name后面填写删除的用户名（删除用户后，所有用户的流量从头开始统计）
- /user/share?name=【用户名】
> 展示用户分享信息

### 编译指令
```shell
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -o stat *.go
```
