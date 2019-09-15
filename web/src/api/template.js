import request from '@/utils/request'

/**
 * @return {Promise}
 */
export function getList(pagination) {
  return request({
    url: '/template/getList',
    method: 'get',
    params: { ...pagination }
  })
}

/**
 * @return {Promise}
 */
export function getOption() {
  return request({
    url: '/template/getOption',
    method: 'get'
  })
}

export function add(data) {
  return request({
    url: '/template/add',
    method: 'post',
    data
  })
}

export function edit(data) {
  return request({
    url: '/template/edit',
    method: 'post',
    data
  })
}

export function remove(id) {
  return request({
    url: '/template/remove',
    method: 'delete',
    data: { id }
  })
}

export function removePackage(templateId, filename) {
  return request({
    url: '/template/removePackage',
    method: 'delete',
    data: { templateId, filename }
  })
}
