import Axios from './axios'

/**
 * @return {Promise}
 */
export function getList({ page, rows }) {
  return Axios.request({
    url: '/monitor/getList',
    method: 'get',
    params: { page, rows },
  })
}

/**
 * @return {Promise}
 */
export function getTotal() {
  return Axios.request({
    url: '/monitor/getTotal',
    method: 'get',
    params: {},
  })
}

export function add(data) {
  return Axios.request({
    url: '/monitor/add',
    method: 'post',
    data,
  })
}

export function edit(data) {
  return Axios.request({
    url: '/monitor/edit',
    method: 'put',
    data,
  })
}

export function check(data) {
  return Axios.request({
    timeout: 100000,
    url: '/monitor/check',
    method: 'post',
    data,
  })
}

export function toggle(id) {
  return Axios.request({
    url: '/monitor/toggle',
    method: 'put',
    data: { id },
  })
}

export function remove(id) {
  return Axios.request({
    url: '/monitor/remove',
    method: 'delete',
    data: { id },
  })
}
