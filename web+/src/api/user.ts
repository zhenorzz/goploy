import Axios from './axios'
import { HttpResponse, Pagination, Total, ID } from './types'

export class Login {
  readonly url = '/user/login'
  readonly method = 'post'
  public param: {
    account: string
    password: string
  }
  public datagram!: {
    namespaceList: { id: number; name: string; role: string }[]
  }
  constructor(param: Login['param']) {
    this.param = param
  }
  public request(): Promise<HttpResponse<this['datagram']>> {
    return Axios.request({
      url: this.url,
      method: this.method,
      data: { ...this.param },
    })
  }
}

export class Info {
  readonly url = '/user/info'
  readonly method = 'get'
  public datagram!: {
    userInfo: {
      account: string
      id: number
      name: string
      superManager: number
    }
  }
  public request(): Promise<HttpResponse<this['datagram']>> {
    return Axios.request({
      url: this.url,
      method: this.method,
    })
  }
}

export class UserData {
  public datagram!: {
    detail: {
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
  }
}

export class UserList {
  readonly url = '/user/getList'
  readonly method = 'get'

  public pagination: Pagination

  public datagram!: {
    list: UserData['datagram']['detail'][]
  }
  constructor(pagination: Pagination) {
    this.pagination = pagination
  }
  public request(): Promise<HttpResponse<this['datagram']>> {
    return Axios.request({
      url: this.url,
      method: this.method,
      params: { ...this.pagination },
    })
  }
}

export class UserTotal {
  readonly url = '/user/getTotal'
  readonly method = 'get'

  public datagram!: Total

  public request(): Promise<HttpResponse<this['datagram']>> {
    return Axios.request({
      url: this.url,
      method: this.method,
    })
  }
}

export class UserOption {
  readonly url = '/user/getOption'
  readonly method = 'get'
  public datagram!: {
    list: {
      id: number
      account: string
      contact: string
      name: string
      password: string
      state: number
      superManager: number
      lastLoginTime: string
      insertTime: string
      updateTime: string
    }[]
  }
  public request(): Promise<HttpResponse<this['datagram']>> {
    return Axios.request({
      url: this.url,
      method: this.method,
    })
  }
}

export class UserAdd {
  readonly url = '/user/add'
  readonly method = 'post'
  public param: {
    account: string
    password: string
    name: string
    contact: string
    superManager: number
  }
  public datagram!: ID
  constructor(param: UserAdd['param']) {
    this.param = param
  }
  public request(): Promise<HttpResponse<this['datagram']>> {
    return Axios.request({
      url: this.url,
      method: this.method,
      data: { ...this.param },
    })
  }
}

export class UserEdit {
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
    this.param = param
  }
  public request(): Promise<HttpResponse<never>> {
    return Axios.request({
      url: this.url,
      method: this.method,
      data: { ...this.param },
    })
  }
}

export class UserRemove {
  readonly url = '/user/remove'
  readonly method = 'delete'
  public param: ID
  constructor(param: ID) {
    this.param = param
  }
  public request(): Promise<HttpResponse<never>> {
    return Axios.request({
      url: this.url,
      method: this.method,
      data: { ...this.param },
    })
  }
}

export class UserChangePassword {
  readonly url = '/user/changePassword'
  readonly method = 'put'
  public param: {
    oldPwd: string
    newPwd: string
  }
  constructor(param: UserChangePassword['param']) {
    this.param = param
  }
  public request(): Promise<HttpResponse<never>> {
    return Axios.request({
      url: this.url,
      method: this.method,
      data: { ...this.param },
    })
  }
}
