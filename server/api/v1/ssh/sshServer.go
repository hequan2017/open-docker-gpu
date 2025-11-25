package ssh

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ssh"
	sshReq "github.com/flipped-aurora/gin-vue-admin/server/model/ssh/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SshServerApi struct{}

// CreateSshServer 创建Linux SSH管理
// @Tags SshServer
// @Summary 创建Linux SSH管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body ssh.SshServer true "创建Linux SSH管理"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /ssh/createSshServer [post]
func (sshApi *SshServerApi) CreateSshServer(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var ssh ssh.SshServer
	err := c.ShouldBindJSON(&ssh)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = sshService.CreateSshServer(ctx, &ssh)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"status": ssh.Status}, "创建成功", c)
}

// DeleteSshServer 删除Linux SSH管理
// @Tags SshServer
// @Summary 删除Linux SSH管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body ssh.SshServer true "删除Linux SSH管理"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /ssh/deleteSshServer [delete]
func (sshApi *SshServerApi) DeleteSshServer(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	err := sshService.DeleteSshServer(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteSshServerByIds 批量删除Linux SSH管理
// @Tags SshServer
// @Summary 批量删除Linux SSH管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /ssh/deleteSshServerByIds [delete]
func (sshApi *SshServerApi) DeleteSshServerByIds(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := sshService.DeleteSshServerByIds(ctx, IDs)
	if err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:"+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateSshServer 更新Linux SSH管理
// @Tags SshServer
// @Summary 更新Linux SSH管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body ssh.SshServer true "更新Linux SSH管理"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /ssh/updateSshServer [put]
func (sshApi *SshServerApi) UpdateSshServer(c *gin.Context) {
	// 从ctx获取标准context进行业务行为
	ctx := c.Request.Context()

	var ssh ssh.SshServer
	err := c.ShouldBindJSON(&ssh)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = sshService.UpdateSshServer(ctx, &ssh)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(gin.H{"status": ssh.Status}, "更新成功", c)
}

// FindSshServer 用id查询Linux SSH管理
// @Tags SshServer
// @Summary 用id查询Linux SSH管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询Linux SSH管理"
// @Success 200 {object} response.Response{data=ssh.SshServer,msg=string} "查询成功"
// @Router /ssh/findSshServer [get]
func (sshApi *SshServerApi) FindSshServer(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	ID := c.Query("ID")
	ressh, err := sshService.GetSshServer(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:"+err.Error(), c)
		return
	}
	response.OkWithData(ressh, c)
}

// GetSshServerList 分页获取Linux SSH管理列表
// @Tags SshServer
// @Summary 分页获取Linux SSH管理列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query sshReq.SshServerSearch true "分页获取Linux SSH管理列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /ssh/getSshServerList [get]
func (sshApi *SshServerApi) GetSshServerList(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	var pageInfo sshReq.SshServerSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := sshService.GetSshServerInfoList(ctx, pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetSshServerPublic 不需要鉴权的Linux SSH管理接口
// @Tags SshServer
// @Summary 不需要鉴权的Linux SSH管理接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /ssh/getSshServerPublic [get]
func (sshApi *SshServerApi) GetSshServerPublic(c *gin.Context) {
	// 创建业务用Context
	ctx := c.Request.Context()

	// 此接口不需要鉴权
	// 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
	sshService.GetSshServerPublic(ctx)
	response.OkWithDetailed(gin.H{
		"info": "不需要鉴权的Linux SSH管理接口信息",
	}, "获取成功", c)
}

// GetGpuInfo 通过SSH获取GPU信息
// @Tags SshServer
// @Summary 通过SSH获取GPU信息
// @Security ApiKeyAuth
// @Produce application/json
// @Param ID query uint true "SSH记录ID"
// @Router /ssh/gpuInfo [get]
func (sshApi *SshServerApi) GetGpuInfo(c *gin.Context) {
	ctx := c.Request.Context()
	ID := c.Query("ID")
	if ID == "" {
		response.FailWithMessage("参数错误", c)
		return
	}
	items, err := sshService.FetchGPUInfo(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("获取GPU信息失败!", zap.Error(err))
		response.FailWithMessage("获取GPU信息失败:"+err.Error(), c)
		return
	}
	response.OkWithData(items, c)
}

// NvidiaSmiText 返回 nvidia-smi 尾部摘要文本
// @Tags SshServer
// @Summary 返回 nvidia-smi 尾部摘要文本
// @Security ApiKeyAuth
// @Produce application/json
// @Param ID query uint true "SSH记录ID"
// @Param tail query int false "尾部行数，默认50"
// @Router /ssh/nvidiaSmiText [get]
func (sshApi *SshServerApi) NvidiaSmiText(c *gin.Context) {
	ctx := c.Request.Context()
	ID := c.Query("ID")
	tailStr := c.DefaultQuery("tail", "50")
	if ID == "" {
		response.FailWithMessage("参数错误", c)
		return
	}
	tail := 50
	if v, err := strconv.Atoi(tailStr); err == nil {
		tail = v
	}
	txt, err := sshService.FetchNvidiaSmiText(ctx, ID, tail)
	if err != nil {
		global.GVA_LOG.Error("获取nvidia-smi失败!", zap.Error(err))
		response.FailWithMessage("获取失败:"+err.Error(), c)
		return
	}
	response.OkWithData(gin.H{"text": txt}, c)
}
