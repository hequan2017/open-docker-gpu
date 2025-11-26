package meter

import (
	"context"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	assetModel "github.com/flipped-aurora/gin-vue-admin/server/model/asset"
	billingModel "github.com/flipped-aurora/gin-vue-admin/server/model/billing"
	dockerSvcPkg "github.com/flipped-aurora/gin-vue-admin/server/service/docker"
	sshSvcPkg "github.com/flipped-aurora/gin-vue-admin/server/service/ssh"
)

type GpuMeterService struct{}

// countContainerGPUs 解析容器的 DeviceRequests 信息，返回设备ID集合、count加总以及是否包含ALL
func countContainerGPUs(reqs []container.DeviceRequest) (deviceIDs []string, countSum int, hasAll bool) {
	for _, dr := range reqs {
		// 仅处理包含gpu能力的请求
		hasGPU := false
		for _, capset := range dr.Capabilities {
			for _, cap := range capset {
				if cap == "gpu" {
					hasGPU = true
					break
				}
			}
			if hasGPU {
				break
			}
		}
		if !hasGPU {
			continue
		}
		if dr.Count == -1 {
			hasAll = true
			// ALL 模式无需继续解析具体ID
			continue
		}
		if len(dr.DeviceIDs) > 0 {
			deviceIDs = append(deviceIDs, dr.DeviceIDs...)
		} else if dr.Count > 0 {
			countSum += dr.Count
		}
	}
	return
}

// CollectOneEndpoint 采集单个 Endpoint 的 GPU 使用率
func (s *GpuMeterService) CollectOneEndpoint(ctx context.Context, endpointID string) error {
	var dockerSvc dockerSvcPkg.DockerService
	var sshSvc sshSvcPkg.SshServerService

	// 查资产映射
	var asset assetModel.ServerAsset
	// endpointID 为字符串形式的自增ID，转换为int64以匹配资产表字段类型
	eid := strToInt64(endpointID)
	if err := global.GVA_DB.Where("endpoint_id = ?", eid).First(&asset).Error; err != nil {
		return err
	}

	// 主机总卡数
	var hostTotal int64
	if asset.SshServerId != nil {
		if items, err := sshSvc.FetchGPUInfo(ctx, toStr(*asset.SshServerId)); err == nil {
			hostTotal = int64(len(items))
		}
	}
	if hostTotal == 0 && asset.GpuCount != nil {
		hostTotal = *asset.GpuCount
	}
	if hostTotal <= 0 {
		// 无法确定主机卡数，记录备注并跳过
		now := time.Now().UTC().Truncate(time.Minute)
		aid := int64(asset.ID)
		rec := billingModel.GPUUsageRecord{
			AssetId:    &aid,
			EndpointId: &eid,
			MeasuredAt: &now,
			Note:       strPtr("主机GPU总卡数未知，采集跳过"),
		}
		return global.GVA_DB.Create(&rec).Error
	}

	// 列出运行中容器
	rows, err := dockerSvc.FetchDockerPs(ctx, endpointID, "running")
	if err != nil {
		return err
	}

	cli, err := dockerSvc.GetClientByEndpointID(ctx, endpointID)
	if err != nil {
		return err
	}

	// 聚合设备ID与count
	idSet := make(map[string]struct{})
	var countSum int
	var hasAll bool

	for _, row := range rows {
		ins, e := cli.ContainerInspect(ctx, row.ID)
		if e != nil {
			continue
		}
		devs, cnt, all := countContainerGPUs(ins.HostConfig.DeviceRequests)
		if all {
			hasAll = true
		}
		for _, id := range devs {
			idSet[id] = struct{}{}
		}
		countSum += cnt
	}

	var usedCards int64
	if hasAll {
		usedCards = hostTotal
	} else {
		// 设备ID集合大小与count累加二者取较大者，但不超过主机总卡数
		size := int64(len(idSet))
		sum := int64(countSum)
		if sum > size {
			usedCards = sum
		} else {
			usedCards = size
		}
		if usedCards > hostTotal {
			usedCards = hostTotal
		}
	}

	// 使用率
	var rate float64
	if hostTotal > 0 {
		rate = float64(usedCards) / float64(hostTotal)
	}

	// 写入记录
	now := time.Now().UTC().Truncate(time.Minute)
	cc := int64(len(rows))
	aid := int64(asset.ID)
	rec := billingModel.GPUUsageRecord{
		AssetId:        &aid,
		EndpointId:     &eid,
		MeasuredAt:     &now,
		ContainerCount: &cc,
		HostGpuTotal:   &hostTotal,
		UsedGpuCards:   &usedCards,
		UsageRate:      &rate,
	}
	return global.GVA_DB.Create(&rec).Error
}

// CollectAllEndpoints 遍历所有状态正常的 DockerEndpoint 进行采集
func (s *GpuMeterService) CollectAllEndpoints(ctx context.Context) error {
	var dockerSvc dockerSvcPkg.DockerService
	eps, err := dockerSvc.ListNormalServers(ctx)
	if err != nil {
		return err
	}
	for _, ep := range eps {
		id := fmtInt(int64(ep.ID))
		if id == "" {
			continue
		}
		_ = s.CollectOneEndpoint(ctx, id)
	}
	return nil
}

// 工具函数
func strToInt64(s string) int64 {
	var v int64
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return 0
		}
		v = v*10 + int64(c-'0')
	}
	return v
}

func strPtr(s string) *string { return &s }

func toStr(i int64) string {
	// 将整型ID转字符串
	if i == 0 {
		return ""
	}
	return fmtInt(i)
}

func fmtInt(v int64) string {
	// 轻量字符串格式化，避免引入fmt
	if v == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for v > 0 {
		i--
		buf[i] = byte('0' + v%10)
		v /= 10
	}
	return string(buf[i:])
}
