import request from '@/utils/request'

/**
 * @return {Promise}
 */
export function getList({ page, rows }, projectName) {
  return request({
    url: '/project/getList',
    method: 'get',
    params: { page, rows, projectName }
  })
}

/**
 * @return {Promise}
 */
export function getTotal(projectName) {
  return request({
    url: '/project/getTotal',
    method: 'get',
    params: { projectName }
  })
}

/**
 * @return {Promise}
 */
export function getRemoteBranchList(url) {
  return request({
    url: '/project/getRemoteBranchList',
    method: 'get',
    timeout: 0,
    params: { url }
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

export function removeUser(projectUserId) {
  return request({
    url: '/project/removeUser',
    method: 'delete',
    data: {
      projectUserId
    }
  })
}

export function removeServer(projectServerId) {
  return request({
    url: '/project/removeServer',
    method: 'delete',
    data: {
      projectServerId
    }
  })
}

export function addTask(data) {
  return request({
    url: '/project/addTask',
    method: 'post',
    data
  })
}

export function editTask(data) {
  return request({
    url: '/project/editTask',
    method: 'post',
    data
  })
}

export function removeTask(id) {
  return request({
    url: '/project/removeTask',
    method: 'post',
    data: { id }
  })
}

/**
 * @param  {id} id
 * @return {Promise}
 */
export function getTaskList({ page, rows }, id) {
  return request({
    url: '/project/getTaskList',
    method: 'get',
    params: { page, rows, id }
  })
}
