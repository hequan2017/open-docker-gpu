package ssh

import api "github.com/flipped-aurora/gin-vue-admin/server/api/v1"

type RouterGroup struct{ SshServerRouter }

var sshApi = api.ApiGroupApp.SshApiGroup.SshServerApi
