package docker

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct{ DockerRouter }

var dockerApi = api.ApiGroupApp.DockerApiGroup.DockerApi
var dockerEndpointApi = api.ApiGroupApp.DockerApiGroup.DockerEndpointApi
