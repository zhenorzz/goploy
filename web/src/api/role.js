import request from '@/utils/request'

export function getOption(params) {
  return request({
    url: '/role/getOption',
    method: 'get',
    params
  })
}
