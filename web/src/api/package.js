import request from '@/utils/request'

/**
 * @return {Promise}
 */
export function getList() {
  return request({
    url: '/package/getList',
    method: 'get',
    params: {}
  })
}

/**
 * @return {Promise}
 */
export function getOption() {
  return request({
    url: '/package/getOption',
    method: 'get'
  })
}
