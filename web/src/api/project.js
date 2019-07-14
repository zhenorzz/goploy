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
 * @param  {id} id
 * @return {Promise}
 */
export function branch(id) {
  return request({
    url: '/project/branch',
    method: 'post',
    data: { id }
  })
}

/**
 * @param  {id} id
 * @param  {string} branch
 * @return {Promise}
 */
export function commit(id, branch) {
  return request({
    url: '/project/commit',
    method: 'post',
    data: { id, branch }
  })
}

/**
 * @param  {string} sha
 * @return {Promise}
 */
export function tree(sha) {
  return request({
    url: '/project/tree',
    method: 'get',
    params: {
      sha: sha
    }
  })
}

/**
 * @param  {string} project
 * @param  {string} owner
 * @param  {string} repository
 * @param  {string} serverIds
 * @return {Promise}
 */
export function add(data) {
  return request({
    url: '/project/add',
    method: 'post',
    data
  })
}
