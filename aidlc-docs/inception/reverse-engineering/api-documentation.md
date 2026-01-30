# API Documentation

## REST APIs

### 系统管理 API

#### 用户认证
| 方法 | 路径 | 描述 |
|------|------|------|
| POST | /base/login | 用户登录 |
| POST | /base/captcha | 获取验证码 |
| POST | /jwt/jsonInBlacklist | JWT 加入黑名单 |

#### 用户管理
| 方法 | 路径 | 描述 |
|------|------|------|
| POST | /user/register | 注册用户 |
| POST | /user/changePassword | 修改密码 |
| GET | /user/getUserList | 获取用户列表 |
| PUT | /user/setUserInfo | 设置用户信息 |
| DELETE | /user/deleteUser | 删除用户 |

#### 角色权限
| 方法 | 路径 | 描述 |
|------|------|------|
| POST | /authority/createAuthority | 创建角色 |
| DELETE | /authority/deleteAuthority | 删除角色 |
| PUT | /authority/updateAuthority | 更新角色 |
| POST | /authority/getAuthorityList | 获取角色列表 |

### 代理管理 API (v2ray_admin)

#### 服务器管理
| 方法 | 路径 | 描述 |
|------|------|------|
| POST | /v2ray_admin/server/createServer | 创建服务器 |
| DELETE | /v2ray_admin/server/deleteServer | 删除服务器 |
| DELETE | /v2ray_admin/server/deleteServerByIds | 批量删除服务器 |
| PUT | /v2ray_admin/server/updateServer | 更新服务器 |
| GET | /v2ray_admin/server/findServer | 查询服务器 |
| GET | /v2ray_admin/server/getServerList | 获取服务器列表 |
| POST | /v2ray_admin/server/getAllServer | 获取所有服务器 |
| PUT | /v2ray_admin/server/restartXray | 重启 xray 服务 |

#### 流量统计
| 方法 | 路径 | 描述 |
|------|------|------|
| POST | /v2ray_admin/stat/createStat | 创建统计记录 |
| DELETE | /v2ray_admin/stat/deleteStat | 删除统计记录 |
| GET | /v2ray_admin/stat/getStatList | 获取统计列表 |
| GET | /v2ray_admin/stat/getStatCharts | 获取统计图表数据 |
| GET | /v2ray_admin/stat/getStatRank | 获取流量排行 |

#### 用户绑定
| 方法 | 路径 | 描述 |
|------|------|------|
| POST | /v2ray_admin/binding/createBinding | 创建绑定 |
| DELETE | /v2ray_admin/binding/deleteBinding | 删除绑定 |
| PUT | /v2ray_admin/binding/updateBinding | 更新绑定 |
| GET | /v2ray_admin/binding/getBindingList | 获取绑定列表 |

### 用户代理 API (v2ray)

#### 流量查询
| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /v2ray/stat/getStatList | 获取个人流量统计 |

#### 绑定查询
| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /v2ray/binding/getBindingList | 获取个人绑定列表 |

### Stat 程序 API (运行在代理服务器上)

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | /stat/traffic | 采集流量数据 |
| POST | /conf/update | 更新 xray 配置 |

## Data Models

### Server (服务器)
```go
type Server struct {
    ID         uint            // 主键
    Ip         string          // 服务器 IP (唯一)
    Remark     string          // 备注
    Port       int64           // 端口
    ResetDate  int             // 流量重置日期
    Config     json.RawMessage // xray 配置 (JSON)
    CreatedAt  time.Time       // 创建时间
    UsedQuota  uint64          // 已使用流量
    TotalQuota uint64          // 总流量配额
}
```

### Stat (流量统计)
```go
type Stat struct {
    ID        uint   // 主键
    Tag       string // 用户标识
    Down      uint64 // 下行流量 (字节)
    Up        uint64 // 上行流量 (字节)
    ServerIp  string // 服务器 IP
    CreatedAt int    // 创建时间戳
}
```

### Binding (用户绑定)
```go
type Binding struct {
    ID        uint      // 主键
    CreatedAt time.Time // 创建时间
    UserID    int       // 用户 ID
    ServerID  int       // 服务器 ID
    AlterID   int64     // Alter ID
    Level     int64     // 等级
    IsLimited bool      // 是否限速
}
```

### SysUser (系统用户)
```go
type SysUser struct {
    ID          uint      // 主键
    UUID        uuid.UUID // UUID
    Username    string    // 用户名
    Password    string    // 密码 (加密)
    NickName    string    // 昵称
    HeaderImg   string    // 头像
    AuthorityId uint      // 角色 ID
    Phone       string    // 手机号
    Email       string    // 邮箱
    Enable      int       // 是否启用
}
```

## Request/Response Format

### 通用响应格式
```json
{
    "code": 0,
    "data": {},
    "msg": "操作成功"
}
```

### 分页请求
```json
{
    "page": 1,
    "pageSize": 10
}
```

### 分页响应
```json
{
    "code": 0,
    "data": {
        "list": [],
        "total": 100,
        "page": 1,
        "pageSize": 10
    },
    "msg": "获取成功"
}
```
