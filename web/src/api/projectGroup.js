import request from '@/utils/request'

/**
 * @return {Promise}
 */
export function getList() {
  return request({
    url: '/projectGroup/getList',
    method: 'get',
    params: {}
  })
}

/**
 * @return {Promise}
 */
export function getOption() {
  return request({
    url: '/projectGroup/getOption',
    method: 'get'
  })
}

export function add(data) {
  return request({
    url: '/projectGroup/add',
    method: 'post',
    data
  })
}

export function edit(data) {
  return request({
    url: '/projectGroup/edit',
    method: 'post',
    data
  })
}

export function remove(id) {
  return request({
    url: '/projectGroup/remove',
    method: 'post',
    data: { id }
  })
}
