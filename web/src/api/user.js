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
 * @param  {object} pagination
 * @return {Promise}
 */
export function get(pagination) {
  return request({
    url: '/user/get',
    method: 'get',
    params: {
      ...pagination,
    },
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

/**
 * @param  {string} account
 * @param  {string} oldPwd
 * @param  {string} newPwd
 * @return {Promise}
 */
export function changePassword(account, oldPwd, newPwd) {
  return request({
    url: '/user/changePassword',
    method: 'post',
    data: {
      account,
      oldPwd,
      newPwd,
    },
  });
}
