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
export const startContainer = (params) => {
  return service({ url: '/docker/startContainer', method: 'post', params })
}
export const stopContainer = (params) => {
  return service({ url: '/docker/stopContainer', method: 'post', params })
}
export const removeContainer = (params) => {
  return service({ url: '/docker/removeContainer', method: 'delete', params })
}
