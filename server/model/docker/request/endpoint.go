package request

import (
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type DockerEndpointSearch struct {
    Label    *string `json:"label" form:"label"`
    Endpoint *string `json:"endpoint" form:"endpoint"`
    UseTLS   *bool   `json:"useTls" form:"useTls"`
    Status   *string `json:"status" form:"status"`
    request.PageInfo
    Sort  string `json:"sort" form:"sort"`
    Order string `json:"order" form:"order"`
}

