import { Request, ID } from './types'

export interface UserData {
  account: string
  contact: string
  id: number
  insertTime: string
  lastLoginTime: string
  name: string
  password: string
  state: number
  superManager: number
  updateTime: string
}

export class Login extends Request {
  readonly url = '/user/login'
  readonly method = 'post'
  public param: {
    account: string
    password: string
  }
  public declare datagram: {
    namespaceList: { id: number; name: string; role_id: number }[]
  }
  constructor(param: Login['param']) {
    super()
    this.param = param
  }
}

export class Info extends Request {
  readonly url = '/user/info'
  readonly method = 'get'
  public declare datagram: {
    userInfo: {
      account: string
      id: number
      name: string
      superManager: number
    }
    namespace: {
      id: number
      permissionIds: number[]
    }
  }
}

export class UserList extends Request {
  readonly url = '/user/getList'
  readonly method = 'get'
  public declare datagram: {
    list: UserData[]
  }
}

export class UserOption extends Request {
  readonly url = '/user/getOption'
  readonly method = 'get'
  public declare datagram: {
    list: UserData[]
  }
}

export class UserAdd extends Request {
  readonly url = '/user/add'
  readonly method = 'post'
  public param: {
    account: string
    password: string
    name: string
    contact: string
    superManager: number
  }
  constructor(param: UserAdd['param']) {
    super()
    this.param = param
  }
}

export class UserEdit extends Request {
  readonly url = '/user/edit'
  readonly method = 'put'
  public param: {
    id: number
    password: string
    name: string
    contact: string
    superManager: number
  }
  constructor(param: UserEdit['param']) {
    super()
    this.param = param
  }
}

export class UserRemove extends Request {
  readonly url = '/user/remove'
  readonly method = 'delete'
  public param: ID
  constructor(param: ID) {
    super()
    this.param = param
  }
}

export class UserChangePassword extends Request {
  readonly url = '/user/changePassword'
  readonly method = 'put'
  public param: {
    oldPwd: string
    newPwd: string
  }
  constructor(param: UserChangePassword['param']) {
    super()
    this.param = param
  }
}
