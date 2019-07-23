import request from '@/utils/request'

export function getOption(params) {
  return request({
    url: '/role/getOption',
    method: 'get',
    params
  })
}

export function getPermissionList() {
  return request({
    url: '/role/getPermissionList',
    method: 'get'
  })
}

export function edit(id, name, remark, permissionList) {
  return request({
    url: '/role/edit',
    method: 'post',
    data: {
      id, name, remark, permissionList
    }
  })
}
