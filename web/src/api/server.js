import request from '@/utils/request'

/**
 * @return {Promise}
 */
export function getList({ page, rows }) {
  return request({
    url: '/server/getList',
    method: 'get',
    params: { page, rows }
  })
}

/**
 * @return {Promise}
 */
export function getTotal() {
  return request({
    url: '/server/getTotal',
    method: 'get',
    params: { }
  })
}

/**
 * @return {Promise}
 */
export function getOption() {
  return request({
    url: '/server/getOption',
    method: 'get'
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

export function check(data) {
  return request({
    timeout: 100000,
    url: '/server/check',
    method: 'post',
    data
  })
}

export function remove(id) {
  return request({
    url: '/server/remove',
    method: 'delete',
    data: { id }
  })
}
