import request from '@/utils/request'

/**
 * @return {Promise}
 */
export function getList(groupId, projectName) {
  return request({
    url: '/deploy/getList',
    method: 'get',
    params: { groupId, projectName }
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
export function getPreview(projectId) {
  return request({
    url: '/deploy/getPreview',
    method: 'get',
    params: {
      projectId
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
    params: { id }
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
