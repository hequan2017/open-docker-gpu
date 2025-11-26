package billing

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	billingModel "github.com/flipped-aurora/gin-vue-admin/server/model/billing"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/gin-gonic/gin"
)

type GpuApi struct{}

// Collect 手动触发单次采集
// @Tags Billing
// @Summary 手动采集指定Endpoint的GPU使用率
// @Security ApiKeyAuth
// @Produce application/json
// @Param endpointId query string true "服务器ID"
// @Router /billing/gpu/collect [post]
func (a *GpuApi) Collect(c *gin.Context) {
	ctx := c.Request.Context()
	endpointID := c.Query("endpointId")
	if endpointID == "" || endpointID == "undefined" {
		response.FailWithMessage("参数错误", c)
		return
	}
	if err := gpuMeterService.CollectOneEndpoint(ctx, endpointID); err != nil {
		response.FailWithMessage("采集失败:"+err.Error(), c)
		return
	}
	response.Ok(c)
}

// Latest 获取最近一次记录
// @Tags Billing
// @Summary 获取最近一次记录
// @Security ApiKeyAuth
// @Produce application/json
// @Param endpointId query string true "服务器ID"
// @Router /billing/gpu/latest [get]
func (a *GpuApi) Latest(c *gin.Context) {
	endpointID := c.Query("endpointId")
	if endpointID == "" || endpointID == "undefined" {
		response.FailWithMessage("参数错误", c)
		return
	}
	eid := strToInt64(endpointID)
	var rec billingModel.GPUUsageRecord
	if err := global.GVA_DB.Where("endpoint_id = ?", eid).Order("measured_at desc").First(&rec).Error; err != nil {
		response.Fail(c)
		return
	}
	response.OkWithData(rec, c)
}

// Records 分页查询区间内记录
// @Tags Billing
// @Summary 分页查询区间内记录
// @Security ApiKeyAuth
// @Produce application/json
// @Param endpointId query string true "服务器ID"
// @Param start query string false "开始时间RFC3339"
// @Param end query string false "结束时间RFC3339"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Router /billing/gpu/records [get]
func (a *GpuApi) Records(c *gin.Context) {
	endpointID := c.Query("endpointId")
	if endpointID == "" || endpointID == "undefined" {
		response.FailWithMessage("参数错误", c)
		return
	}
	eid := strToInt64(endpointID)
	startStr := c.Query("start")
	endStr := c.Query("end")
	page := parseIntDefault(c.Query("page"), 1)
	pageSize := parseIntDefault(c.Query("pageSize"), 20)
	db := global.GVA_DB.Model(&billingModel.GPUUsageRecord{})
	db = db.Where("endpoint_id = ?", eid)
	if t, ok := parseTime(startStr); ok {
		db = db.Where("measured_at >= ?", t)
	}
	if t, ok := parseTime(endStr); ok {
		db = db.Where("measured_at <= ?", t)
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		response.Fail(c)
		return
	}
	var list []billingModel.GPUUsageRecord
	if err := db.Order("measured_at desc").Limit(pageSize).Offset((page - 1) * pageSize).Find(&list).Error; err != nil {
		response.Fail(c)
		return
	}
	response.OkWithDetailed(response.PageResult{List: list, Total: total, Page: page, PageSize: pageSize}, "获取成功", c)
}

// Summary 聚合查询（hour/day）
// @Tags Billing
// @Summary 聚合查询（hour/day）
// @Security ApiKeyAuth
// @Produce application/json
// @Param endpointId query string true "服务器ID"
// @Param groupBy query string false "hour|day"
// @Param start query string false "开始时间RFC3339"
// @Param end query string false "结束时间RFC3339"
// @Router /billing/gpu/summary [get]
func (a *GpuApi) Summary(c *gin.Context) {
	endpointID := c.Query("endpointId")
	if endpointID == "" || endpointID == "undefined" {
		response.FailWithMessage("参数错误", c)
		return
	}
	eid := strToInt64(endpointID)
	groupBy := c.DefaultQuery("groupBy", "hour")
	startStr := c.Query("start")
	endStr := c.Query("end")
	start, _ := parseTime(startStr)
	end, _ := parseTime(endStr)
	if end.IsZero() {
		end = time.Now().UTC()
	}
	if start.IsZero() {
		start = end.Add(-24 * time.Hour)
	}
	// 简单聚合：按小时/天分组计算平均使用率与容器数
	type item struct {
		Bucket        string  `json:"bucket"`
		AvgRate       float64 `json:"avgRate"`
		AvgUsed       float64 `json:"avgUsed"`
		AvgCnt        float64 `json:"avgCnt"`
		CostMinuteAvg float64 `json:"costMinuteAvg"`
	}
	var out []item
	// 直接在应用层聚合，避免引入特定数据库函数差异
	var list []billingModel.GPUUsageRecord
	if err := global.GVA_DB.Where("endpoint_id = ? AND measured_at BETWEEN ? AND ?", eid, start, end).Order("measured_at").Find(&list).Error; err != nil {
		response.Fail(c)
		return
	}
	buckets := make(map[string][]billingModel.GPUUsageRecord)
	layout := "2006-01-02 15"
	if groupBy == "day" {
		layout = "2006-01-02"
	}
	for _, r := range list {
		if r.MeasuredAt == nil {
			continue
		}
		key := r.MeasuredAt.UTC().Format(layout)
		buckets[key] = append(buckets[key], r)
	}
	var sysParamSvc system.SysParamsService
	rateCard := getParamFloat(&sysParamSvc, "gpu.rate.perCard.perMinute")
	rateUsage := getParamFloat(&sysParamSvc, "gpu.rate.perUsagePercent.perMinute")
	for k, rs := range buckets {
		var sumRate, sumUsed, sumCnt float64
		for _, r := range rs {
			if r.UsageRate != nil {
				sumRate += *r.UsageRate
			}
			if r.UsedGpuCards != nil {
				sumUsed += float64(*r.UsedGpuCards)
			}
			if r.ContainerCount != nil {
				sumCnt += float64(*r.ContainerCount)
			}
		}
		n := float64(len(rs))
		if n == 0 {
			continue
		}
		avgRate := sumRate / n
		avgUsed := sumUsed / n
		cost := 0.0
		if rateCard > 0 {
			cost = avgUsed * rateCard
		}
		if rateUsage > 0 {
			cost = avgRate * rateUsage
		}
		out = append(out, item{Bucket: k, AvgRate: avgRate, AvgUsed: avgUsed, AvgCnt: sumCnt / n, CostMinuteAvg: cost})
	}
	response.OkWithData(out, c)
}

// 工具
func parseIntDefault(s string, def int) int {
	v := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return def
		}
		v = v*10 + int(c-'0')
	}
	if v == 0 {
		return def
	}
	return v
}

func parseTime(s string) (time.Time, bool) {
	if s == "" {
		return time.Time{}, false
	}
	t, e := time.Parse(time.RFC3339, s)
	if e != nil {
		return time.Time{}, false
	}
	return t, true
}

func strToInt64(s string) int64 {
	var v int64
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return 0
		}
		v = v*10 + int64(c-'0')
	}
	return v
}

func getParamFloat(svc *system.SysParamsService, key string) float64 {
	p, err := svc.GetSysParam(key)
	if err != nil {
		return 0
	}
	// 解析浮点（支持简单整数/小数）
	s := p.Value
	if s == "" {
		return 0
	}
	// 使用标准库需引入 strconv；为避免依赖，此处简化实现
	// 找到小数点位置
	dot := -1
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			dot = i
			break
		}
		if s[i] < '0' || s[i] > '9' {
			return 0
		}
	}
	v := 0.0
	if dot == -1 {
		for i := 0; i < len(s); i++ {
			v = v*10 + float64(s[i]-'0')
		}
		return v
	}
	for i := 0; i < dot; i++ {
		v = v*10 + float64(s[i]-'0')
	}
	frac := 0.0
	pow := 1.0
	for i := dot + 1; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return 0
		}
		frac = frac*10 + float64(s[i]-'0')
		pow *= 10
	}
	return v + frac/pow
}
