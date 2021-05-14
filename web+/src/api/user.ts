import Axios from './axios'
import { HttpResponse, Pagination, Total, ID } from './types'

type LoginParams = {
  account: string
  password: string
}

type LoginResp = {
  namespaceList: Array<{ id: number; name: string }>
}

export function login(data: LoginParams): Promise<HttpResponse<LoginResp>> {
  return Axios.request({
    url: '/user/login',
    method: 'post',
    data,
  })
}

export function getInfo(): Promise<HttpResponse> {
  return Axios.request({
    url: '/user/info',
    method: 'get',
  })
}

type listFilterParams = Record<string, unknown>
export function getList({
  page,
  rows,
}: Pagination | listFilterParams): Promise<HttpResponse> {
  return Axios.request({
    url: '/user/getList',
    method: 'get',
    params: { page, rows },
  })
}

export function getTotal(): Promise<HttpResponse<Total>> {
  return Axios.request({
    url: '/user/getTotal',
    method: 'get',
    params: {},
  })
}

export function getOption(): Promise<HttpResponse> {
  return Axios.request({
    url: '/user/getOption',
    method: 'get',
  })
}

export function add(data): Promise<HttpResponse<ID>> {
  return Axios.request({
    url: '/user/add',
    method: 'post',
    data,
  })
}

export function edit(data): Promise<HttpResponse<never>> {
  return Axios.request({
    url: '/user/edit',
    method: 'put',
    data,
  })
}

export function remove(id: number): Promise<HttpResponse<never>> {
  return Axios.request({
    url: '/user/remove',
    method: 'delete',
    data: { id },
  })
}

export function changePassword(
  oldPwd: string,
  newPwd: string
): Promise<HttpResponse<never>> {
  return Axios.request({
    url: '/user/changePassword',
    method: 'put',
    data: {
      oldPwd,
      newPwd,
    },
  })
}
