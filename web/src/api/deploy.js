import request from '@/utils/request'

/**
 * @param  {string}    projectName
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
 * @param  {string}    lastPublishToken
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
 * @param  {object}  pagination
 * @param  {object}  params
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
export function getCommitList(id, branch) {
  return request({
    url: '/deploy/getCommitList',
    method: 'get',
    params: { id, branch },
    timeout: 0
  })
}

/**
 * @param  {int}    id
 * @return {Promise}
 */
export function getBranchList(id) {
  return request({
    url: '/deploy/getBranchList',
    method: 'get',
    params: { id },
    timeout: 0
  })
}

/**
 * @param  {int}    id
 * @return {Promise}
 */
export function getTagList(id) {
  return request({
    url: '/deploy/getTagList',
    method: 'get',
    params: { id },
    timeout: 0
  })
}

/**
 * @param  {int}      projectId
 * @param  {string}   commit
 * @return {Promise}
 */
export function publish(projectId, branch, commit) {
  return request({
    url: '/deploy/publish',
    method: 'post',
    data: { projectId, branch, commit }
  })
}

/**
 * @param  {int}      projectId
 * @return {Promise}
 */
export function resetState(projectId) {
  return request({
    url: '/deploy/resetState',
    method: 'post',
    data: { projectId }
  })
}

/**
 * @param  {int}      projectId
 * @param  {string}   commit
 * @param  {Array}    serverIds
 * @return {Promise}
 */
export function greyPublish(projectId, commit, serverIds) {
  return request({
    url: '/deploy/greyPublish',
    method: 'post',
    data: { projectId, commit, serverIds }
  })
}

/**
 * @param  {int}    projectReviewId
 * @param  {int}    state
 * @return {Promise}
 */
export function review(projectReviewId, state) {
  return request({
    url: '/deploy/review',
    method: 'post',
    data: { projectReviewId, state }
  })
}
