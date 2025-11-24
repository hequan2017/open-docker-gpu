package ssh

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct{ SshServerApi }

var sshService = service.ServiceGroupApp.SshServiceGroup.SshServerService
