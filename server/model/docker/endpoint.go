package docker

import (
    "github.com/flipped-aurora/gin-vue-admin/server/global"
)

// DockerEndpoint 远程 Docker SDK 连接配置
type DockerEndpoint struct {
    global.GVA_MODEL
    Label      *string `json:"label" form:"label" gorm:"comment:服务器标签名;size:64"`
    Endpoint   *string `json:"endpoint" form:"endpoint" gorm:"comment:连接地址;size:256"`
    UseTLS     *bool   `json:"useTls" form:"useTls" gorm:"comment:是否使用TLS"`
    CACert     *string `json:"caCert" form:"caCert" gorm:"comment:CA证书;type:text"`
    ClientCert *string `json:"clientCert" form:"clientCert" gorm:"comment:客户端证书;type:text"`
    ClientKey  *string `json:"clientKey" form:"clientKey" gorm:"comment:客户端私钥;type:text"`
    Status     *string `json:"status" form:"status" gorm:"comment:连接状态;size:16"`
}

// TableName 自定义表名
func (DockerEndpoint) TableName() string { return "docker_endpoints" }

