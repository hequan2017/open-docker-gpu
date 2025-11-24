package docker

import (
    "github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    model "github.com/flipped-aurora/gin-vue-admin/server/model/docker"
    req "github.com/flipped-aurora/gin-vue-admin/server/model/docker/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type DockerEndpointApi struct{}

// Create
// @Tags DockerEndpoint
// @Summary 创建Docker SDK配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.DockerEndpoint true "Docker SDK配置"
// @Router /docker/createDockerEndpoint [post]
func (a *DockerEndpointApi) CreateDockerEndpoint(c *gin.Context) {
    ctx := c.Request.Context()
    var m model.DockerEndpoint
    if err := c.ShouldBindJSON(&m); err != nil { response.FailWithMessage(err.Error(), c); return }
    if err := dockerEndpointService.CreateDockerEndpoint(ctx, &m); err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
        response.FailWithMessage("创建失败:"+err.Error(), c)
        return
    }
    response.OkWithDetailed(gin.H{"status": m.Status}, "创建成功", c)
}

// Delete
// @Tags DockerEndpoint
// @Summary 删除Docker SDK配置
// @Security ApiKeyAuth
// @Router /docker/deleteDockerEndpoint [delete]
func (a *DockerEndpointApi) DeleteDockerEndpoint(c *gin.Context) {
    ctx := c.Request.Context()
    id := c.Query("ID")
    if err := dockerEndpointService.DeleteDockerEndpoint(ctx, id); err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
        response.FailWithMessage("删除失败:"+err.Error(), c)
        return
    }
    response.OkWithMessage("删除成功", c)
}

// DeleteByIds
// @Tags DockerEndpoint
// @Summary 批量删除Docker SDK配置
// @Security ApiKeyAuth
// @Router /docker/deleteDockerEndpointByIds [delete]
func (a *DockerEndpointApi) DeleteDockerEndpointByIds(c *gin.Context) {
    ctx := c.Request.Context()
    IDs := c.QueryArray("IDs[]")
    if err := dockerEndpointService.DeleteDockerEndpointByIds(ctx, IDs); err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
        response.FailWithMessage("批量删除失败:"+err.Error(), c)
        return
    }
    response.OkWithMessage("批量删除成功", c)
}

// Update
// @Tags DockerEndpoint
// @Summary 更新Docker SDK配置
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.DockerEndpoint true "Docker SDK配置"
// @Router /docker/updateDockerEndpoint [put]
func (a *DockerEndpointApi) UpdateDockerEndpoint(c *gin.Context) {
    ctx := c.Request.Context()
    var m model.DockerEndpoint
    if err := c.ShouldBindJSON(&m); err != nil { response.FailWithMessage(err.Error(), c); return }
    if err := dockerEndpointService.UpdateDockerEndpoint(ctx, &m); err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
        response.FailWithMessage("更新失败:"+err.Error(), c)
        return
    }
    response.OkWithDetailed(gin.H{"status": m.Status}, "更新成功", c)
}

// Find
// @Tags DockerEndpoint
// @Summary 用id查询Docker SDK配置
// @Security ApiKeyAuth
// @Router /docker/findDockerEndpoint [get]
func (a *DockerEndpointApi) FindDockerEndpoint(c *gin.Context) {
    ctx := c.Request.Context()
    id := c.Query("ID")
    m, err := dockerEndpointService.GetDockerEndpoint(ctx, id)
    if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
        response.FailWithMessage("查询失败:"+err.Error(), c)
        return
    }
    response.OkWithData(m, c)
}

// List
// @Tags DockerEndpoint
// @Summary 分页获取Docker SDK配置
// @Security ApiKeyAuth
// @Param data query req.DockerEndpointSearch true "分页条件"
// @Router /docker/getDockerEndpointList [get]
func (a *DockerEndpointApi) GetDockerEndpointList(c *gin.Context) {
    ctx := c.Request.Context()
    var q req.DockerEndpointSearch
    if err := c.ShouldBindQuery(&q); err != nil { response.FailWithMessage(err.Error(), c); return }
    list, total, err := dockerEndpointService.GetDockerEndpointList(ctx, q)
    if err != nil {
        global.GVA_LOG.Error("获取失败!", zap.Error(err))
        response.FailWithMessage("获取失败:"+err.Error(), c)
        return
    }
    response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: q.Page, PageSize: q.PageSize}, "获取成功", c)
}

