package v1

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/asset"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/docker"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/example"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/ssh"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/system"
)

var ApiGroupApp = new(ApiGroup)

type ApiGroup struct {
	SystemApiGroup  system.ApiGroup
	ExampleApiGroup example.ApiGroup
	DockerApiGroup  docker.ApiGroup
	SshApiGroup     ssh.ApiGroup
	AssetApiGroup   asset.ApiGroup
}
