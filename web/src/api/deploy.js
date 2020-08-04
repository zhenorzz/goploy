import request from '@/utils/request'

/**
 * @return {Promise}
 */
export function getList(projectName) {
  return request({
    url: '/deploy/getList',
    method: 'get',
    params: { projectName }
  })
}

/**
 * @param  {int}    id
 * @return {Promise}
 */
export function getDetail(lastPublishToken) {
  return request({
    url: '/deploy/getDetail',
    method: 'get',
    params: {
      lastPublishToken
    }
  })
}

/**
 * @param  {int}    id
 * @return {Promise}
 */
export function getPreview({ page, rows }, params) {
  return request({
    url: '/deploy/getPreview',
    method: 'get',
    params: {
      page, rows, ...params
    }
  })
}

/**
 * @param  {int}    id
 * @return {Promise}
 */
export function getCommitList(id) {
  return request({
    url: '/deploy/getCommitList',
    method: 'get',
    params: { id },
    timeout: 0
  })
}

/**
 * @param  {int}    id
 * @return {Promise}
 */
export function publish(projectId, commit) {
  return request({
    url: '/deploy/publish',
    method: 'post',
    data: { projectId, commit }
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
