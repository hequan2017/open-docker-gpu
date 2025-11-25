# open-docker-gpu

> docker gpu 开源管理系统 
## 启动

> cd server/ go run main.go

> cd web  /  npm run dev

## 新增功能：Linux SSH 管理

- 模块：Linux SSH 管理
- 字段：
  - 服务器IP地址 `ip`
  - SSH端口 `port`（默认 `22`）
  - 登录账号 `username`（默认 `root`）
  - 登录密码 `password`
  - 状态 `status`（创建/更新后自动检测：正常/异常）
  - 服务器标签名 `label`
  - 服务器地区 `region`
- 行为：创建或更新时，系统使用 `ip/port/username/password` 尝试 SSH 登录，成功则记录状态为“正常”，失败记录为“异常”。
- 接口：
  - `POST /ssh/createSshServer` 新增
  - `DELETE /ssh/deleteSshServer?ID=<id>` 删除
  - `DELETE /ssh/deleteSshServerByIds?IDs[]=<id1>&IDs[]=<id2>` 批量删除
  - `PUT /ssh/updateSshServer` 更新
  - `GET /ssh/findSshServer?ID=<id>` 查询单条
  - `GET /ssh/getSshServerList?page=<n>&pageSize=<m>&ip=<like>&label=<like>&region=<like>&status=正常` 列表
  - `GET /ssh/gpuInfo?ID=<id>` 获取 GPU 信息
  - `GET /ssh/nvidiaSmiText?ID=<id>&tail=<n>` 获取 `nvidia-smi` 文本摘要（尾部 n 行）
- 前端：
  - 菜单：Linux SSH 管理（路由 `ssh`）
  - 列表列：IP（宽度 200）、端口、账号、服务器标签名、服务器地区、状态
  - 表单默认值：`port=22`、`username='root'`
  - 密码输入：编辑窗口显示且字符隐藏（`type=password`）
  - 状态：表单只读显示后端检测结果
  - 交互与查询时机：
    - 首次进入页面：自动加载并仅展示列表基础数据（轻量）
    - 点击“查询”、分页、排序、重置：触发列表刷新（查询）
    - 展开行与详情：点击后按需请求；展开行支持 `nvidia-smi` 文本预览、手动刷新与自动刷新（可设置间隔与尾部行数）
- 示例（创建）：
  - 请求体：`{"ip":"192.168.1.10","port":22,"username":"root","password":"your_password","label":"生产A","region":"杭州"}`

## 新增功能：容器 WebSSH 终端

- 前端：
  - 页面：`Docker 管理`（`view/docker/dockerManage.vue`）
  - 功能：`xterm.js` 终端渲染、复制选中、清屏、Shell 切换（`/bin/sh`/`/bin/bash`）、尺寸预设与自适应、断线自动重连与心跳保活
- 后端：
  - 接口：`GET /docker/execWs`（WebSocket，参数：`endpointId`、`ID`、`shell`）
  - 说明：创建 `exec` 并桥接容器 STDIN/STDOUT；支持 `{"type":"resize","cols":<int>,"rows":<int>}` 同步终端尺寸，`{"type":"ping"}` 心跳
- 相关 API：
  - `GET /docker/servers` 获取正常服务器
  - `GET /docker/ps` 获取容器列表（支持 `scope=running|exited|all`）
  - `POST /docker/createContainer` 创建容器
  - `POST /docker/startContainer` 启动容器
  - `POST /docker/stopContainer` 停止容器
  - `DELETE /docker/removeContainer` 删除容器
  - `GET /docker/findDockerEndpoint` 查询 Endpoint
  - `GET /docker/getDockerEndpointList` Endpoint 列表
  - `POST /docker/createDockerEndpoint` 创建 Endpoint
  - `DELETE /docker/deleteDockerEndpoint` 删除 Endpoint
  - `DELETE /docker/deleteDockerEndpointByIds` 批量删除 Endpoint
  - `PUT /docker/updateDockerEndpoint` 更新 Endpoint

## 初始化数据更新

- API 初始化：在 `server/source/system/api.go` 的 `entities` 中新增上述 Docker 接口（含 `execWs`）
- 菜单初始化：在 `server/source/system/menu.go` 的 `childMenus` 中新增 `Docker 管理` 菜单，父级为 `systemTools`，组件路径 `view/docker/dockerManage.vue`
