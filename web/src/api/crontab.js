import request from '@/utils/request'

/**
 * @return {Promise}
 */
export function getList({ page, rows }, command) {
  return request({
    url: '/crontab/getList',
    method: 'get',
    params: { page, rows, command }
  })
}

/**
 * @return {Promise}
 */
export function getTotal(command) {
  return request({
    url: '/crontab/getTotal',
    method: 'get',
    params: { command }
  })
}

/**
 * @return {Promise}
 */
export function getRemoteServerList(serverId) {
  return request({
    url: '/crontab/getRemoteServerList',
    method: 'get',
    params: { serverId },
    timeout: 0
  })
}

/**
 * @param  {id} id
 * @return {Promise}
 */
export function getBindServerList(id) {
  return request({
    url: '/crontab/getBindServerList',
    method: 'get',
    params: { id }
  })
}

export function add(data) {
  return request({
    url: '/crontab/add',
    method: 'post',
    data
  })
}

export function edit(data) {
  return request({
    url: '/crontab/edit',
    method: 'post',
    data
  })
}

export function importCrontab(data) {
  return request({
    url: '/crontab/import',
    method: 'post',
    data
  })
}

export function remove(data) {
  return request({
    url: '/crontab/remove',
    method: 'delete',
    data
  })
}

export function addServer(data) {
  return request({
    url: '/crontab/addServer',
    method: 'post',
    data
  })
}

export function removeCrontabServer(data) {
  return request({
    url: '/crontab/removeCrontabServer',
    method: 'delete',
    data
  })
}
