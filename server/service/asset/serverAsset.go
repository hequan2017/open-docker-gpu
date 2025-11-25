
package asset

import (
	"context"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/asset"
    assetReq "github.com/flipped-aurora/gin-vue-admin/server/model/asset/request"
)

type ServerAssetService struct {}
// CreateServerAsset 创建服务器资产表记录
// Author [yourname](https://github.com/yourname)
func (SrvAstService *ServerAssetService) CreateServerAsset(ctx context.Context, SrvAst *asset.ServerAsset) (err error) {
	err = global.GVA_DB.Create(SrvAst).Error
	return err
}

// DeleteServerAsset 删除服务器资产表记录
// Author [yourname](https://github.com/yourname)
func (SrvAstService *ServerAssetService)DeleteServerAsset(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&asset.ServerAsset{},"id = ?",ID).Error
	return err
}

// DeleteServerAssetByIds 批量删除服务器资产表记录
// Author [yourname](https://github.com/yourname)
func (SrvAstService *ServerAssetService)DeleteServerAssetByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]asset.ServerAsset{},"id in ?",IDs).Error
	return err
}

// UpdateServerAsset 更新服务器资产表记录
// Author [yourname](https://github.com/yourname)
func (SrvAstService *ServerAssetService)UpdateServerAsset(ctx context.Context, SrvAst asset.ServerAsset) (err error) {
	err = global.GVA_DB.Model(&asset.ServerAsset{}).Where("id = ?",SrvAst.ID).Updates(&SrvAst).Error
	return err
}

// GetServerAsset 根据ID获取服务器资产表记录
// Author [yourname](https://github.com/yourname)
func (SrvAstService *ServerAssetService)GetServerAsset(ctx context.Context, ID string) (SrvAst asset.ServerAsset, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&SrvAst).Error
	return
}
// GetServerAssetInfoList 分页获取服务器资产表记录
// Author [yourname](https://github.com/yourname)
func (SrvAstService *ServerAssetService)GetServerAssetInfoList(ctx context.Context, info assetReq.ServerAssetSearch) (list []asset.ServerAsset, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&asset.ServerAsset{})
    var SrvAsts []asset.ServerAsset
    // 如果有条件搜索 下方会自动创建搜索语句
    if len(info.CreatedAtRange) == 2 {
     db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
    }
    
    if info.Name != nil && *info.Name != "" {
        db = db.Where("name LIKE ?", "%"+ *info.Name+"%")
    }
    if info.Hostname != nil && *info.Hostname != "" {
        db = db.Where("hostname LIKE ?", "%"+ *info.Hostname+"%")
    }
    if info.Ip != nil && *info.Ip != "" {
        db = db.Where("ip LIKE ?", "%"+ *info.Ip+"%")
    }
    if info.Port != nil {
        db = db.Where("port = ?", *info.Port)
    }
    if info.OsType != nil && *info.OsType != "" {
        db = db.Where("os_type = ?", *info.OsType)
    }
    if info.CpuCores != nil {
        db = db.Where("cpu_cores >= ?", *info.CpuCores)
    }
    if info.MemoryGB != nil {
        db = db.Where("memory_gb >= ?", *info.MemoryGB)
    }
    if info.DiskGB != nil {
        db = db.Where("disk_gb >= ?", *info.DiskGB)
    }
    if info.GpuCount != nil {
        db = db.Where("gpu_count >= ?", *info.GpuCount)
    }
    if info.Status != nil && *info.Status != "" {
        db = db.Where("status = ?", *info.Status)
    }
    if info.Region != nil && *info.Region != "" {
        db = db.Where("region LIKE ?", "%"+ *info.Region+"%")
    }
    if info.Label != nil && *info.Label != "" {
        db = db.Where("label LIKE ?", "%"+ *info.Label+"%")
    }
    if info.Remark != "" {
        // TODO 数据类型为复杂类型，请根据业务需求自行实现复杂类型的查询业务
    }
    if info.SshServerId != nil {
        db = db.Where("ssh_server_id = ?", *info.SshServerId)
    }
    if info.EndpointId != nil {
        db = db.Where("endpoint_id = ?", *info.EndpointId)
    }
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }
        var OrderStr string
        orderMap := make(map[string]bool)
           orderMap["id"] = true
           orderMap["created_at"] = true
         	orderMap["name"] = true
         	orderMap["cpu_cores"] = true
         	orderMap["memory_gb"] = true
         	orderMap["disk_gb"] = true
         	orderMap["gpu_count"] = true
         	orderMap["status"] = true
         	orderMap["region"] = true
         	orderMap["label"] = true
       if orderMap[info.Sort] {
          OrderStr = info.Sort
          if info.Order == "descending" {
             OrderStr = OrderStr + " desc"
          }
          db = db.Order(OrderStr)
       }

	if limit != 0 {
       db = db.Limit(limit).Offset(offset)
    }

	err = db.Find(&SrvAsts).Error
	return  SrvAsts, total, err
}
func (SrvAstService *ServerAssetService)GetServerAssetDataSource(ctx context.Context) (res map[string][]map[string]any, err error) {
	res = make(map[string][]map[string]any)
	
	   endpointId := make([]map[string]any, 0)
	   
       
       global.GVA_DB.Table("docker_endpoints").Where("deleted_at IS NULL").Select("label as label,id as value").Scan(&endpointId)
	   res["endpointId"] = endpointId
	   sshServerId := make([]map[string]any, 0)
	   
       
       global.GVA_DB.Table("ssh_servers").Where("deleted_at IS NULL").Select("ip as label,id as value").Scan(&sshServerId)
	   res["sshServerId"] = sshServerId
	return
}
func (SrvAstService *ServerAssetService)GetServerAssetPublic(ctx context.Context) {
    // 此方法为获取数据源定义的数据
    // 请自行实现
}
