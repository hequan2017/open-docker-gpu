package docker

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	DockerApi
	DockerEndpointApi
}

var dockerService = service.ServiceGroupApp.DockerServiceGroup.DockerService
var dockerEndpointService = service.ServiceGroupApp.DockerServiceGroup.DockerEndpointService
