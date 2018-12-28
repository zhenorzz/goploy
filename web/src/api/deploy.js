import request from '@/utils/request';

/**
 * @return {Promise}
 */
export function get() {
  return request({
    url: '/deploy/get',
    method: 'get',
    params: {},
  });
}

/**
 * @param  {int}    projectId
 * @param  {string} branch
 * @param  {string} commit
 * @param  {string} commitSha
 * @param  {int}    type
 * @return {Promise}
 */
export function add(projectId, branch, commit, commitSha, type) {
  return request({
    url: '/deploy/add',
    method: 'post',
    data: {
      projectId: projectId,
      branch: branch,
      commit: commit,
      commitSha: commitSha,
      type: type,
    },
  });
}
