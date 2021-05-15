import Axios from './axios'
import { HttpResponse, Pagination, Total, ID } from './types'

export class Login {
  public param: {
    account: string
    password: string
  }
  public datagram!: {
    namespaceList: Array<{ id: number; name: string; role: string }>
  }
  constructor(param: Login['param']) {
    this.param = param
  }
  public request(): Promise<HttpResponse<this['datagram']>> {
    return Axios.request({
      url: '/user/login',
      method: 'post',
      data: { ...this.param },
    })
  }
}

export class Info {
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
      url: '/user/info',
      method: 'get',
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
  public pagination: Pagination
  public datagram!: {
    list: UserData['datagram']['detail'][]
  }
  constructor(pagination: Pagination) {
    this.pagination = pagination
  }
  public request(): Promise<HttpResponse<this['datagram']>> {
    return Axios.request({
      url: '/user/getList',
      method: 'get',
      params: { ...this.pagination },
    })
  }
}

export class UserTotal {
  public datagram!: Total

  public request(): Promise<HttpResponse<this['datagram']>> {
    return Axios.request({
      url: '/user/getTotal',
      method: 'get',
    })
  }
}

export class UserOption {
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
      url: '/user/getOption',
      method: 'get',
    })
  }
}

export class UserAdd {
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
      url: '/user/add',
      method: 'post',
      data: { ...this.param },
    })
  }
}

export class UserEdit {
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
      url: '/user/edit',
      method: 'post',
      data: { ...this.param },
    })
  }
}

export class UserRemove {
  public param: ID
  constructor(param: ID) {
    this.param = param
  }
  public request(): Promise<HttpResponse<never>> {
    return Axios.request({
      url: '/user/remove',
      method: 'delete',
      data: { ...this.param },
    })
  }
}

export class UserChangePassword {
  public param: {
    oldPwd: string
    newPwd: string
  }
  constructor(param: UserChangePassword['param']) {
    this.param = param
  }
  public request(): Promise<HttpResponse<never>> {
    return Axios.request({
      url: '/user/changePassword',
      method: 'put',
      data: { ...this.param },
    })
  }
}
