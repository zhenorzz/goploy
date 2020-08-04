import request from '@/utils/request'

/**
 * @return {Promise}
 */
export function getList({ page, rows }) {
  return request({
    url: '/monitor/getList',
    method: 'get',
    params: { page, rows }
  })
}

/**
 * @return {Promise}
 */
export function getTotal() {
  return request({
    url: '/monitor/getTotal',
    method: 'get',
    params: { }
  })
}

export function add(data) {
  return request({
    url: '/monitor/add',
    method: 'post',
    data
  })
}

export function edit(data) {
  return request({
    url: '/monitor/edit',
    method: 'post',
    data
  })
}

export function check(data) {
  return request({
    timeout: 100000,
    url: '/monitor/check',
    method: 'post',
    data
  })
}

export function toggle(id) {
  return request({
    url: '/monitor/toggle',
    method: 'post',
    data: { id }
  })
}

export function remove(id) {
  return request({
    url: '/monitor/remove',
    method: 'delete',
    data: { id }
  })
}
