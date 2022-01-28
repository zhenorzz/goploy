import { Request, Pagination, ID, Total } from './types'

export class NamespaceData {
  public datagram!: {
    id: number
    name: string
    userId: number
    role: string
    insertTime: string
    updateTime: string
  }
}

export class NamespaceUserData {
  public datagram!: {
    id: number
    namespaceId: number
    namespaceName: string
    userId: number
    userName: string
    role: string
    insertTime: string
    updateTime: string
  }
}

export class NamespaceList extends Request {
  readonly url = '/namespace/getList'
  readonly method = 'get'

  public pagination: Pagination

  public datagram!: {
    list: NamespaceData['datagram'][]
  }
  constructor(pagination: Pagination) {
    super()
    this.pagination = pagination
    this.param = { ...pagination }
  }
}

export class NamespaceTotal extends Request {
  readonly url = '/namespace/getTotal'
  readonly method = 'get'

  public datagram!: Total
}

export class NamespaceUserOption extends Request {
  readonly url = '/namespace/getUserOption'
  readonly method = 'get'

  public datagram!: {
    list: NamespaceUserData['datagram'][]
  }
}

export class NamespaceOption extends Request {
  readonly url = '/namespace/getOption'
  readonly method = 'get'

  public datagram!: {
    list: NamespaceUserData['datagram'][]
  }
}

export class NamespaceUserList extends Request {
  readonly url = '/namespace/getBindUserList'
  readonly method = 'get'
  public param: ID
  public datagram!: {
    list: NamespaceUserData['datagram'][]
  }
  constructor(param: NamespaceUserList['param']) {
    super()
    this.param = param
  }
}

export class NamespaceAdd extends Request {
  readonly url = '/namespace/add'
  readonly method = 'post'
  public param: {
    name: string
  }
  public datagram!: ID
  constructor(param: NamespaceAdd['param']) {
    super()
    this.param = param
  }
}

export class NamespaceEdit extends Request {
  readonly url = '/namespace/edit'
  readonly method = 'put'
  public param: {
    id: number
    name: string
  }
  constructor(param: NamespaceEdit['param']) {
    super()
    this.param = param
  }
}

export class NamespaceUserAdd extends Request {
  readonly url = '/namespace/addUser'
  readonly method = 'post'
  public param: {
    namespaceId: number
    userIds: number[]
    role: string
  }
  constructor(param: NamespaceUserAdd['param']) {
    super()
    this.param = param
  }
}

export class NamespaceUserRemove extends Request {
  readonly url = '/namespace/removeUser'
  readonly method = 'delete'
  public param: {
    namespaceUserId: number
  }
  constructor(param: NamespaceUserRemove['param']) {
    super()
    this.param = param
  }
}
