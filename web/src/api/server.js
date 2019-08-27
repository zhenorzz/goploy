import request from '@/utils/request'

/**
 * @return {Promise}
 */
export function getList() {
  return request({
    url: '/server/getList',
    method: 'get',
    params: {}
  })
}

/**
 * @return {Promise}
 */
export function getInstallPreview(serverId) {
  return request({
    url: '/server/getInstallPreview',
    method: 'get',
    params: { serverId }
  })
}

/**
 * @return {Promise}
 */
export function getInstallList(token) {
  return request({
    url: '/server/getInstallList',
    method: 'get',
    params: { token }
  })
}

/**
 * @return {Promise}
 */
export function getOption() {
  return request({
    url: '/server/getOption',
    method: 'get'
  })
}

export function add(data) {
  return request({
    url: '/server/add',
    method: 'post',
    data
  })
}

export function edit(data) {
  return request({
    url: '/server/edit',
    method: 'post',
    data
  })
}

export function remove(id) {
  return request({
    url: '/server/remove',
    method: 'post',
    data: { id }
  })
}

export function install(serverId, templateId) {
  return request({
    url: '/server/install',
    method: 'post',
    data: { serverId, templateId }
  })
}
