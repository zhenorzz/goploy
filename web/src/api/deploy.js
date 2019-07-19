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
export function publish(id) {
  return request({
    url: '/deploy/publish',
    method: 'post',
    data: { id }
  })
}

/**
 * @param  {int}    projectId
 * @param  {string} branch
 * @param  {string} commit
 * @param  {string} commitSha
 * @param  {int}    serverId
 * @param  {int}    type
 * @return {Promise}
 */
export function add(projectId, branch, commit, commitSha, serverId, type) {
  return request({
    url: '/deploy/add',
    method: 'post',
    data: {
      projectId: projectId,
      branch: branch,
      commit: commit,
      commitSha: commitSha,
      serverId: serverId,
      type: type
    }
  })
}
