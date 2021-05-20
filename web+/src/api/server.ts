import Axios from './axios'
import { Request, Pagination, ID, Total } from './types'

export class ServerData {
  public datagram!: {
    detail: {
      id: number
      name: string
      ip: string
      port: number
      owner: string
      path: string
      password: string
      namespaceId: number
      description: string
      insertTime: string
      updateTime: string
    }
  }
}

/**
 * @return {Promise}
 */
export function getList({ page, rows }) {
  return Axios.request({
    url: '/server/getList',
    method: 'get',
    params: { page, rows },
  })
}

/**
 * @return {Promise}
 */
export function getTotal() {
  return Axios.request({
    url: '/server/getTotal',
    method: 'get',
    params: {},
  })
}

/**
 * @return {Promise}
 */
export function getPublicKey(path) {
  return Axios.request({
    url: '/server/getPublicKey',
    method: 'get',
    params: { path },
  })
}

export class ServerOption extends Request {
  readonly url = '/server/getOption'
  readonly method = 'get'

  public datagram!: {
    list: ServerData['datagram']['detail'][]
  }
}

/**
 * @return {Promise}
 */
export function getOption() {
  return Axios.request({
    url: '/server/getOption',
    method: 'get',
  })
}

export function add(data) {
  return Axios.request({
    url: '/server/add',
    method: 'post',
    data,
  })
}

export function edit(data) {
  return Axios.request({
    url: '/server/edit',
    method: 'put',
    data,
  })
}

export function check(data) {
  return Axios.request({
    timeout: 100000,
    url: '/server/check',
    method: 'post',
    data,
  })
}

export function remove(id) {
  return Axios.request({
    url: '/server/remove',
    method: 'delete',
    data: { id },
  })
}
