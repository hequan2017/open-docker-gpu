// 自动生成模板SshServer
package ssh

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// Linux SSH管理 结构体  SshServer
type SshServer struct {
	global.GVA_MODEL
	Ip       *string `json:"ip" form:"ip" gorm:"comment:服务器IP;column:ip;size:64;" binding:"required"`
	Port     *int64  `json:"port" form:"port" gorm:"default:22;comment:端口;column:port;"`
	Username *string `json:"username" form:"username" gorm:"default:root;comment:账号;column:username;size:64;"`
	Password *string `json:"password" form:"password" gorm:"comment:密码;column:password;size:128;" binding:"required"`
	Status   *string `json:"status" form:"status" gorm:"comment:状态;column:status;size:16;"`
	Label    *string `json:"label" form:"label" gorm:"comment:服务器标签名;column:label;size:64;"`
	Region   *string `json:"region" form:"region" gorm:"comment:服务器地区;column:region;size:64;"`
}

// TableName Linux SSH管理 SshServer自定义表名 ssh_servers
func (SshServer) TableName() string {
	return "ssh_servers"
}
