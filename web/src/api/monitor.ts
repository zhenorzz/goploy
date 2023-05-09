import { Request, ID } from './types'

export interface MonitorData {
  id: number
  namespaceId: number
  name: string
  type: number
  target: any
  second: number
  times: number
  silentCycle: number
  notifyType: number
  notifyTarget: string
  description: string
  errorContent: string
  successServerId: number
  successScript: string
  failServerId: number
  failScript: string
  state: number
  insertTime: string
  updateTime: string
}

export class MonitorList extends Request {
  readonly url = '/monitor/getList'
  readonly method = 'get'

  public declare datagram: {
    list: MonitorData[]
  }
}

export class MonitorAdd extends Request {
  readonly url = '/monitor/add'
  readonly method = 'post'
  public param: {
    name: string
    type: number
    target: string
    second: number
    times: number
    silentCycle: number
    notifyType: number
    notifyTarget: string
    description: string
  }
  constructor(param: MonitorAdd['param']) {
    super()
    this.param = param
  }
}

export class MonitorEdit extends Request {
  readonly url = '/monitor/edit'
  readonly method = 'put'
  public param: {
    id: number
    name: string
    type: number
    target: string
    second: number
    times: number
    silentCycle: number
    notifyType: number
    notifyTarget: string
    description: string
  }
  constructor(param: MonitorEdit['param']) {
    super()
    this.param = param
  }
}

export class MonitorRemove extends Request {
  readonly url = '/monitor/remove'
  readonly method = 'delete'
  public param: ID
  constructor(param: ID) {
    super()
    this.param = param
  }
}

export class MonitorCheck extends Request {
  readonly url = '/monitor/check'
  readonly method = 'post'
  readonly timeout = 100000
  public param: {
    type: number
    target: string
    successServerId: number
    successScript: string
    failServerId: number
    failScript: string
  }
  constructor(param: MonitorCheck['param']) {
    super()
    this.param = param
  }
}

export class MonitorToggle extends Request {
  readonly url = '/monitor/toggle'
  readonly method = 'put'
  public param: {
    id: number
    state: number
  }
  constructor(param: MonitorToggle['param']) {
    super()
    this.param = param
  }
}
