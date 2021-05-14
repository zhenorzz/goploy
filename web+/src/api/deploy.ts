import Axios from './axios'
/**
 * @param  {string}    projectName
 * @return {Promise}
 */
export function getList(projectName) {
  return Axios.request({
    url: '/deploy/getList',
    method: 'get',
    params: { projectName },
  })
}

/**
 * @param  {string}    lastPublishToken
 * @return {Promise}
 */
export function getPublishTrace(lastPublishToken) {
  return Axios.request({
    url: '/deploy/getPublishTrace',
    method: 'get',
    params: {
      lastPublishToken,
    },
  })
}

/**
 * @param  {Number}    publish_trace_id
 * @return {Promise}
 */
export function getPublishTraceDetail(publish_trace_id) {
  return Axios.request({
    url: '/deploy/getPublishTraceDetail',
    method: 'get',
    params: {
      publish_trace_id,
    },
    timeout: 0,
  })
}

/**
 * @param  {object}  pagination
 * @param  {object}  params
 * @return {Promise}
 */
export function getPreview({ page, rows }, params) {
  return Axios.request({
    url: '/deploy/getPreview',
    method: 'get',
    params: {
      page,
      rows,
      ...params,
    },
  })
}

/**
 * @param  {int}    id
 * @return {Promise}
 */
export function getCommitList(id, branch) {
  return Axios.request({
    url: '/deploy/getCommitList',
    method: 'get',
    params: { id, branch },
    timeout: 0,
  })
}

/**
 * @param  {int}    id
 * @return {Promise}
 */
export function getBranchList(id) {
  return Axios.request({
    url: '/deploy/getBranchList',
    method: 'get',
    params: { id },
    timeout: 0,
  })
}

/**
 * @param  {int}    id
 * @return {Promise}
 */
export function getTagList(id) {
  return Axios.request({
    url: '/deploy/getTagList',
    method: 'get',
    params: { id },
    timeout: 0,
  })
}

/**
 * @param  {int}      projectId
 * @return {Promise}
 */
export function resetState(projectId) {
  return Axios.request({
    url: '/deploy/resetState',
    method: 'put',
    data: { projectId },
  })
}

/**
 * @param  {int}      projectId
 * @param  {string}   commit
 * @return {Promise}
 */
export function publish(projectId, branch, commit) {
  return Axios.request({
    url: '/deploy/publish',
    method: 'post',
    data: { projectId, branch, commit },
  })
}

/**
 * @param  {string}   token
 * @return {Promise}
 */
export function rebuild(projectId, token) {
  return Axios.request({
    url: '/deploy/rebuild',
    method: 'post',
    data: { projectId, token },
  })
}

/**
 * @param  {int}      projectId
 * @param  {string}   commit
 * @param  {Array}    serverIds
 * @return {Promise}
 */
export function greyPublish(projectId, commit, serverIds) {
  return Axios.request({
    url: '/deploy/greyPublish',
    method: 'post',
    data: { projectId, commit, serverIds },
  })
}

/**
 * @param  {int}    projectReviewId
 * @param  {int}    state
 * @return {Promise}
 */
export function review(projectReviewId, state) {
  return Axios.request({
    url: '/deploy/review',
    method: 'put',
    data: { projectReviewId, state },
  })
}
