# web

`web` 是 `xray-admin` 的 Vue 前端。生产发布时需要先编译前端，再把 `dist` 内容复制到后端资源目录。

## 安装依赖

```bash
cd gin-vue-admin/web
npm install --legacy-peer-deps
```

如果在 Windows PowerShell 下遇到 `npm.ps1` 执行策略限制，可以改用：

```powershell
npm.cmd install --legacy-peer-deps
```

## 本地开发

```bash
npm run serve
```

开发环境接口代理配置在 `.env.development`。

## 生产构建

```bash
npm run build
```

构建结果会生成在：

```text
gin-vue-admin/web/dist
```

## 复制到后端资源目录

当前后端读取的不是 `server/dist`，而是：

```text
gin-vue-admin/server/resource/page
```

构建完成后执行：

```bash
rm -rf ../server/resource/page
mkdir -p ../server/resource/page
cp -r dist/* ../server/resource/page/
```

复制后应存在：

```text
gin-vue-admin/server/resource/page/index.html
gin-vue-admin/server/resource/page/assets/
```

`server/resource/page` 是构建产物目录，已被 `.gitignore` 忽略；`server/resource/resource.go` 是后端源码，需要保留。
