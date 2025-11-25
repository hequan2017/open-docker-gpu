package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/router/asset"
	"github.com/flipped-aurora/gin-vue-admin/server/router/docker"
	"github.com/flipped-aurora/gin-vue-admin/server/router/example"
	"github.com/flipped-aurora/gin-vue-admin/server/router/ssh"
	"github.com/flipped-aurora/gin-vue-admin/server/router/system"
)

var RouterGroupApp = new(RouterGroup)

type RouterGroup struct {
	System  system.RouterGroup
	Example example.RouterGroup
	Docker  docker.RouterGroup
	Ssh     ssh.RouterGroup
	Asset   asset.RouterGroup
}
