package billing

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/gin-gonic/gin"
)

type GpuRouter struct{}

func (r *GpuRouter) InitBillingGpuRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	gpuApi := v1.ApiGroupApp.BillingApiGroup.GpuApi
	{
		Router.POST("/billing/gpu/collect", gpuApi.Collect)
		Router.GET("/billing/gpu/latest", gpuApi.Latest)
		Router.GET("/billing/gpu/records", gpuApi.Records)
		Router.GET("/billing/gpu/summary", gpuApi.Summary)
	}
}
