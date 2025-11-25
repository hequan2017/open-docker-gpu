import service from '@/utils/request'

export const getDockerServers = () => {
  return service({
    url: '/docker/servers',
    method: 'get'
  })
}

export const getDockerPs = (params) => {
  return service({
    url: '/docker/ps',
    method: 'get',
    params
  })
}

export const createContainer = (data) => {
  return service({ url: '/docker/createContainer', method: 'post', data })
}
export const createContainerByDockerfile = (data) => {
  return service({ url: '/docker/createContainerByDockerfile', method: 'post', data })
}
export const createContainerWithOptions = (data) => {
  return service({ url: '/docker/createContainerWithOptions', method: 'post', data })
}
export const startContainer = (params) => {
  return service({ url: '/docker/startContainer', method: 'post', params })
}
export const stopContainer = (params) => {
  return service({ url: '/docker/stopContainer', method: 'post', params })
}
export const removeContainer = (params) => {
  return service({ url: '/docker/removeContainer', method: 'delete', params })
}

export const getDockerLogs = (params) => {
  return service({ url: '/docker/logs', method: 'get', params })
}

export const getPreferredShell = (params) => {
  return service({ url: '/docker/preferredShell', method: 'get', params })
}

export const getDockerImages = (params) => {
  return service({ url: '/docker/images', method: 'get', params })
}

export const getServerAssetList = (params) => {
  return service({ url: '/SrvAst/getServerAssetList', method: 'get', params })
}

export const getServerAssetDataSource = () => {
  return service({ url: '/SrvAst/getServerAssetDataSource', method: 'get' })
}
