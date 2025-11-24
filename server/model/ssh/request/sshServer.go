package request

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type SshServerSearch struct {
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	Ip             *string     `json:"ip" form:"ip"`
	Port           *int        `json:"port" form:"port"`
	Username       *string     `json:"username" form:"username"`
	Password       *string     `json:"password" form:"password"`
	Status         *string     `json:"status" form:"status"`
	Label          *string     `json:"label" form:"label"`
	Region         *string     `json:"region" form:"region"`
	request.PageInfo
	Sort  string `json:"sort" form:"sort"`
	Order string `json:"order" form:"order"`
}
