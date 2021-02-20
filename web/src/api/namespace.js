import request from '@/utils/request'

/**
 * @return {Promise}
 */
export function getList({ page, rows }) {
  return request({
    url: '/namespace/getList',
    method: 'get',
    params: { page, rows }
  })
}

/**
 * @return {Promise}
 */
export function getTotal() {
  return request({
    url: '/namespace/getTotal',
    method: 'get',
    params: { }
  })
}

/**
 * @return {Promise}
 */
export function getUserOption() {
  return request({
    url: '/namespace/getUserOption',
    method: 'get'
  })
}

/**
 * @param  {id} id
 * @return {Promise}
 */
export function getBindUserList(id) {
  return request({
    url: '/namespace/getBindUserList',
    method: 'get',
    params: { id }
  })
}

export function add(data) {
  return request({
    url: '/namespace/add',
    method: 'post',
    data
  })
}

export function edit(data) {
  return request({
    url: '/namespace/edit',
    method: 'put',
    data
  })
}

export function addUser(data) {
  return request({
    url: '/namespace/addUser',
    method: 'post',
    data
  })
}

export function removeUser(namespaceUserId) {
  return request({
    url: '/namespace/removeUser',
    method: 'delete',
    data: {
      namespaceUserId
    }
  })
}
