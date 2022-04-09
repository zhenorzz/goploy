import { Request, ID } from './types'

export interface NamespaceData {
  id: number
  name: string
  userId: number
  insertTime: string
  updateTime: string
}

export interface NamespaceUserData {
  id: number
  namespaceId: number
  namespaceName: string
  userId: number
  userName: string
  roleId: number
  roleName: string
  insertTime: string
  updateTime: string
}

export class NamespaceList extends Request {
  readonly url = '/namespace/getList'
  readonly method = 'get'

  public declare datagram: {
    list: NamespaceData[]
  }
}

export class NamespaceUserOption extends Request {
  readonly url = '/namespace/getUserOption'
  readonly method = 'get'

  public declare datagram: {
    list: NamespaceUserData[]
  }
}

export class NamespaceOption extends Request {
  readonly url = '/namespace/getOption'
  readonly method = 'get'

  public declare datagram: {
    list: NamespaceUserData[]
  }
}

export class NamespaceUserList extends Request {
  readonly url = '/namespace/getBindUserList'
  readonly method = 'get'
  public param: ID
  public declare datagram: {
    list: NamespaceUserData[]
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
  public declare datagram: ID
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
    roleId: number
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
