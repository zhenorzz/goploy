import request from '@/utils/request'

/**
 * @return {Promise}
 */
export function getList(pagination) {
  return request({
    url: '/package/getList',
    method: 'get',
    params: { ...pagination }
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
