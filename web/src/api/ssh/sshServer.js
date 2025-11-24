import service from '@/utils/request'
// @Tags SshServer
// @Summary 创建Linux SSH管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SshServer true "创建Linux SSH管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /ssh/createSshServer [post]
export const createSshServer = (data) => {
  return service({
    url: '/ssh/createSshServer',
    method: 'post',
    data
  })
}

// @Tags SshServer
// @Summary 删除Linux SSH管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SshServer true "删除Linux SSH管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /ssh/deleteSshServer [delete]
export const deleteSshServer = (params) => {
  return service({
    url: '/ssh/deleteSshServer',
    method: 'delete',
    params
  })
}

// @Tags SshServer
// @Summary 批量删除Linux SSH管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除Linux SSH管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /ssh/deleteSshServer [delete]
export const deleteSshServerByIds = (params) => {
  return service({
    url: '/ssh/deleteSshServerByIds',
    method: 'delete',
    params
  })
}

// @Tags SshServer
// @Summary 更新Linux SSH管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body model.SshServer true "更新Linux SSH管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /ssh/updateSshServer [put]
export const updateSshServer = (data) => {
  return service({
    url: '/ssh/updateSshServer',
    method: 'put',
    data
  })
}

// @Tags SshServer
// @Summary 用id查询Linux SSH管理
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query model.SshServer true "用id查询Linux SSH管理"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /ssh/findSshServer [get]
export const findSshServer = (params) => {
  return service({
    url: '/ssh/findSshServer',
    method: 'get',
    params
  })
}

// @Tags SshServer
// @Summary 分页获取Linux SSH管理列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取Linux SSH管理列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /ssh/getSshServerList [get]
export const getSshServerList = (params) => {
  return service({
    url: '/ssh/getSshServerList',
    method: 'get',
    params
  })
}

// @Tags SshServer
// @Summary 不需要鉴权的Linux SSH管理接口
// @Accept application/json
// @Produce application/json
// @Param data query sshReq.SshServerSearch true "分页获取Linux SSH管理列表"
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /ssh/getSshServerPublic [get]
export const getSshServerPublic = () => {
  return service({
    url: '/ssh/getSshServerPublic',
    method: 'get',
  })
}
