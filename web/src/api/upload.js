import request from '@/utils/request'

export function getOss() {
  return request({
    url: '/upload/assumeRole',
    method: 'get'
  })
}
