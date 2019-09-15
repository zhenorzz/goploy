import request from '@/utils/request'

/**
 * @return {Promise}
 */
export function getList(pagination) {
  return request({
    url: '/group/getList',
    method: 'get',
    params: { ...pagination }
  })
}

/**
 * @return {Promise}
 */
export function getOption() {
  return request({
    url: '/group/getOption',
    method: 'get'
  })
}

export function add(data) {
  return request({
    url: '/group/add',
    method: 'post',
    data
  })
}

export function edit(data) {
  return request({
    url: '/group/edit',
    method: 'post',
    data
  })
}

export function remove(id) {
  return request({
    url: '/group/remove',
    method: 'delete',
    data: { id }
  })
}
