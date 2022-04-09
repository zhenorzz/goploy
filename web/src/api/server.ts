import { HttpResponse, Request, Pagination, Total, ID } from './types'

export interface ServerData {
  id: number
  label: string
  name: string
  ip: string
  port: number
  owner: string
  path: string
  password: string
  jumpIP: string
  jumpPort: number
  jumpOwner: string
  jumpPath: string
  jumpPassword: string
  namespaceId: number
  description: string
  state: number
  insertTime: string
  updateTime: string
}

export interface ServerMonitorData {
  id: number
  serverId: number
  item: string
  formula: string
  operator: string
  value: string
  groupCycle: number
  lastCycle: number
  silentCycle: number
  startTime: string
  endTime: string
  notifyType: number
  notifyTarget: string
  insertTime: string
  updateTime: string
}

export class ServerList extends Request {
  readonly url = '/server/getList'
  readonly method = 'get'

  public declare datagram: {
    list: ServerData[]
  }

  public request(): Promise<HttpResponse<this['datagram']>> {
    return super.request().then((response) => {
      response.data.list = response.data.list.map((element) => {
        element.label =
          element.name +
          (element.description.length > 0
            ? '(' + element.description + ')'
            : '')
        return element
      })
      return response
    })
  }
}

export class ServerTotal extends Request {
  readonly url = '/server/getTotal'
  readonly method = 'get'
  public declare datagram: Total
}

export class ServerPublicKey extends Request {
  readonly url = '/server/getPublicKey'
  readonly method = 'get'

  public param: {
    path: string
  }

  public declare datagram: {
    key: string
  }

  constructor(param: ServerPublicKey['param']) {
    super()
    this.param = param
  }
}

export class ServerOption extends Request {
  readonly url = '/server/getOption'
  readonly method = 'get'

  public declare datagram: {
    list: ServerData[]
  }
  public request(): Promise<HttpResponse<this['datagram']>> {
    return super.request().then((response) => {
      response.data.list = response.data.list.map((element) => {
        element.label =
          element.name +
          (element.description.length > 0
            ? '(' + element.description + ')'
            : '')
        return element
      })
      return response
    })
  }
}

export class ServerAdd extends Request {
  readonly url = '/server/add'
  readonly method = 'post'
  public param: {
    namespaceId: number
    name: string
    ip: string
    port: number
    owner: string
    path: string
    password: string
    jumpIP: string
    jumpPort: number
    jumpOwner: string
    jumpPath: string
    jumpPassword: string
    description: string
  }
  constructor(param: ServerAdd['param']) {
    super()
    this.param = param
  }
}

export class ServerEdit extends Request {
  readonly url = '/server/edit'
  readonly method = 'put'
  public param: {
    id: number
    namespaceId: number
    name: string
    ip: string
    port: number
    owner: string
    path: string
    password: string
    jumpIP: string
    jumpPort: number
    jumpOwner: string
    jumpPath: string
    jumpPassword: string
    description: string
  }
  constructor(param: ServerEdit['param']) {
    super()
    this.param = param
  }
}

export class ServerCheck extends Request {
  readonly url = '/server/check'
  readonly method = 'post'
  public param: {
    ip: string
    port: number
    owner: string
    path: string
    password: string
    jumpIP: string
    jumpPort: number
    jumpOwner: string
    jumpPath: string
    jumpPassword: string
  }
  constructor(param: ServerCheck['param']) {
    super()
    this.param = param
  }
}

export class ServerToggle extends Request {
  readonly url = '/server/toggle'
  readonly method = 'put'
  public param: {
    id: number
    state: number
  }
  constructor(param: ServerToggle['param']) {
    super()
    this.param = param
  }
}

export class ServerInstallAgent extends Request {
  readonly url = '/server/installAgent'
  readonly method = 'post'
  public param: {
    ids: number[]
    installPath: string
    tool: string
    reportURL: string
    webPort: string
  }
  constructor(param: ServerInstallAgent['param']) {
    super()
    this.param = param
  }
}

export class ServerReport extends Request {
  readonly url = '/server/report'
  readonly method = 'get'
  public param: {
    serverId: number
    type: number
    datetimeRange: string
  }
  constructor(param: ServerReport['param']) {
    super()
    this.param = param
  }
}

export class ServerMonitorList extends Request {
  readonly url = '/server/getAllMonitor'
  readonly method = 'get'

  public param: {
    serverId: number
  }

  public declare datagram: {
    list: ServerMonitorData[]
  }

  constructor(param: ServerMonitorList['param']) {
    super()
    this.param = param
  }
}

export class ServerMonitorAdd extends Request {
  readonly url = '/server/addMonitor'
  readonly method = 'post'
  public param: {
    serverId: number
    item: string
    formula: string
    operator: string
    value: string
    groupCycle: number
    lastCycle: number
    silentCycle: number
    startTime: string
    endTime: string
    notifyType: number
    notifyTarget: string
  }
  constructor(param: ServerMonitorAdd['param']) {
    super()
    this.param = param
  }
}

export class ServerMonitorEdit extends Request {
  readonly url = '/server/editMonitor'
  readonly method = 'put'
  public param: {
    id: number
    serverId: number
    item: string
    formula: string
    operator: string
    value: string
    groupCycle: number
    lastCycle: number
    silentCycle: number
    startTime: string
    endTime: string
    notifyType: number
    notifyTarget: string
  }
  constructor(param: ServerMonitorEdit['param']) {
    super()
    this.param = param
  }
}

export class ServerMonitorDelete extends Request {
  readonly url = '/server/deleteMonitor'
  readonly method = 'delete'
  public param: ID
  constructor(param: ServerMonitorDelete['param']) {
    super()
    this.param = param
  }
}
