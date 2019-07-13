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

export function add(data) {
  return request({
    url: '/server/add',
    method: 'post',
    data
  })
}

export function edit(data) {
  return request({
    url: '/server/edit',
    method: 'post',
    data
  })
}
