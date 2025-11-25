package asset

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ServerAssetRouter struct {}

// InitServerAssetRouter 初始化 服务器资产表 路由信息
func (s *ServerAssetRouter) InitServerAssetRouter(Router *gin.RouterGroup,PublicRouter *gin.RouterGroup) {
	SrvAstRouter := Router.Group("SrvAst").Use(middleware.OperationRecord())
	SrvAstRouterWithoutRecord := Router.Group("SrvAst")
	SrvAstRouterWithoutAuth := PublicRouter.Group("SrvAst")
	{
		SrvAstRouter.POST("createServerAsset", SrvAstApi.CreateServerAsset)   // 新建服务器资产表
		SrvAstRouter.DELETE("deleteServerAsset", SrvAstApi.DeleteServerAsset) // 删除服务器资产表
		SrvAstRouter.DELETE("deleteServerAssetByIds", SrvAstApi.DeleteServerAssetByIds) // 批量删除服务器资产表
		SrvAstRouter.PUT("updateServerAsset", SrvAstApi.UpdateServerAsset)    // 更新服务器资产表
	}
	{
		SrvAstRouterWithoutRecord.GET("findServerAsset", SrvAstApi.FindServerAsset)        // 根据ID获取服务器资产表
		SrvAstRouterWithoutRecord.GET("getServerAssetList", SrvAstApi.GetServerAssetList)  // 获取服务器资产表列表
	}
	{
	    SrvAstRouterWithoutAuth.GET("getServerAssetDataSource", SrvAstApi.GetServerAssetDataSource)  // 获取服务器资产表数据源
	    SrvAstRouterWithoutAuth.GET("getServerAssetPublic", SrvAstApi.GetServerAssetPublic)  // 服务器资产表开放接口
	}
}
