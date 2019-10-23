import request from '@/utils/request'

/**
 * @return {Promise}
 */
export function getList(pagination) {
  return request({
    url: '/project/getList',
    method: 'get',
    params: { ...pagination }
  })
}

/**
 * @return {Promise}
 */
export function getOption() {
  return request({
    url: '/project/getOption',
    method: 'get',
    params: { }
  })
}

/**
 * @param  {id} id
 * @return {Promise}
 */
export function getBindServerList(id) {
  return request({
    url: '/project/getBindServerList',
    method: 'get',
    params: { id }
  })
}

/**
 * @param  {id} id
 * @return {Promise}
 */
export function getBindUserList(id) {
  return request({
    url: '/project/getBindUserList',
    method: 'get',
    params: { id }
  })
}

/**
 * @param  {string} project
 * @param  {string} owner
 * @param  {string} repository
 * @param  {string} serverIds
 * @param  {string} userIds
 * @return {Promise}
 */
export function add(data) {
  return request({
    url: '/project/add',
    method: 'post',
    data
  })
}

/**
 * @param  {string} project
 * @param  {string} owner
 * @param  {string} repository
 * @return {Promise}
 */
export function edit(data) {
  return request({
    url: '/project/edit',
    method: 'post',
    data
  })
}

export function remove(id) {
  return request({
    url: '/project/remove',
    method: 'delete',
    data: { id }
  })
}

export function addServer(data) {
  return request({
    url: '/project/addServer',
    method: 'post',
    data
  })
}

export function addUser(data) {
  return request({
    url: '/project/addUser',
    method: 'post',
    data
  })
}

export function removeProjectUser(projectUserId) {
  return request({
    url: '/project/removeProjectUser',
    method: 'delete',
    data: {
      projectUserId
    }
  })
}

export function removeProjectServer(projectServerId) {
  return request({
    url: '/project/removeProjectServer',
    method: 'delete',
    data: {
      projectServerId
    }
  })
}
