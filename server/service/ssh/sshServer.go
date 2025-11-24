package ssh

import (
	"context"
	"fmt"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ssh"
	sshReq "github.com/flipped-aurora/gin-vue-admin/server/model/ssh/request"
	sshc "golang.org/x/crypto/ssh"
)

type SshServerService struct{}

// CreateSshServer 创建Linux SSH管理记录
// Author [yourname](https://github.com/yourname)
func (sshService *SshServerService) CreateSshServer(ctx context.Context, record *ssh.SshServer) (err error) {
	if record.Port == nil {
		var p int64 = 22
		record.Port = &p
	}
	if record.Username == nil || *record.Username == "" {
		u := "root"
		record.Username = &u
	}
	if record.Ip != nil && record.Password != nil && record.Port != nil && record.Username != nil {
		addr := fmt.Sprintf("%s:%d", *record.Ip, int(*record.Port))
		cfg := &sshc.ClientConfig{User: *record.Username, Auth: []sshc.AuthMethod{sshc.Password(*record.Password)}, HostKeyCallback: sshc.InsecureIgnoreHostKey(), Timeout: 5 * time.Second}
		if c, e := sshc.Dial("tcp", addr, cfg); e == nil {
			s := "正常"
			record.Status = &s
			_ = c.Close()
		} else {
			s := "异常"
			record.Status = &s
		}
	}
	err = global.GVA_DB.Create(record).Error
	return err
}

// DeleteSshServer 删除Linux SSH管理记录
// Author [yourname](https://github.com/yourname)
func (sshService *SshServerService) DeleteSshServer(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&ssh.SshServer{}, "id = ?", ID).Error
	return err
}

// DeleteSshServerByIds 批量删除Linux SSH管理记录
// Author [yourname](https://github.com/yourname)
func (sshService *SshServerService) DeleteSshServerByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]ssh.SshServer{}, "id in ?", IDs).Error
	return err
}

// UpdateSshServer 更新Linux SSH管理记录
// Author [yourname](https://github.com/yourname)
func (sshService *SshServerService) UpdateSshServer(ctx context.Context, record *ssh.SshServer) (err error) {
	if record.Port == nil {
		var p int64 = 22
		record.Port = &p
	}
	if record.Username == nil || *record.Username == "" {
		u := "root"
		record.Username = &u
	}
	if record.Ip != nil && record.Password != nil && record.Port != nil && record.Username != nil {
		addr := fmt.Sprintf("%s:%d", *record.Ip, int(*record.Port))
		cfg := &sshc.ClientConfig{User: *record.Username, Auth: []sshc.AuthMethod{sshc.Password(*record.Password)}, HostKeyCallback: sshc.InsecureIgnoreHostKey(), Timeout: 5 * time.Second}
		if c, e := sshc.Dial("tcp", addr, cfg); e == nil {
			s := "正常"
			record.Status = &s
			_ = c.Close()
		} else {
			s := "异常"
			record.Status = &s
		}
	}
	err = global.GVA_DB.Model(&ssh.SshServer{}).Where("id = ?", record.ID).Updates(record).Error
	return err
}

// GetSshServer 根据ID获取Linux SSH管理记录
// Author [yourname](https://github.com/yourname)
func (sshService *SshServerService) GetSshServer(ctx context.Context, ID string) (record ssh.SshServer, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&record).Error
	return record, err
}

// GetSshServerInfoList 分页获取Linux SSH管理记录
// Author [yourname](https://github.com/yourname)
func (sshService *SshServerService) GetSshServerInfoList(ctx context.Context, info sshReq.SshServerSearch) (list []ssh.SshServer, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&ssh.SshServer{})
	var sshs []ssh.SshServer
	// 如果有条件搜索 下方会自动创建搜索语句
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	if info.Ip != nil && *info.Ip != "" {
		db = db.Where("ip LIKE ?", "%"+*info.Ip+"%")
	}
	if info.Port != nil {
		db = db.Where("port = ?", *info.Port)
	}
	if info.Username != nil && *info.Username != "" {
		db = db.Where("username LIKE ?", "%"+*info.Username+"%")
	}
	if info.Password != nil && *info.Password != "" {
		db = db.Where("password = ?", *info.Password)
	}
	if info.Status != nil && *info.Status != "" {
		db = db.Where("status = ?", *info.Status)
	}
	if info.Label != nil && *info.Label != "" {
		db = db.Where("label LIKE ?", "%"+*info.Label+"%")
	}
	if info.Region != nil && *info.Region != "" {
		db = db.Where("region LIKE ?", "%"+*info.Region+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	var OrderStr string
	orderMap := make(map[string]bool)
	orderMap["id"] = true
	orderMap["created_at"] = true
	orderMap["ip"] = true
	orderMap["status"] = true
	orderMap["label"] = true
	orderMap["region"] = true
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

	err = db.Find(&sshs).Error
	return sshs, total, err
}
func (sshService *SshServerService) GetSshServerPublic(ctx context.Context) {
	// 此方法为获取数据源定义的数据
	// 请自行实现
}
