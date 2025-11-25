package ssh

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type SshServerRouter struct{}

// InitSshServerRouter 初始化 Linux SSH管理 路由信息
func (s *SshServerRouter) InitSshServerRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	sshRouter := Router.Group("ssh").Use(middleware.OperationRecord())
	sshRouterWithoutRecord := Router.Group("ssh")
	sshRouterWithoutAuth := PublicRouter.Group("ssh")
	{
		sshRouter.POST("createSshServer", sshApi.CreateSshServer)             // 新建Linux SSH管理
		sshRouter.DELETE("deleteSshServer", sshApi.DeleteSshServer)           // 删除Linux SSH管理
		sshRouter.DELETE("deleteSshServerByIds", sshApi.DeleteSshServerByIds) // 批量删除Linux SSH管理
		sshRouter.PUT("updateSshServer", sshApi.UpdateSshServer)              // 更新Linux SSH管理
	}
	{
		sshRouterWithoutRecord.GET("findSshServer", sshApi.FindSshServer)       // 根据ID获取Linux SSH管理
		sshRouterWithoutRecord.GET("getSshServerList", sshApi.GetSshServerList) // 获取Linux SSH管理列表
		sshRouterWithoutRecord.GET("gpuInfo", sshApi.GetGpuInfo)                // 获取GPU信息
		sshRouterWithoutRecord.GET("nvidiaSmiText", sshApi.NvidiaSmiText)       // 获取nvidia-smi文本摘要
	}
	{
		sshRouterWithoutAuth.GET("getSshServerPublic", sshApi.GetSshServerPublic) // Linux SSH管理开放接口
	}
}
