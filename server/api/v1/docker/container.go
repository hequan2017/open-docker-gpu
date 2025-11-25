package docker

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types/container"
	imageTypes "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type CreateContainerReq struct {
	EndpointID string `json:"endpointId" form:"endpointId"`
	Image      string `json:"image" form:"image"`
	Name       string `json:"name" form:"name"`
}

type BuildContainerReq struct {
	EndpointID string `json:"endpointId" form:"endpointId"`
	Dockerfile string `json:"dockerfile" form:"dockerfile"`
	Tag        string `json:"tag" form:"tag"`
	Name       string `json:"name" form:"name"`
}

type CreateContainerWithOptionsReq struct {
	EndpointID string   `json:"endpointId" form:"endpointId"`
	Image      string   `json:"image" form:"image"`
	Name       string   `json:"name" form:"name"`
	WorkingDir string   `json:"workingDir" form:"workingDir"`
	Env        []string `json:"env"`
	Ports      []string `json:"ports"`
	Volumes    []string `json:"volumes"`
	Cmd        []string `json:"cmd"`
	GpuEnabled bool     `json:"gpuEnabled"`
	GpuMode    string   `json:"gpuMode"`
	GpuCount   int      `json:"gpuCount"`
	GpuDevices []string `json:"gpuDevices"`
}

func (a *DockerApi) CreateContainer(c *gin.Context) {
	ctx := c.Request.Context()
	var req CreateContainerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	id, err := dockerService.CreateContainer(ctx, req.EndpointID, req.Image, req.Name)
	if err != nil {
		global.GVA_LOG.Error("创建容器失败!", zap.Error(err))
		response.FailWithMessage("创建容器失败:"+err.Error(), c)
		return
	}
	response.OkWithData(gin.H{"id": id}, c)
}

func (a *DockerApi) BuildContainerByDockerfile(c *gin.Context) {
	ctx := c.Request.Context()
	var req BuildContainerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	id, err := dockerService.BuildContainerFromDockerfile(ctx, req.EndpointID, req.Dockerfile, req.Tag, req.Name)
	if err != nil {
		global.GVA_LOG.Error("创建容器失败!", zap.Error(err))
		response.FailWithMessage("创建容器失败:"+err.Error(), c)
		return
	}
	response.OkWithData(gin.H{"id": id}, c)
}

func (a *DockerApi) CreateContainerWithOptions(c *gin.Context) {
	ctx := c.Request.Context()
	var req CreateContainerWithOptionsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var gpuEnabled bool
	var gpuAll bool
	var gpuCount int
	var gpuDevices []string
	if strings.ToLower(req.GpuMode) != "" && req.GpuEnabled {
		gpuEnabled = true
		switch strings.ToLower(req.GpuMode) {
		case "all":
			gpuAll = true
		case "count":
			gpuCount = req.GpuCount
		case "devices":
			gpuDevices = req.GpuDevices
		}
	}
	id, err := dockerService.CreateContainerWithOptions(ctx, req.EndpointID, req.Image, req.Name, req.WorkingDir, req.Env, req.Ports, req.Volumes, req.Cmd, gpuEnabled, gpuAll, gpuCount, gpuDevices)
	if err != nil {
		global.GVA_LOG.Error("创建容器失败!", zap.Error(err))
		response.FailWithMessage("创建容器失败:"+err.Error(), c)
		return
	}
	response.OkWithData(gin.H{"id": id}, c)
}

func (a *DockerApi) StartContainer(c *gin.Context) {
	ctx := c.Request.Context()
	endpointID := c.Query("endpointId")
	cid := c.Query("ID")
	if err := dockerService.StartContainer(ctx, endpointID, cid); err != nil {
		global.GVA_LOG.Error("启动容器失败!", zap.Error(err))
		response.FailWithMessage("启动容器失败:"+err.Error(), c)
		return
	}
	response.Ok(c)
}

func (a *DockerApi) StopContainer(c *gin.Context) {
	ctx := c.Request.Context()
	endpointID := c.Query("endpointId")
	cid := c.Query("ID")
	if err := dockerService.StopContainer(ctx, endpointID, cid); err != nil {
		global.GVA_LOG.Error("停止容器失败!", zap.Error(err))
		response.FailWithMessage("停止容器失败:"+err.Error(), c)
		return
	}
	response.Ok(c)
}

func (a *DockerApi) RemoveContainer(c *gin.Context) {
	ctx := c.Request.Context()
	endpointID := c.Query("endpointId")
	cid := c.Query("ID")
	force := c.Query("force") == "true"
	if err := dockerService.RemoveContainer(ctx, endpointID, cid, force); err != nil {
		global.GVA_LOG.Error("删除容器失败!", zap.Error(err))
		response.FailWithMessage("删除容器失败:"+err.Error(), c)
		return
	}
	response.Ok(c)
}

type wsWriter struct{ ws *websocket.Conn }

func (w *wsWriter) Write(p []byte) (int, error) {
	if err := w.ws.WriteMessage(websocket.TextMessage, p); err != nil {
		return 0, err
	}
	return len(p), nil
}

func (a *DockerApi) ContainerLogsWs(c *gin.Context) {
	ctx := c.Request.Context()
	endpointID := c.Query("endpointId")
	cid := c.Query("ID")
	if endpointID == "" || cid == "" || endpointID == "undefined" {
		response.FailWithMessage("参数错误", c)
		return
	}
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	cli, err := dockerService.GetClientByEndpointID(ctx, endpointID)
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("连接Docker失败:"+err.Error()))
		ws.Close()
		return
	}
	stdout := c.DefaultQuery("stdout", "true") == "true"
	stderr := c.DefaultQuery("stderr", "true") == "true"
	follow := c.DefaultQuery("follow", "true") == "true"
	tail := c.DefaultQuery("tail", "200")
	since := c.DefaultQuery("since", "")
	timestamps := c.DefaultQuery("timestamps", "false") == "true"
	rc, err := cli.ContainerLogs(ctx, cid, container.LogsOptions{ShowStdout: stdout, ShowStderr: stderr, Follow: follow, Tail: tail, Since: since, Timestamps: timestamps})
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("获取日志失败:"+err.Error()))
		ws.Close()
		return
	}
	defer rc.Close()
	w := &wsWriter{ws: ws}
	go func() {
		_, _ = stdcopy.StdCopy(w, w, rc)
		ws.WriteMessage(websocket.TextMessage, []byte("日志流结束"))
		ws.Close()
	}()
	for {
		_, _, e := ws.ReadMessage()
		if e != nil {
			break
		}
	}
}

// GetContainerLogs 获取容器尾部日志（HTTP一次性返回）
// @Tags Docker
// @Summary 获取容器尾部日志
// @Security ApiKeyAuth
// @Produce application/json
// @Param endpointId query string true "服务器ID"
// @Param ID query string true "容器ID"
// @Param tail query string false "尾部行数，默认200"
// @Param since query string false "开始时间"
// @Param timestamps query bool false "是否包含时间戳"
// @Param stdout query bool false "是否包含stdout"
// @Param stderr query bool false "是否包含stderr"
// @Router /docker/logs [get]
func (a *DockerApi) GetContainerLogs(c *gin.Context) {
	ctx := c.Request.Context()
	endpointID := c.Query("endpointId")
	cid := c.Query("ID")
	if endpointID == "" || cid == "" || endpointID == "undefined" {
		response.FailWithMessage("参数错误", c)
		return
	}
	stdout := c.DefaultQuery("stdout", "true") == "true"
	stderr := c.DefaultQuery("stderr", "true") == "true"
	tail := c.DefaultQuery("tail", "200")
	since := c.DefaultQuery("since", "")
	timestamps := c.DefaultQuery("timestamps", "false") == "true"
	cli, err := dockerService.GetClientByEndpointID(ctx, endpointID)
	if err != nil {
		response.FailWithMessage("连接Docker失败:"+err.Error(), c)
		return
	}
	rc, err := cli.ContainerLogs(ctx, cid, container.LogsOptions{ShowStdout: stdout, ShowStderr: stderr, Follow: false, Tail: tail, Since: since, Timestamps: timestamps})
	if err != nil {
		response.FailWithMessage("获取日志失败:"+err.Error(), c)
		return
	}
	defer rc.Close()
	var bufStdout, bufStderr strings.Builder
	_, _ = stdcopy.StdCopy(&bufStdout, &bufStderr, rc)
	combined := bufStdout.String()
	if bufStderr.Len() > 0 {
		combined = combined + bufStderr.String()
	}
	response.OkWithData(gin.H{"text": combined}, c)
}

// GetDockerImages 获取Docker镜像列表
// @Tags Docker
// @Summary 获取镜像列表
// @Security ApiKeyAuth
// @Produce application/json
// @Param endpointId query string true "服务器ID"
// @Router /docker/images [get]
func (a *DockerApi) GetDockerImages(c *gin.Context) {
	ctx := c.Request.Context()
	endpointID := c.Query("endpointId")
	if endpointID == "" || endpointID == "undefined" {
		response.FailWithMessage("参数错误", c)
		return
	}
	cli, err := dockerService.GetClientByEndpointID(ctx, endpointID)
	if err != nil {
		response.FailWithMessage("连接Docker失败:"+err.Error(), c)
		return
	}
	list, e := cli.ImageList(ctx, imageTypes.ListOptions{})
	if e != nil {
		response.FailWithMessage("获取镜像失败:"+e.Error(), c)
		return
	}
	type item struct {
		ID    string   `json:"id"`
		Tags  []string `json:"tags"`
		Size  int64    `json:"size"`
		DSize int64    `json:"dsize"`
	}
	var out []item
	for _, im := range list {
		out = append(out, item{ID: im.ID, Tags: im.RepoTags, Size: im.Size, DSize: im.SharedSize})
	}
	response.OkWithData(out, c)
}

// PreferredShell 预检测容器内可用的交互Shell
// @Tags Docker
// @Summary 预检测容器内可用的交互Shell
// @Security ApiKeyAuth
// @Produce application/json
// @Param endpointId query string true "服务器ID"
// @Param ID query string true "容器ID"
// @Router /docker/preferredShell [get]
func (a *DockerApi) PreferredShell(c *gin.Context) {
	ctx := c.Request.Context()
	endpointID := c.Query("endpointId")
	cid := c.Query("ID")
	if endpointID == "" || cid == "" || endpointID == "undefined" {
		response.FailWithMessage("参数错误", c)
		return
	}
	shell, err := dockerService.DetectPreferredShell(ctx, endpointID, cid)
	if err != nil {
		response.FailWithMessage("检测失败:"+err.Error(), c)
		return
	}
	response.OkWithData(gin.H{"shell": shell}, c)
}

// ContainerTerminalWs WebSocket 进入容器交互终端
// @Tags Docker
// @Summary WebSocket 进入容器交互终端
// @Security ApiKeyAuth
// @Produce application/json
// @Param endpointId query string true "服务器ID"
// @Param ID query string true "容器ID"
// @Router /docker/execWs [get]
func (a *DockerApi) ContainerTerminalWs(c *gin.Context) {
	ctx := c.Request.Context()
	endpointID := c.Query("endpointId")
	cid := c.Query("ID")
	shell := c.DefaultQuery("shell", "/bin/sh")
	if endpointID == "" || cid == "" || endpointID == "undefined" {
		response.FailWithMessage("参数错误", c)
		return
	}
	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	cli, err := dockerService.GetClientByEndpointID(ctx, endpointID)
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("连接Docker失败:"+err.Error()))
		ws.Close()
		return
	}
	execResp, err := cli.ContainerExecCreate(ctx, cid, container.ExecOptions{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Cmd:          []string{shell},
	})
	if err != nil {
		if shell != "/bin/sh" {
			execResp2, err2 := cli.ContainerExecCreate(ctx, cid, container.ExecOptions{
				AttachStdin:  true,
				AttachStdout: true,
				AttachStderr: true,
				Tty:          true,
				Cmd:          []string{"/bin/sh"},
			})
			if err2 == nil {
				_ = ws.WriteMessage(websocket.TextMessage, []byte("已自动切换到 /bin/sh"))
				execResp = execResp2
				shell = "/bin/sh"
			} else {
				ws.WriteMessage(websocket.TextMessage, []byte("创建终端失败:"+err.Error()))
				ws.Close()
				return
			}
		} else {
			ws.WriteMessage(websocket.TextMessage, []byte("创建终端失败:"+err.Error()))
			ws.Close()
			return
		}
	}
	attach, err := cli.ContainerExecAttach(ctx, execResp.ID, container.ExecAttachOptions{Tty: true})
	if err != nil {
		if shell != "/bin/sh" {
			execResp2, err2 := cli.ContainerExecCreate(ctx, cid, container.ExecOptions{
				AttachStdin:  true,
				AttachStdout: true,
				AttachStderr: true,
				Tty:          true,
				Cmd:          []string{"/bin/sh"},
			})
			if err2 == nil {
				_ = ws.WriteMessage(websocket.TextMessage, []byte("已自动切换到 /bin/sh"))
				execResp = execResp2
				shell = "/bin/sh"
				attach, err = cli.ContainerExecAttach(ctx, execResp.ID, container.ExecAttachOptions{Tty: true})
			}
		}
		if err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte("连接终端失败:"+err.Error()))
			ws.Close()
			return
		}
	}
	defer attach.Close()
	_ = cli.ContainerExecStart(ctx, execResp.ID, container.ExecStartOptions{Tty: true})

	go func() {
		buf := make([]byte, 4096)
		for {
			n, e := attach.Reader.Read(buf)
			if n > 0 {
				_ = ws.WriteMessage(websocket.BinaryMessage, buf[:n])
			}
			if e != nil {
				_ = ws.WriteMessage(websocket.TextMessage, []byte("终端已结束"))
				ws.Close()
				break
			}
		}
	}()
	for {
		mt, data, e := ws.ReadMessage()
		if e != nil {
			break
		}
		if mt == websocket.TextMessage {
			type resizeMsg struct {
				Type string `json:"type"`
				Cols int    `json:"cols"`
				Rows int    `json:"rows"`
			}
			var rm resizeMsg
			if json.Unmarshal(data, &rm) == nil {
				if rm.Type == "resize" && rm.Cols > 0 && rm.Rows > 0 {
					_ = cli.ContainerExecResize(ctx, execResp.ID, container.ResizeOptions{Width: uint(rm.Cols), Height: uint(rm.Rows)})
					continue
				}
				if rm.Type == "ping" {
					_ = ws.WriteMessage(websocket.TextMessage, []byte("pong"))
					continue
				}
			}
			if attach.Conn != nil {
				_, _ = attach.Conn.Write(data)
			}
		} else if mt == websocket.BinaryMessage {
			if attach.Conn != nil {
				_, _ = attach.Conn.Write(data)
			}
		}
	}
}
