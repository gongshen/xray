# 介绍
`gin-vue-admin` 是使用gin-vue-admin做的管理前端
`stat` 是用来收集流量的后端程序

# 一键安装
```shell
wget -N --no-check-certificate -q -O install.sh "https://raw.githubusercontent.com/gongshen/dino/main/base/install.sh" && chmod +x install.sh && bash install.sh
```

# 统计信息
### 统计信息访问地址
- 【服务器ip地址】:8080
- 【域名】:【暴露的端口】

### 统计信息接口
- /user/stat

    查看用户流量使用信息

- /user/new?name=【用户名】

    新增用户，name后面填写新增的用户名

- /user/del?name=【用户名】

    删除用户，name后面填写删除的用户名

- /user/share?name=【用户名】

    展示用户分享信息
- /user/reset_stat
  
  手动重置流量

### 统计信息其他功能
- 流量每个月定时重置
- 支持程序重启流量累计计算
- 支持v2rayN使用

### 样例
```http request
http://127.0.0.1:8080/user/stat
```
```http request
http://mydomain.com:64002/user/stat
```

### TODO（后续迭代需要实现的功能）
- 每个月定时流量重置后通过telegram上报用户流量统计信息
- 流量重置前保存历史数据
- 支持手动指定用户充值流量