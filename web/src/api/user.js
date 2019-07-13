import request from '@/utils/request'

export function login(data) {
  return request({
    url: '/user/login',
    method: 'post',
    data
  })
}

export function getInfo() {
  return request({
    url: '/user/info',
    method: 'get'
  })
}

export function logout() {
  return request({
    url: '/user/logout',
    method: 'get'
  })
}

export function isShowPhrase() {
  return request({
    url: '/user/isShowPhrase',
    method: 'get'
  })
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
      ...pagination
    }
  })
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
      role
    }
  })
}

export function changePassword(oldPwd, newPwd) {
  return request({
    url: '/user/changePassword',
    method: 'post',
    data: {
      oldPwd,
      newPwd
    }
  })
}
