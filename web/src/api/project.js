import request from '@/utils/request';

/**
 * @return {Promise}
 */
export function get() {
  return request({
    url: '/project/get',
    method: 'get',
    params: {},
  });
}

/**
 * @param  {id} id
 * @return {Promise}
 */
export function create(id) {
  return request({
    url: '/project/create',
    method: 'get',
    params: {id},
  });
}

/**
 * @return {Promise}
 */
export function branch() {
  return request({
    url: '/project/branch',
    method: 'get',
    params: {},
  });
}

/**
 * @return {Promise}
 */
export function commit() {
  return request({
    url: '/project/commit',
    method: 'get',
    params: {},
  });
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
      sha: sha,
    },
  });
}

/**
 * @param  {string} project
 * @param  {string} owner
 * @param  {string} repository
 * @return {Promise}
 */
export function add(project, owner, repository) {
  return request({
    url: '/project/add',
    method: 'post',
    data: {
      project: project,
      owner: owner,
      repository: repository,
    },
  });
}
