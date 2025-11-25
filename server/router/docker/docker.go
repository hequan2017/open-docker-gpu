package docker

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type DockerRouter struct{}

// InitDockerRouter 初始化 Docker 服务管理 路由
func (r *DockerRouter) InitDockerRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	// 读接口（需鉴权）
	groupWithoutRecord := Router.Group("docker")
	{
		groupWithoutRecord.GET("servers", dockerApi.GetNormalServers)
		groupWithoutRecord.GET("ps", dockerApi.GetDockerPs)
		groupWithoutRecord.GET("preferredShell", dockerApi.PreferredShell)
		groupWithoutRecord.GET("logs", dockerApi.GetContainerLogs)
		groupWithoutRecord.GET("findDockerEndpoint", dockerEndpointApi.FindDockerEndpoint)
		groupWithoutRecord.GET("getDockerEndpointList", dockerEndpointApi.GetDockerEndpointList)
	}
	// 读接口（不鉴权）：WebSocket 终端
	groupPublic := PublicRouter.Group("docker")
	{
		groupPublic.GET("execWs", dockerApi.ContainerTerminalWs)
		groupPublic.GET("logsWs", dockerApi.ContainerLogsWs)
	}
	// 写接口
	group := Router.Group("docker").Use(middleware.OperationRecord())
	{
		group.POST("createDockerEndpoint", dockerEndpointApi.CreateDockerEndpoint)
		group.DELETE("deleteDockerEndpoint", dockerEndpointApi.DeleteDockerEndpoint)
		group.DELETE("deleteDockerEndpointByIds", dockerEndpointApi.DeleteDockerEndpointByIds)
		group.PUT("updateDockerEndpoint", dockerEndpointApi.UpdateDockerEndpoint)
		group.POST("createContainer", dockerApi.CreateContainer)
		group.POST("createContainerByDockerfile", dockerApi.BuildContainerByDockerfile)
		group.POST("createContainerWithOptions", dockerApi.CreateContainerWithOptions)
		group.POST("startContainer", dockerApi.StartContainer)
		group.POST("stopContainer", dockerApi.StopContainer)
		group.DELETE("removeContainer", dockerApi.RemoveContainer)
	}
}
