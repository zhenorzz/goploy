import { Request, ID } from './types'

export interface RoleData {
  id: number
  name: string
  description: string
  insertTime: string
  updateTime: string
}

export interface PermissionData {
  id: number
  pid: number
  name: string
  sort: number
  description: string
  insertTime: string
  updateTime: string
}

export interface RolePermissionData {
  id: number
  roleId: number
  permissionId: number
}

export class RoleList extends Request {
  readonly url = '/role/getList'
  readonly method = 'get'

  public declare datagram: {
    list: RoleData[]
  }
}

export class RoleOption extends Request {
  readonly url = '/role/getOption'
  readonly method = 'get'
  public declare datagram: {
    list: RoleData[]
  }
}

export class RoleAdd extends Request {
  readonly url = '/role/add'
  readonly method = 'post'
  public param: {
    name: string
    description: string
  }
  public declare datagram: ID
  constructor(param: RoleAdd['param']) {
    super()
    this.param = param
  }
}

export class RoleEdit extends Request {
  readonly url = '/role/edit'
  readonly method = 'put'
  public param: {
    id: number
    name: string
    description: string
  }
  constructor(param: RoleEdit['param']) {
    super()
    this.param = param
  }
}

export class RoleRemove extends Request {
  readonly url = '/role/remove'
  readonly method = 'delete'
  public param: {
    id: number
  }
  constructor(param: RoleRemove['param']) {
    super()
    this.param = param
  }
}

export class RolePermissionList extends Request {
  readonly url = '/role/getPermissionList'
  readonly method = 'get'

  public declare datagram: {
    list: PermissionData[]
  }
}

export class RolePermissionBindings extends Request {
  readonly url = '/role/getPermissionBindings'
  readonly method = 'get'
  public param: {
    roleId: number
  }
  public declare datagram: {
    list: RolePermissionData[]
  }
  constructor(param: RolePermissionBindings['param']) {
    super()
    this.param = param
  }
}

export class RolePermissionChange extends Request {
  readonly url = '/role/changePermission'
  readonly method = 'put'
  public param: {
    permissionIds: number[]
    roleId: number
  }
  constructor(param: RolePermissionChange['param']) {
    super()
    this.param = param
  }
}
