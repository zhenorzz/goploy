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
