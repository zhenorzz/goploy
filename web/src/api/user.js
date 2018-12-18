import request from '@/utils/request';
import store from '@/store';

/**
 */
export function getInfo() {
  const token = store.getters.token;
  return request({
    url: '/user/info',
    method: 'get',
    params: {token},
  });
}
