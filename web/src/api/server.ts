import { HttpResponse, Request, Pagination, Total } from './types'

export class ServerData {
  public datagram!: {
    detail: {
      id: number
      label: string
      name: string
      ip: string
      port: number
      owner: string
      path: string
      password: string
      namespaceId: number
      description: string
      state: number
      insertTime: string
      updateTime: string
    }
  }
}

export class ServerList extends Request {
  readonly url = '/server/getList'
  readonly method = 'get'

  public pagination: Pagination

  public datagram!: {
    list: ServerData['datagram']['detail'][]
  }
  constructor(pagination: Pagination) {
    super()
    this.pagination = pagination
    this.param = { ...pagination }
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
  public datagram!: Total
}

export class ServerPublicKey extends Request {
  readonly url = '/server/getPublicKey'
  readonly method = 'get'

  public param: {
    path: string
  }

  public datagram!: {
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

  public datagram!: {
    list: ServerData['datagram']['detail'][]
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
