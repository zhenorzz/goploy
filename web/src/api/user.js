import request from '@/utils/request';

/**
 * @return {Promise}
 */
export function getInfo() {
  return request({
    url: '/user/info',
    method: 'get',
    params: {},
  });
}

/**
 * @return {Promise}
 */
export function get() {
  return request({
    url: '/user/get',
    method: 'get',
    params: {},
  });
}

/**
 * @param  {string} account
 * @param  {string} password
 * @param  {string} name
 * @param  {string} email
 * @param  {string} role
 * @return {Promise}
 */
export function add(account, password, name, email, role) {
  return request({
    url: '/user/add',
    method: 'post',
    data: {
      name,
      account,
      password,
      email,
      role,
    },
  });
}
