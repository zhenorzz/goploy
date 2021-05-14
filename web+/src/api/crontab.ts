import Axios from './axios'

/**
 * @return {Promise}
 */
export function getList({ page, rows }, command) {
  return Axios.request({
    url: '/crontab/getList',
    method: 'get',
    params: { page, rows, command },
  })
}

/**
 * @return {Promise}
 */
export function getTotal(command) {
  return Axios.request({
    url: '/crontab/getTotal',
    method: 'get',
    params: { command },
  })
}

/**
 * @return {Promise}
 */
export function getRemoteServerList(serverId) {
  return Axios.request({
    url: '/crontab/getRemoteServerList',
    method: 'get',
    params: { serverId },
    timeout: 0,
  })
}

/**
 * @param  {id} id
 * @return {Promise}
 */
export function getBindServerList(id) {
  return Axios.request({
    url: '/crontab/getBindServerList',
    method: 'get',
    params: { id },
  })
}

export function add(data) {
  return Axios.request({
    url: '/crontab/add',
    method: 'post',
    data,
  })
}

export function edit(data) {
  return Axios.request({
    url: '/crontab/edit',
    method: 'put',
    data,
  })
}

export function importCrontab(data) {
  return Axios.request({
    url: '/crontab/import',
    method: 'post',
    data,
  })
}

export function remove(data) {
  return Axios.request({
    url: '/crontab/remove',
    method: 'delete',
    data,
  })
}

export function addServer(data) {
  return Axios.request({
    url: '/crontab/addServer',
    method: 'post',
    data,
  })
}

export function removeCrontabServer(data) {
  return Axios.request({
    url: '/crontab/removeCrontabServer',
    method: 'delete',
    data,
  })
}
