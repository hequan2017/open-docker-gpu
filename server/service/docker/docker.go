package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	imageTypes "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	dockermodel "github.com/flipped-aurora/gin-vue-admin/server/model/docker"
)

type DockerService struct{}

// ListNormalServers 返回状态为正常的DockerEndpoint列表
func (s *DockerService) ListNormalServers(ctx context.Context) (list []dockermodel.DockerEndpoint, err error) {
	err = global.GVA_DB.Where("status = ?", "正常").Find(&list).Error
	return
}

// DockerPsRow 代表 docker ps -a --format '{{json .}}' 的单行结果
type DockerPsRow struct {
	ID         string `json:"ID"`
	Image      string `json:"Image"`
	Command    string `json:"Command"`
	CreatedAt  string `json:"CreatedAt"`
	RunningFor string `json:"RunningFor"`
	Status     string `json:"Status"`
	Ports      string `json:"Ports"`
	Names      string `json:"Names"`
	Labels     string `json:"Labels"`
	LogsWs     string `json:"LogsWs"`
}

// FetchDockerPs 通过Docker SDK获取容器列表（基于DockerEndpoint配置）
func (s *DockerService) FetchDockerPs(ctx context.Context, endpointID string, scope string) (rows []DockerPsRow, err error) {
	cli, err := s.getClientByEndpointID(ctx, endpointID)
	if err != nil {
		return
	}

	list, e2 := cli.ContainerList(ctx, container.ListOptions{All: true})
	if e2 != nil {
		err = e2
		return
	}

	now := time.Now()
	for _, c := range list {
		if scope == "running" && strings.ToLower(c.State) != "running" {
			continue
		}
		if scope == "exited" && strings.ToLower(c.State) == "running" {
			continue
		}
		var ports []string
		for _, p := range c.Ports {
			pp := strconv.Itoa(int(p.PublicPort))
			hp := strconv.Itoa(int(p.PrivatePort))
			if p.PublicPort == 0 {
				pp = ""
			}
			proto := p.Type
			if pp != "" {
				ports = append(ports, fmt.Sprintf("%s->%s/%s", pp, hp, proto))
			} else {
				ports = append(ports, fmt.Sprintf("%s/%s", hp, proto))
			}
		}
		names := strings.Join(c.Names, ",")
		created := time.Unix(c.Created, 0)
		runningFor := now.Sub(created).Round(time.Second)
		var labelPairs []string
		for k, v := range c.Labels {
			labelPairs = append(labelPairs, k+"="+v)
		}
		row := DockerPsRow{
			ID:         c.ID,
			Image:      c.Image,
			Command:    c.Command,
			CreatedAt:  created.Format("2006-01-02 15:04:05"),
			RunningFor: runningFor.String(),
			Status:     c.Status,
			Ports:      strings.Join(ports, ","),
			Names:      names,
			Labels:     strings.Join(labelPairs, ","),
			LogsWs:     "/docker/logsWs?endpointId=" + endpointID + "&ID=" + c.ID,
		}
		rows = append(rows, row)
	}
	return
}

func (s *DockerService) getClientByEndpointID(ctx context.Context, endpointID string) (*client.Client, error) {
	var ep dockermodel.DockerEndpoint
	if err := global.GVA_DB.Where("id = ?", endpointID).First(&ep).Error; err != nil {
		return nil, err
	}
	if ep.Endpoint == nil || *ep.Endpoint == "" {
		return nil, fmt.Errorf("连接地址为空")
	}
	opts := []client.Opt{client.WithHost(*ep.Endpoint), client.WithVersion("1.44")}
	if ep.UseTLS != nil && *ep.UseTLS {
		tlsConf := &tls.Config{}
		if ep.CACert != nil && *ep.CACert != "" {
			pool := x509.NewCertPool()
			pool.AppendCertsFromPEM([]byte(*ep.CACert))
			tlsConf.RootCAs = pool
		}
		if ep.ClientCert != nil && ep.ClientKey != nil && *ep.ClientCert != "" && *ep.ClientKey != "" {
			if pair, e := tls.X509KeyPair([]byte(*ep.ClientCert), []byte(*ep.ClientKey)); e == nil {
				tlsConf.Certificates = []tls.Certificate{pair}
			}
		}
		httpClient := &http.Client{Transport: &http.Transport{TLSClientConfig: tlsConf}}
		opts = append(opts, client.WithHTTPClient(httpClient))
	}
	return client.NewClientWithOpts(opts...)
}

// GetClientByEndpointID 导出获取Docker客户端的方法，供API层复用
func (s *DockerService) GetClientByEndpointID(ctx context.Context, endpointID string) (*client.Client, error) {
	return s.getClientByEndpointID(ctx, endpointID)
}

func (s *DockerService) CreateContainer(ctx context.Context, endpointID, image, name string) (id string, err error) {
	cli, err := s.getClientByEndpointID(ctx, endpointID)
	if err != nil {
		return
	}
	rc, e := cli.ImagePull(ctx, image, imageTypes.PullOptions{})
	if e == nil && rc != nil {
		io.Copy(io.Discard, rc)
		rc.Close()
	}
	resp, e2 := cli.ContainerCreate(ctx, &container.Config{Image: image}, &container.HostConfig{}, nil, nil, name)
	if e2 != nil {
		err = e2
		return
	}
	if e3 := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); e3 != nil {
		err = e3
		return
	}
	id = resp.ID
	return
}

func (s *DockerService) StartContainer(ctx context.Context, endpointID, cid string) (err error) {
	cli, err := s.getClientByEndpointID(ctx, endpointID)
	if err != nil {
		return
	}
	err = cli.ContainerStart(ctx, cid, container.StartOptions{})
	return
}

func (s *DockerService) StopContainer(ctx context.Context, endpointID, cid string) (err error) {
	cli, err := s.getClientByEndpointID(ctx, endpointID)
	if err != nil {
		return
	}
	err = cli.ContainerStop(ctx, cid, container.StopOptions{})
	return
}

func (s *DockerService) BuildContainerFromDockerfile(ctx context.Context, endpointID, dockerfile, tag, name string) (id string, err error) {
	cli, err := s.getClientByEndpointID(ctx, endpointID)
	if err != nil {
		return
	}
	if tag == "" {
		tag = "dockerfile:latest"
	}
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	hdr := &tar.Header{Name: "Dockerfile", Mode: 0600, Size: int64(len(dockerfile))}
	if e := tw.WriteHeader(hdr); e != nil {
		err = e
		tw.Close()
		return
	}
	if _, e := tw.Write([]byte(dockerfile)); e != nil {
		err = e
		tw.Close()
		return
	}
	if e := tw.Close(); e != nil {
		err = e
		return
	}
	buildResp, e := cli.ImageBuild(ctx, bytes.NewReader(buf.Bytes()), types.ImageBuildOptions{Dockerfile: "Dockerfile", Tags: []string{tag}})
	if e != nil {
		err = e
		return
	}
	if buildResp.Body != nil {
		io.Copy(io.Discard, buildResp.Body)
		buildResp.Body.Close()
	}
	resp, e2 := cli.ContainerCreate(ctx, &container.Config{Image: tag}, &container.HostConfig{}, nil, nil, name)
	if e2 != nil {
		err = e2
		return
	}
	if e3 := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); e3 != nil {
		err = e3
		return
	}
	id = resp.ID
	return
}

func (s *DockerService) CreateContainerWithOptions(ctx context.Context, endpointID, image, name, workingDir string, env, ports, volumes, cmd []string, gpuEnabled bool, gpuAll bool, gpuCount int, gpuDevices []string) (id string, err error) {
	cli, err := s.getClientByEndpointID(ctx, endpointID)
	if err != nil {
		return
	}
	if image == "" {
		err = fmt.Errorf("镜像为空")
		return
	}
	if rc, e := cli.ImagePull(ctx, image, imageTypes.PullOptions{}); e == nil && rc != nil {
		io.Copy(io.Discard, rc)
		rc.Close()
	}
	cfg := &container.Config{Image: image, Env: env, WorkingDir: workingDir}
	if len(cmd) > 0 {
		cfg.Cmd = cmd
	}
	exposed := nat.PortSet{}
	bindings := nat.PortMap{}
	for _, p := range ports {
		if p == "" {
			continue
		}
		proto := "tcp"
		s2 := p
		if strings.Contains(p, "/") {
			parts := strings.SplitN(p, "/", 2)
			s2 = parts[0]
			proto = strings.ToLower(parts[1])
		}
		hp := ""
		cp := ""
		seg := strings.SplitN(s2, ":", 2)
		if len(seg) == 2 {
			hp = seg[0]
			cp = seg[1]
		} else {
			cp = seg[0]
		}
		portKey, e := nat.NewPort(proto, cp)
		if e != nil {
			continue
		}
		exposed[portKey] = struct{}{}
		if _, ok := bindings[portKey]; !ok {
			bindings[portKey] = []nat.PortBinding{}
		}
		bindings[portKey] = append(bindings[portKey], nat.PortBinding{HostPort: hp})
	}
	if len(exposed) > 0 {
		cfg.ExposedPorts = exposed
	}
	host := &container.HostConfig{}
	if len(bindings) > 0 {
		host.PortBindings = bindings
	}
	if len(volumes) > 0 {
		host.Binds = volumes
	}
	if gpuEnabled {
		dr := container.DeviceRequest{Capabilities: [][]string{{"gpu"}}}
		if gpuAll {
			dr.Count = -1
		} else if len(gpuDevices) > 0 {
			dr.DeviceIDs = gpuDevices
		} else if gpuCount > 0 {
			dr.Count = gpuCount
		}
		host.DeviceRequests = append(host.DeviceRequests, dr)
	}
	resp, e2 := cli.ContainerCreate(ctx, cfg, host, nil, nil, name)
	if e2 != nil {
		err = e2
		return
	}
	if e3 := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); e3 != nil {
		err = e3
		return
	}
	id = resp.ID
	return
}

func (s *DockerService) RemoveContainer(ctx context.Context, endpointID, cid string, force bool) (err error) {
	cli, err := s.getClientByEndpointID(ctx, endpointID)
	if err != nil {
		return
	}
	err = cli.ContainerRemove(ctx, cid, container.RemoveOptions{Force: force, RemoveVolumes: true})
	return
}

func (s *DockerService) tryExec(cli *client.Client, ctx context.Context, cid string, cmd []string) error {
	execResp, err := cli.ContainerExecCreate(ctx, cid, container.ExecOptions{
		AttachStdin:  false,
		AttachStdout: false,
		AttachStderr: false,
		Tty:          true,
		Cmd:          cmd,
	})
	if err != nil {
		return err
	}
	return cli.ContainerExecStart(ctx, execResp.ID, container.ExecStartOptions{Tty: true})
}

func (s *DockerService) DetectPreferredShell(ctx context.Context, endpointID, cid string) (shell string, err error) {
	cli, err := s.getClientByEndpointID(ctx, endpointID)
	if err != nil {
		return "", err
	}
	if e := s.tryExec(cli, ctx, cid, []string{"/bin/bash", "-lc", "echo ready"}); e == nil {
		return "/bin/bash", nil
	}
	if e := s.tryExec(cli, ctx, cid, []string{"/bin/sh", "-lc", "echo ready"}); e == nil {
		return "/bin/sh", nil
	}
	return "/bin/sh", nil
}
