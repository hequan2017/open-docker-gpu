package asset

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct{ ServerAssetApi }

var SrvAstService = service.ServiceGroupApp.AssetServiceGroup.ServerAssetService
