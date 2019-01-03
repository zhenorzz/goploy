import request from '@/utils/request';

/**
 * @param  {string} date
 * @return {Promise}
 */
export function get(date) {
  return request({
    url: '/index/get',
    method: 'get',
    params: {
      date,
    },
  });
}
