import request from '@/utils/request'

/**
 * @return {Promise}
 */
export function get() {
  return request({
    url: '/server/get',
    method: 'get',
    params: {}
  })
}

/**
 * @param  {string} name
 * @param  {string} ip
 * @param  {string} path
 * @return {Promise}
 */
export function add(name, ip, path) {
  return request({
    url: '/server/add',
    method: 'post',
    data: {
      name: name,
      ip: ip,
      path: path
    }
  })
}
