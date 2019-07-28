import request from '@/utils/request'

/**
 * @return {Promise}
 */
export function getList() {
  return request({
    url: '/deploy/getList',
    method: 'get',
    params: {}
  })
}

/**
 * @param  {int}    id
 * @return {Promise}
 */
export function getDetail(id) {
  return request({
    url: '/deploy/getDetail',
    method: 'get',
    params: { id }
  })
}

/**
 * @param  {int}    id
 * @return {Promise}
 */
export function getSyncDetail(gitTraceId) {
  return request({
    url: '/deploy/getSyncDetail',
    method: 'get',
    params: { gitTraceId }
  })
}

/**
 * @param  {int}    id
 * @return {Promise}
 */
export function publish(projectId) {
  return request({
    url: '/deploy/publish',
    method: 'post',
    data: { projectId }
  })
}

/**
 * @param  {int}    id
 * @return {Promise}
 */
export function rollback(projectId, commit) {
  return request({
    url: '/deploy/rollback',
    method: 'post',
    data: { projectId, commit }
  })
}
