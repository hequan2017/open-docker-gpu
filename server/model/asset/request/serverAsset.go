
package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type ServerAssetSearch struct{
    CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
      Name  *string `json:"name" form:"name"` 
      Hostname  *string `json:"hostname" form:"hostname"` 
      Ip  *string `json:"ip" form:"ip"` 
      Port  *int `json:"port" form:"port"` 
      OsType  *string `json:"osType" form:"osType"` 
      CpuCores  *int `json:"cpuCores" form:"cpuCores"` 
      MemoryGB  *int `json:"memoryGB" form:"memoryGB"` 
      DiskGB  *int `json:"diskGB" form:"diskGB"` 
      GpuCount  *int `json:"gpuCount" form:"gpuCount"` 
      Status  *string `json:"status" form:"status"` 
      Region  *string `json:"region" form:"region"` 
      Label  *string `json:"label" form:"label"` 
      Remark  string `json:"remark" form:"remark"` 
      SshServerId  *int `json:"sshServerId" form:"sshServerId"` 
      EndpointId  *int `json:"endpointId" form:"endpointId"` 
    request.PageInfo
    Sort  string `json:"sort" form:"sort"`
    Order string `json:"order" form:"order"`
}
