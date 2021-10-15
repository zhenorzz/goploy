import { Request, Pagination, ID, Total } from './types'

export class MonitorData {
  public datagram!: {
    detail: {
      id: number
      namespaceId: number
      name: string
      url: string
      second: number
      times: number
      notifyType: number
      notifyTarget: string
      notifyTimes: number
      description: string
      errorContent: string
      state: number
      insertTime: string
      updateTime: string
    }
  }
}

export class MonitorList extends Request {
  readonly url = '/monitor/getList'
  readonly method = 'get'
  public pagination: Pagination

  public datagram!: {
    list: MonitorData['datagram']['detail'][]
  }
  constructor(pagination: Pagination) {
    super()
    this.pagination = pagination
    this.param = { ...pagination }
  }
}

export class MonitorTotal extends Request {
  readonly url = '/monitor/getTotal'
  readonly method = 'get'

  public datagram!: Total
}

export class MonitorAdd extends Request {
  readonly url = '/monitor/add'
  readonly method = 'post'
  public param: {
    name: string
    url: string
    second: number
    times: number
    notifyType: number
    notifyTarget: string
    notifyTimes: number
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
    url: string
    second: number
    times: number
    notifyType: number
    notifyTarget: string
    notifyTimes: number
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
    url: string
  }
  constructor(param: MonitorCheck['param']) {
    super()
    this.param = param
  }
}

export class MonitorToggle extends Request {
  readonly url = '/monitor/toggle'
  readonly method = 'put'
  public param: ID
  constructor(param: ID) {
    super()
    this.param = param
  }
}
