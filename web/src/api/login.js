import request from '@/utils/request';

export function login(account, password) {
  return request({
    url: '/user/login',
    method: 'post',
    data: {
      account,
      password,
    },
  });
}

export function logout() {
  return request({
    url: '/user/logout',
    method: 'post',
  });
}
