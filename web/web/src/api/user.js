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
