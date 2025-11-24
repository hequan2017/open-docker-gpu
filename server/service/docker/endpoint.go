package docker

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"strings"

	"github.com/docker/docker/client"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	model "github.com/flipped-aurora/gin-vue-admin/server/model/docker"
	req "github.com/flipped-aurora/gin-vue-admin/server/model/docker/request"
)

type DockerEndpointService struct{}

// CreateDockerEndpoint 创建配置并检测连接状态
func (s *DockerEndpointService) CreateDockerEndpoint(ctx context.Context, m *model.DockerEndpoint) (err error) {
	s.checkAndSetStatus(ctx, m)
	err = global.GVA_DB.Create(m).Error
	return
}

// UpdateDockerEndpoint 更新配置并检测连接状态
func (s *DockerEndpointService) UpdateDockerEndpoint(ctx context.Context, m *model.DockerEndpoint) (err error) {
	s.checkAndSetStatus(ctx, m)
	err = global.GVA_DB.Save(m).Error
	return
}

// DeleteDockerEndpoint 删除
func (s *DockerEndpointService) DeleteDockerEndpoint(ctx context.Context, ID string) (err error) {
	err = global.GVA_DB.Delete(&model.DockerEndpoint{}, "id = ?", ID).Error
	return
}

// DeleteDockerEndpointByIds 批量删除
func (s *DockerEndpointService) DeleteDockerEndpointByIds(ctx context.Context, IDs []string) (err error) {
	err = global.GVA_DB.Delete(&[]model.DockerEndpoint{}, "id in ?", IDs).Error
	return
}

// Find 根据ID查询
func (s *DockerEndpointService) GetDockerEndpoint(ctx context.Context, ID string) (m model.DockerEndpoint, err error) {
	err = global.GVA_DB.Where("id = ?", ID).First(&m).Error
	return
}

// List 列表
func (s *DockerEndpointService) GetDockerEndpointList(ctx context.Context, info req.DockerEndpointSearch) (list []model.DockerEndpoint, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&model.DockerEndpoint{})
	if info.Label != nil && *info.Label != "" {
		db = db.Where("label LIKE ?", "%"+*info.Label+"%")
	}
	if info.Endpoint != nil && *info.Endpoint != "" {
		db = db.Where("endpoint LIKE ?", "%"+*info.Endpoint+"%")
	}
	if info.UseTLS != nil {
		db = db.Where("use_tls = ?", *info.UseTLS)
	}
	if info.Status != nil && *info.Status != "" {
		db = db.Where("status = ?", *info.Status)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	var order string
	if info.Sort != "" {
		order = info.Sort
		if strings.ToLower(info.Order) == "descending" {
			order += " desc"
		}
		db = db.Order(order)
	}
	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err = db.Find(&list).Error
	return
}

// checkAndSetStatus 尝试连接并设置状态
func (s *DockerEndpointService) checkAndSetStatus(ctx context.Context, m *model.DockerEndpoint) {
	status := "异常"
	if m.Endpoint != nil && *m.Endpoint != "" {
		opts := []client.Opt{client.WithHost(*m.Endpoint), client.WithVersion("1.44")}
		if m.UseTLS != nil && *m.UseTLS {
			tlsConf := &tls.Config{}
			// CA
			if m.CACert != nil && *m.CACert != "" {
				pool := x509.NewCertPool()
				pool.AppendCertsFromPEM([]byte(*m.CACert))
				tlsConf.RootCAs = pool
			}
			// client cert/key
			if m.ClientCert != nil && m.ClientKey != nil && *m.ClientCert != "" && *m.ClientKey != "" {
				if pair, e := tls.X509KeyPair([]byte(*m.ClientCert), []byte(*m.ClientKey)); e == nil {
					tlsConf.Certificates = []tls.Certificate{pair}
				}
			}
			httpClient := &http.Client{Transport: &http.Transport{TLSClientConfig: tlsConf}}
			opts = append(opts, client.WithHTTPClient(httpClient))
		}
		if cli, e := client.NewClientWithOpts(opts...); e == nil {
			if _, e2 := cli.Ping(ctx); e2 == nil {
				status = "正常"
			}
		}
	}
	m.Status = &status
}
