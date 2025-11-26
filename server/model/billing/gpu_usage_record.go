package billing

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// GPUUsageRecord 每分钟GPU使用率采集记录
type GPUUsageRecord struct {
	global.GVA_MODEL
	AssetId        *int64     `json:"assetId" form:"assetId" gorm:"comment:server_assets关联ID;column:asset_id;"`
	EndpointId     *int64     `json:"endpointId" form:"endpointId" gorm:"comment:docker_endpoints关联ID;column:endpoint_id;uniqueIndex:ux_endpoint_time"`
	MeasuredAt     *time.Time `json:"measuredAt" form:"measuredAt" gorm:"comment:采集时间(分钟精度);column:measured_at;uniqueIndex:ux_endpoint_time"`
	ContainerCount *int64     `json:"containerCount" form:"containerCount" gorm:"comment:运行中容器数量;column:container_count;"`
	HostGpuTotal   *int64     `json:"hostGpuTotal" form:"hostGpuTotal" gorm:"comment:主机GPU总卡数;column:host_gpu_total;"`
	UsedGpuCards   *int64     `json:"usedGpuCards" form:"usedGpuCards" gorm:"comment:容器分配GPU卡数;column:used_gpu_cards;"`
	UsageRate      *float64   `json:"usageRate" form:"usageRate" gorm:"comment:使用率(0~1);column:usage_rate;"`
	Note           *string    `json:"note" form:"note" gorm:"comment:备注/异常说明;column:note;type:text;"`
}

// TableName 自定义表名
func (GPUUsageRecord) TableName() string { return "gpu_usage_records" }
