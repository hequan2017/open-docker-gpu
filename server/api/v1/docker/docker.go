package docker

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DockerApi struct{}

// GetNormalServers 获取状态为正常的服务器列表
// @Tags Docker
// @Summary 获取状态为正常的服务器列表
// @Security ApiKeyAuth
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string}
// @Router /docker/servers [get]
func (a *DockerApi) GetNormalServers(c *gin.Context) {
	ctx := c.Request.Context()
	list, err := dockerService.ListNormalServers(ctx)
	if err != nil {
		global.GVA_LOG.Error("获取服务器失败!", zap.Error(err))
		response.FailWithMessage("获取服务器失败:"+err.Error(), c)
		return
	}
	response.OkWithData(list, c)
}

// GetDockerPs 获取指定服务器的 docker ps -a 信息
// @Tags Docker
// @Summary 获取指定服务器的 docker ps -a 信息
// @Security ApiKeyAuth
// @Produce application/json
// @Param ID query string true "服务器ID"
// @Success 200 {object} response.Response{data=object,msg=string}
// @Router /docker/ps [get]
func (a *DockerApi) GetDockerPs(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Query("ID")
	if id == "" || id == "undefined" {
		response.FailWithMessage("未选择Endpoint", c)
		return
	}
	scope := c.DefaultQuery("scope", "running")
	rows, err := dockerService.FetchDockerPs(ctx, id, scope)
	if err != nil {
		global.GVA_LOG.Error("获取容器失败!", zap.Error(err))
		response.FailWithMessage("获取容器失败:"+err.Error(), c)
		return
	}
	response.OkWithData(rows, c)
}
