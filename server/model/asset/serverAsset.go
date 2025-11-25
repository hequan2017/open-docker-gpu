
// 自动生成模板ServerAsset
package asset
import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 服务器资产表 结构体  ServerAsset
type ServerAsset struct {
    global.GVA_MODEL
  Name  *string `json:"name" form:"name" gorm:"comment:服务器名称;column:name;size:64;" binding:"required"`  //服务器名称
  Hostname  *string `json:"hostname" form:"hostname" gorm:"comment:主机名;column:hostname;size:128;"`  //主机名
  Ip  *string `json:"ip" form:"ip" gorm:"comment:IP地址;column:ip;size:64;" binding:"required"`  //IP地址
  Port  *int64 `json:"port" form:"port" gorm:"default:22;comment:端口;column:port;"`  //端口
  OsType  *string `json:"osType" form:"osType" gorm:"comment:操作系统类型;column:os_type;size:32;"`  //操作系统类型
  CpuCores  *int64 `json:"cpuCores" form:"cpuCores" gorm:"comment:CPU核数;column:cpu_cores;"`  //CPU核数
  MemoryGB  *int64 `json:"memoryGB" form:"memoryGB" gorm:"comment:内存容量GB;column:memory_gb;"`  //内存容量GB
  DiskGB  *int64 `json:"diskGB" form:"diskGB" gorm:"comment:磁盘容量GB;column:disk_gb;"`  //磁盘容量GB
  GpuCount  *int64 `json:"gpuCount" form:"gpuCount" gorm:"comment:GPU数量;column:gpu_count;"`  //GPU数量
  Status  *string `json:"status" form:"status" gorm:"comment:资产状态;column:status;size:16;"`  //资产状态
  Region  *string `json:"region" form:"region" gorm:"comment:地区;column:region;size:64;"`  //地区
  Label  *string `json:"label" form:"label" gorm:"comment:标签;column:label;size:64;"`  //标签
  Remark  *string `json:"remark" form:"remark" gorm:"comment:备注;column:remark;type:text;"`  //备注
  SshServerId  *int64 `json:"sshServerId" form:"sshServerId" gorm:"comment:ssh_servers关联ID;column:ssh_server_id;"`  //关联SSH服务器
  EndpointId  *int64 `json:"endpointId" form:"endpointId" gorm:"comment:docker_endpoints关联ID;column:endpoint_id;"`  //关联Docker端点
}


// TableName 服务器资产表 ServerAsset自定义表名 server_assets
func (ServerAsset) TableName() string {
    return "server_assets"
}





