package billing

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct{ GpuApi }

var gpuMeterService = service.ServiceGroupApp.MeterServiceGroup.GpuMeterService
