import Axios from './axios'

/**
 * @return {Promise}
 */
export function getList({ page, rows }, projectName) {
  return Axios.request({
    url: '/project/getList',
    method: 'get',
    params: { page, rows, projectName },
  })
}

/**
 * @return {Promise}
 */
export function getTotal(projectName) {
  return Axios.request({
    url: '/project/getTotal',
    method: 'get',
    params: { projectName },
  })
}

/**
 * @return {Promise}
 */
export function getRemoteBranchList(url) {
  return Axios.request({
    url: '/project/getRemoteBranchList',
    method: 'get',
    timeout: 0,
    params: { url },
  })
}

/**
 * @param  {id} id
 * @return {Promise}
 */
export function getBindServerList(id) {
  return Axios.request({
    url: '/project/getBindServerList',
    method: 'get',
    params: { id },
  })
}

/**
 * @param  {id} id
 * @return {Promise}
 */
export function getBindUserList(id) {
  return Axios.request({
    url: '/project/getBindUserList',
    method: 'get',
    params: { id },
  })
}

/**
 * @param  {id} id
 * @return {Promise}
 */
export function getProjectFileList(id) {
  return Axios.request({
    url: '/project/getProjectFileList',
    method: 'get',
    params: { id },
  })
}

/**
 * @param  {id} id
 * @return {Promise}
 */
export function getProjectFileContent(id) {
  return Axios.request({
    url: '/project/getProjectFileContent',
    method: 'get',
    params: { id },
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
  return Axios.request({
    url: '/project/add',
    method: 'post',
    data,
  })
}

/**
 * @param  {string} project
 * @param  {string} owner
 * @param  {string} repository
 * @return {Promise}
 */
export function edit(data) {
  return Axios.request({
    url: '/project/edit',
    method: 'put',
    data,
  })
}

export function setAutoDeploy(data) {
  return Axios.request({
    url: '/project/setAutoDeploy',
    method: 'put',
    data,
  })
}

export function remove(id) {
  return Axios.request({
    url: '/project/remove',
    method: 'delete',
    data: { id },
  })
}

export function addUser(data) {
  return Axios.request({
    url: '/project/addUser',
    method: 'post',
    data,
  })
}

export function removeUser(projectUserId) {
  return Axios.request({
    url: '/project/removeUser',
    method: 'delete',
    data: {
      projectUserId,
    },
  })
}

export function addServer(data) {
  return Axios.request({
    url: '/project/addServer',
    method: 'post',
    data,
  })
}

export function removeServer(projectServerId) {
  return Axios.request({
    url: '/project/removeServer',
    method: 'delete',
    data: {
      projectServerId,
    },
  })
}

export function addFile(data) {
  return Axios.request({
    url: '/project/addFile',
    method: 'post',
    data,
  })
}

export function editFile(data) {
  return Axios.request({
    url: '/project/editFile',
    method: 'put',
    data,
  })
}

export function removeFile(projectFileId) {
  return Axios.request({
    url: '/project/removeFile',
    method: 'delete',
    data: {
      projectFileId,
    },
  })
}

export function addTask(data) {
  return Axios.request({
    url: '/project/addTask',
    method: 'post',
    data,
  })
}

export function removeTask(id) {
  return Axios.request({
    url: '/project/removeTask',
    method: 'delete',
    data: { id },
  })
}

/**
 * @param  {id} id
 * @return {Promise}
 */
export function getTaskList({ page, rows }, id) {
  return Axios.request({
    url: '/project/getTaskList',
    method: 'get',
    params: { page, rows, id },
  })
}

/**
 * @param  {id} id
 * @return {Promise}
 */
export function getReviewList({ page, rows }, id) {
  return Axios.request({
    url: '/project/getReviewList',
    method: 'get',
    params: { page, rows, id },
  })
}
