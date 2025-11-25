package asset

import (
	
	"github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/model/asset"
    assetReq "github.com/flipped-aurora/gin-vue-admin/server/model/asset/request"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

type ServerAssetApi struct {}



// CreateServerAsset 创建服务器资产表
// @Tags ServerAsset
// @Summary 创建服务器资产表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body asset.ServerAsset true "创建服务器资产表"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /SrvAst/createServerAsset [post]
func (SrvAstApi *ServerAssetApi) CreateServerAsset(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var SrvAst asset.ServerAsset
	err := c.ShouldBindJSON(&SrvAst)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = SrvAstService.CreateServerAsset(ctx,&SrvAst)
	if err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败:" + err.Error(), c)
		return
	}
    response.OkWithMessage("创建成功", c)
}

// DeleteServerAsset 删除服务器资产表
// @Tags ServerAsset
// @Summary 删除服务器资产表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body asset.ServerAsset true "删除服务器资产表"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /SrvAst/deleteServerAsset [delete]
func (SrvAstApi *ServerAssetApi) DeleteServerAsset(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	err := SrvAstService.DeleteServerAsset(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteServerAssetByIds 批量删除服务器资产表
// @Tags ServerAsset
// @Summary 批量删除服务器资产表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /SrvAst/deleteServerAssetByIds [delete]
func (SrvAstApi *ServerAssetApi) DeleteServerAssetByIds(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	err := SrvAstService.DeleteServerAssetByIds(ctx,IDs)
	if err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateServerAsset 更新服务器资产表
// @Tags ServerAsset
// @Summary 更新服务器资产表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body asset.ServerAsset true "更新服务器资产表"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /SrvAst/updateServerAsset [put]
func (SrvAstApi *ServerAssetApi) UpdateServerAsset(c *gin.Context) {
    // 从ctx获取标准context进行业务行为
    ctx := c.Request.Context()

	var SrvAst asset.ServerAsset
	err := c.ShouldBindJSON(&SrvAst)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = SrvAstService.UpdateServerAsset(ctx,SrvAst)
	if err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败:" + err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindServerAsset 用id查询服务器资产表
// @Tags ServerAsset
// @Summary 用id查询服务器资产表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "用id查询服务器资产表"
// @Success 200 {object} response.Response{data=asset.ServerAsset,msg=string} "查询成功"
// @Router /SrvAst/findServerAsset [get]
func (SrvAstApi *ServerAssetApi) FindServerAsset(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	ID := c.Query("ID")
	reSrvAst, err := SrvAstService.GetServerAsset(ctx,ID)
	if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败:" + err.Error(), c)
		return
	}
	response.OkWithData(reSrvAst, c)
}
// GetServerAssetList 分页获取服务器资产表列表
// @Tags ServerAsset
// @Summary 分页获取服务器资产表列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query assetReq.ServerAssetSearch true "分页获取服务器资产表列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /SrvAst/getServerAssetList [get]
func (SrvAstApi *ServerAssetApi) GetServerAssetList(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

	var pageInfo assetReq.ServerAssetSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := SrvAstService.GetServerAssetInfoList(ctx,pageInfo)
	if err != nil {
	    global.GVA_LOG.Error("获取失败!", zap.Error(err))
        response.FailWithMessage("获取失败:" + err.Error(), c)
        return
    }
    response.OkWithDetailed(response.PageResult{
        List:     list,
        Total:    total,
        Page:     pageInfo.Page,
        PageSize: pageInfo.PageSize,
    }, "获取成功", c)
}
// GetServerAssetDataSource 获取ServerAsset的数据源
// @Tags ServerAsset
// @Summary 获取ServerAsset的数据源
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "查询成功"
// @Router /SrvAst/getServerAssetDataSource [get]
func (SrvAstApi *ServerAssetApi) GetServerAssetDataSource(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口为获取数据源定义的数据
    dataSource, err := SrvAstService.GetServerAssetDataSource(ctx)
    if err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
   		response.FailWithMessage("查询失败:" + err.Error(), c)
   		return
    }
   response.OkWithData(dataSource, c)
}

// GetServerAssetPublic 不需要鉴权的服务器资产表接口
// @Tags ServerAsset
// @Summary 不需要鉴权的服务器资产表接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /SrvAst/getServerAssetPublic [get]
func (SrvAstApi *ServerAssetApi) GetServerAssetPublic(c *gin.Context) {
    // 创建业务用Context
    ctx := c.Request.Context()

    // 此接口不需要鉴权
    // 示例为返回了一个固定的消息接口，一般本接口用于C端服务，需要自己实现业务逻辑
    SrvAstService.GetServerAssetPublic(ctx)
    response.OkWithDetailed(gin.H{
       "info": "不需要鉴权的服务器资产表接口信息",
    }, "获取成功", c)
}
