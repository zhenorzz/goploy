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
