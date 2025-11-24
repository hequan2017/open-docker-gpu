import service from '@/utils/request'

export const createDockerEndpoint = (data) => {
  return service({ url: '/docker/createDockerEndpoint', method: 'post', data })
}
export const updateDockerEndpoint = (data) => {
  return service({ url: '/docker/updateDockerEndpoint', method: 'put', data })
}
export const deleteDockerEndpoint = (params) => {
  return service({ url: '/docker/deleteDockerEndpoint', method: 'delete', params })
}
export const deleteDockerEndpointByIds = (params) => {
  return service({ url: '/docker/deleteDockerEndpointByIds', method: 'delete', params })
}
export const findDockerEndpoint = (params) => {
  return service({ url: '/docker/findDockerEndpoint', method: 'get', params })
}
export const getDockerEndpointList = (params) => {
  return service({ url: '/docker/getDockerEndpointList', method: 'get', params })
}

