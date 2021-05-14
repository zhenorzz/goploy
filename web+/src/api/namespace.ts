import Axios from './axios'

/**
 * @return {Promise}
 */
export function getList({ page, rows }) {
  return Axios.request({
    url: '/namespace/getList',
    method: 'get',
    params: { page, rows },
  })
}

/**
 * @return {Promise}
 */
export function getTotal() {
  return Axios.request({
    url: '/namespace/getTotal',
    method: 'get',
    params: {},
  })
}

/**
 * @return {Promise}
 */
export function getUserOption() {
  return Axios.request({
    url: '/namespace/getUserOption',
    method: 'get',
  })
}

/**
 * @param  {id} id
 * @return {Promise}
 */
export function getBindUserList(id) {
  return Axios.request({
    url: '/namespace/getBindUserList',
    method: 'get',
    params: { id },
  })
}

export function add(data) {
  return Axios.request({
    url: '/namespace/add',
    method: 'post',
    data,
  })
}

export function edit(data) {
  return Axios.request({
    url: '/namespace/edit',
    method: 'put',
    data,
  })
}

export function addUser(data) {
  return Axios.request({
    url: '/namespace/addUser',
    method: 'post',
    data,
  })
}

export function removeUser(namespaceUserId) {
  return Axios.request({
    url: '/namespace/removeUser',
    method: 'delete',
    data: {
      namespaceUserId,
    },
  })
}
