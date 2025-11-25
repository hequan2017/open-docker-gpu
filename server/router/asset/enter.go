package asset

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct{ ServerAssetRouter }

var SrvAstApi = api.ApiGroupApp.AssetApiGroup.ServerAssetApi
