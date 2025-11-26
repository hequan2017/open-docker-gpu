## 目标与度量
- 每分钟采集：每台 Server 的“容器数量”和“容器请求的 GPU 卡数”，再除以该机器的总 GPU 卡数，得到使用率并保存。
- 按天/按小时汇总，为计费提供基础数据与查询接口。

## 数据来源
- 容器列表与配置：复用已有 Docker 管理能力，列容器与检查容器 HostConfig（`DeviceRequests`）。
  - 容器列表：`service/docker/docker.go:47` 的 `FetchDockerPs`。
  - Docker 客户端：`service/docker/docker.go:106` 的 `getClientByEndpointID`。
- 主机 GPU 数量：复用 SSH GPU 信息采集。
  - `service/ssh/sshServer.go:245` 的 `FetchGPUInfo` 返回 `[]GPUInfo`，长度即主机卡数。
- Server 与 Endpoint 关联：复用资产表映射。
  - `model/asset/serverAsset.go:24-25` 的 `SshServerId/EndpointId` 字段。

## 使用率计算规则
- 对“机器总卡数”：优先使用 `FetchGPUInfo(SshServerId)` 的 GPU 数量；失败时回退到 `ServerAsset.GpuCount` 字段；还失败则标记为未知并跳过计费。
- 对“容器用了多少卡”：对每个运行中容器读取 `HostConfig.DeviceRequests`（Docker SDK `ContainerInspect`）：
  - `Capabilities` 含 `"gpu"` 才计入。
  - 若 `Count = -1`（all）：视为使用“主机总卡数”。
  - 若 `DeviceIDs` 非空：以 `DeviceIDs` 的长度计数；做全局去重以避免重复 GPU 号计算（同一 GPU 可能被多容器共享，按“分配卡数”口径仍计 1）。
  - 若仅 `Count = n`：按 `n` 计数（无具体 ID 时无法去重，最终对所有容器的 count 求和并 `cap` 到主机总卡数）。
- 计算公式：`usage_rate = used_cards / total_cards`，四舍五入保留 4 位小数。
- 容器数量：取运行中容器数（`scope=running`）。

## 数据模型
- 新增表 `gpu_usage_records`：
  - `id`、`created_at`（GVA 模型）
  - `asset_id`（关联 `server_assets.id`）
  - `endpoint_id`（关联 `docker_endpoints.id`）
  - `measured_at`（时间戳，精确到分钟）
  - `container_count`（整数）
  - `host_gpu_total`（整数）
  - `used_gpu_cards`（整数）
  - `usage_rate`（小数，0~1）
  - `note`（文本，记录异常或回退情况）
- 索引：`(endpoint_id, measured_at)`、`(asset_id, measured_at)`。

## 采集服务
- 新增 `service/meter/gpu.go`（示例命名）实现：
  - `CollectOneEndpoint(endpointID string) error`
    1. 通过 `endpointID` 找到 `ServerAsset`（`EndpointId` 匹配），取出 `SshServerId` 与主机信息。
    2. 计算主机总卡数：`ssh.FetchGPUInfo(ctx, SshServerId)` → `len(items)`；失败则回退至 `ServerAsset.GpuCount`。
    3. 列出运行中容器：`dockerService.FetchDockerPs(ctx, endpointID, "running")`（`service/docker/docker.go:47`）。
    4. 针对每个容器：使用 `cli.ContainerInspect` 解析 `HostConfig.DeviceRequests`（与创建时的 `CreateContainerWithOptions` 呼应，参考 `service/docker/docker.go:285-295`）。
    5. 聚合得到 `used_gpu_cards` 与 `container_count`，计算 `usage_rate`。
    6. 写入 `gpu_usage_records`。
  - `CollectAllEndpoints() error`：遍历 `ListNormalServers(ctx)`（`service/docker/docker.go:27-31`），逐个调用 `CollectOneEndpoint`，对失败项写入 `note` 并继续。

## 定时任务
- 在 `initialize/timer.go` 注册每分钟任务：
  - 使用 `global.GVA_Timer.AddTaskByFuncWithSecond("GpuBilling", "0 */1 * * * *", func() { meterService.CollectAllEndpoints() }, "GPU使用率采集", option...)`（`utils/timer/timed_task.go:75-95`）。
  - 任务确保入库幂等：`measured_at` 以分钟精度作为唯一键（同一 endpoint 重复采集覆盖或忽略）。

## API 接口
- 读接口（仅查询）：
  - `GET /billing/gpu/latest?endpointId=...`：返回最近一次记录。
  - `GET /billing/gpu/records?endpointId=...&start=...&end=...&page=...`：分页返回区间内记录。
  - `GET /billing/gpu/summary?endpointId=...&groupBy=hour|day`：返回聚合（容器均值、使用率均值/最大值、卡数均值）。
- 写接口（可选）：
  - `POST /billing/gpu/collect?endpointId=...`：手动触发单次采集（便于验证）。

## 计费与汇总
- 计费参数：使用系统参数表存储（如 `sys_params`）
  - `gpu.rate.perCard.perMinute`（卡·分钟单价）或 `gpu.rate.perUsagePercent.perMinute`（使用率·分钟单价）。
- 计费口径建议：
  - 按“分配卡数”计费：`cost_minute = used_gpu_cards * rate_card_min`。
  - 或按“使用率”计费：`cost_minute = usage_rate * rate_usage_min`。
- 汇总服务：
  - 按小时/天聚合分钟记录，生成账单快照（可扩展导出 CSV/Excel）。

## 错误处理与监控
- 可达性与异常：
  - Docker Endpoint ping 检查已内置：`service/docker/endpoint.go: checkAndSetStatus`。
  - SSH 采集失败时回退 `ServerAsset.GpuCount` 并记 `note`。
- 数据一致性：
  - `DeviceIDs` 去重；`Count` 求和后 `cap` 到主机总卡数；遇到 `all` 直接取主机总卡数。
  - 同一分钟重复写入采用“覆盖”策略或唯一索引避免脏数据。

## 验证方案
- 手动触发 `POST /billing/gpu/collect`，对比返回的容器数量与主机 GPU 数量（`GET /docker/ps` 与 `GET /ssh/gpuInfo`）。
- 在一台测试机上创建容器分别使用 `gpuMode=devices/count/all`（参考创建接口：`api/v1/docker/container.go:78-107`），验证解析与计数正确性。
- 对边界情况测试：0 容器、无 GPU 主机、多个容器共享同一 `DeviceID`、包含 `all` 的容器。

## 变更点列表（代码路径）
- 新增：`model/billing/gpu_usage_record.go`（表结构）。
- 新增：`service/meter/gpu.go`（采集与计算）。
- 变更：`initialize/timer.go`（注册每分钟采集任务）。
- 新增：`api/v1/billing/gpu.go` 与路由 `router/billing/gpu.go`（查询/手动采集接口）。

## 说明与假设
- “用了多少卡”按“容器分配的 GPU 卡”口径计算，非实时 GPU 利用率；如需真实负载可后续接入 `nvidia-smi pmon/DCGM` 做进阶计费。
- 多容器共享同一 GPU 的情况，按卡数口径仍计为 1 张卡。
- 如某 Endpoint 未在 `server_assets` 建立映射，采集时仅能依据 `Count`/`DeviceIDs`估算，主机总卡数不可得将跳过该 endpoint 并记录异常。