import Axios from './axios'
import { Request, Pagination, ID, Total } from './types'

export class CrontabServerData {
  public datagram!: {
    detail: {
      id: number
      crontabId: number
      serverId: number
      serverName: string
      serverIP: string
      serverPort: number
      serverOwner: string
      serverDescription: string
      insertTime: string
      updateTime: string
    }
  }
}

/**
 * @return {Promise}
 */
export function getList({ page, rows }, command) {
  return Axios.request({
    url: '/crontab/getList',
    method: 'get',
    params: { page, rows, command },
  })
}

/**
 * @return {Promise}
 */
export function getTotal(command) {
  return Axios.request({
    url: '/crontab/getTotal',
    method: 'get',
    params: { command },
  })
}

export class CrontabsInRemoteServer extends Request {
  readonly url = '/crontab/getCrontabsInRemoteServer'
  readonly method = 'get'
  public param: {
    serverId: number
  }
  public datagram!: {
    list: string[]
  }
  constructor(param: CrontabsInRemoteServer['param']) {
    super()
    this.param = param
  }
}

export class CrontabServerList extends Request {
  readonly url = '/crontab/getBindServerList'
  readonly method = 'get'
  public param: ID
  public datagram!: {
    list: CrontabServerData['datagram']['detail'][]
  }
  constructor(param: CrontabServerList['param']) {
    super()
    this.param = param
  }
}

/**
 * @param  {id} id
 * @return {Promise}
 */
export function getBindServerList(id) {
  return Axios.request({
    url: '/crontab/getBindServerList',
    method: 'get',
    params: { id },
  })
}

export function add(data) {
  return Axios.request({
    url: '/crontab/add',
    method: 'post',
    data,
  })
}

export function edit(data) {
  return Axios.request({
    url: '/crontab/edit',
    method: 'put',
    data,
  })
}

export function remove(data) {
  return Axios.request({
    url: '/crontab/remove',
    method: 'delete',
    data,
  })
}

export class CrontabImport extends Request {
  readonly url = '/crontab/import'
  readonly method = 'post'
  public param: {
    serverId: number
    commands: string[]
  }
  public datagram!: {
    list: string[]
  }
  constructor(param: CrontabImport['param']) {
    super()
    this.param = param
  }
}

export class CrontabServerAdd extends Request {
  readonly url = '/crontab/addServer'
  readonly method = 'post'
  public param: {
    crontabId: number
    serverIds: number[]
  }
  constructor(param: CrontabServerAdd['param']) {
    super()
    this.param = param
  }
}

export class CrontabServerRemove extends Request {
  readonly url = '/crontab/removeCrontabServer'
  readonly method = 'delete'
  public param: {
    crontabServerId: number
    crontabId: number
    serverId: number
  }
  constructor(param: CrontabServerRemove['param']) {
    super()
    this.param = param
  }
}
