package docker

import (
	"encoding/json"
	"net/http"

	"github.com/docker/docker/api/types/container"
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
		ws.WriteMessage(websocket.TextMessage, []byte("创建终端失败:"+err.Error()))
		ws.Close()
		return
	}
	attach, err := cli.ContainerExecAttach(ctx, execResp.ID, container.ExecAttachOptions{Tty: true})
	if err != nil {
		ws.WriteMessage(websocket.TextMessage, []byte("连接终端失败:"+err.Error()))
		ws.Close()
		return
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
