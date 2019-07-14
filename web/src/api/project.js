import request from '@/utils/request'

/**
 * @return {Promise}
 */
export function getList() {
  return request({
    url: '/project/getList',
    method: 'get',
    params: {}
  })
}

/**
 * @param  {id} id
 * @return {Promise}
 */
export function getDetail(id) {
  return request({
    url: '/project/getDetail',
    method: 'get',
    params: { id }
  })
}

/**
 * @param  {id} id
 * @return {Promise}
 */
export function create(id) {
  return request({
    url: '/project/create',
    method: 'post',
    data: { id }
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
