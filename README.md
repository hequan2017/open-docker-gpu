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

## 新增功能：GPU 使用率计费系统

- 目标：每分钟采集每台服务器的“运行中容器数量”和“容器分配的 GPU 卡数”，除以该机器的总 GPU 卡数得到使用率并保存。
- 数据来源：
  - 容器与配置：通过 Docker SDK 读取容器列表与 `HostConfig.DeviceRequests`（`gpu` 能力）。
  - 主机 GPU 总卡数：优先通过 SSH 执行 `nvidia-smi/rocm-smi` 获取；失败回退资产表 `server_assets.gpu_count`。
- 定时任务：每分钟自动采集任务 `GpuBilling` 已在后端注册（`initialize/timer.go`）。
- 数据表：`gpu_usage_records`
  - 字段：`asset_id`、`endpoint_id`、`measured_at`（分钟精度）、`container_count`、`host_gpu_total`、`used_gpu_cards`、`usage_rate`、`note`
  - 唯一索引：`(endpoint_id, measured_at)` 保证同一分钟幂等。
- 接口：
  - `POST /billing/gpu/collect?endpointId=<id>` 手动触发一次采集
  - `GET /billing/gpu/latest?endpointId=<id>` 获取最近一次记录
  - `GET /billing/gpu/records?endpointId=<id>&start=<RFC3339>&end=<RFC3339>&page=<n>&pageSize=<m>` 分页查询区间记录
  - `GET /billing/gpu/summary?endpointId=<id>&groupBy=hour|day&start=<RFC3339>&end=<RFC3339>` 聚合查询（返回每桶平均使用率、平均分配卡数、容器均值与分钟均成本）
- 计费参数（系统参数表 `sys_params`）：
  - `gpu.rate.perCard.perMinute` 卡·分钟单价（按分配卡数计费）
  - `gpu.rate.perUsagePercent.perMinute` 使用率·分钟单价（按使用率计费）
  - 聚合接口会读取并计算 `costMinuteAvg`（设置任一非零参数即可生效）。
- 计算规则摘要：
  - `DeviceRequests` 中包含 `gpu` 能力才计入；若 `Count=-1`（all），视为使用主机总卡数。
  - 若 `DeviceIDs` 非空按唯一设备 ID 个数计；仅 `Count=n` 则按 `n` 计数。
  - 总分配卡数取两者较大者，并不超过主机总卡数。
- 验证示例：
  - 手动采集：`POST /billing/gpu/collect?endpointId=1`
  - 最新记录：`GET /billing/gpu/latest?endpointId=1`
  - 分页查询：`GET /billing/gpu/records?endpointId=1&page=1&pageSize=20`
  - 小时聚合：`GET /billing/gpu/summary?endpointId=1&groupBy=hour`



### 设计
#### 镜像管理
镜像名字
镜像地址

#### 服务器管理

名字
区域
CPU
内存
系统盘容量
数据盘容量
ip地址公网
ip地址内网
端口
用户名
密码
显卡名称   
显卡数量
docker连接地址
使用TLS   
CA证书  text
客户端证书 text
客户端私钥 text

### 实例管理

实例模版

GPU型号
GPU数量
CPU核心数
内存
系统盘容量
数据盘容器
价格/小时



购买实例

选择镜像
选择显卡
选择显卡数量
根据显卡数量，选择实例模板
展示满足显卡，显卡数量和实例模板(CPU 内存  系统盘  数据盘  大小)可以购买的实例模版


创建实例，后端去对应的服务器创建容器。 每个用户只能看到自己创建的实例。admin 管理员可以看到所有

实例管理
实例创建的用户ID
实例名字
实例关联 docker名字
状态 （运行中/已停止/已删除）
显卡名称
显卡数量
CPU
内存
操作 （启动/停止/删除）







