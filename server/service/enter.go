package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/service/asset"
	"github.com/flipped-aurora/gin-vue-admin/server/service/docker"
	"github.com/flipped-aurora/gin-vue-admin/server/service/example"
	"github.com/flipped-aurora/gin-vue-admin/server/service/ssh"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
)

var ServiceGroupApp = new(ServiceGroup)

type ServiceGroup struct {
	SystemServiceGroup  system.ServiceGroup
	ExampleServiceGroup example.ServiceGroup
	SshServiceGroup     ssh.ServiceGroup
	DockerServiceGroup  docker.ServiceGroup
	AssetServiceGroup   asset.ServiceGroup
}
