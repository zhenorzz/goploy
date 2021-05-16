import Axios from './axios'
import { HttpResponse, Pagination, Total, ID } from './types'

export class NamespaceData {
  public datagram!: {
    detail: {
      id: number
      name: string
      userId: number
      role: string
      insertTime: string
      updateTime: string
    }
  }
}

export class NamespaceUserData {
  public datagram!: {
    detail: {
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
}

export class NamespaceList {
  readonly url = '/namespace/getList'
  readonly method = 'get'

  public pagination: Pagination

  public datagram!: {
    list: NamespaceData['datagram']['detail'][]
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

export class NamespaceTotal {
  readonly url = '/namespace/getTotal'
  readonly method = 'get'

  public datagram!: Total

  public request(): Promise<HttpResponse<this['datagram']>> {
    return Axios.request({
      url: this.url,
      method: this.method,
    })
  }
}

export class NamespaceUserOption {
  readonly url = '/namespace/getUserOption'
  readonly method = 'get'

  public datagram!: {
    list: NamespaceUserData['datagram']['detail'][]
  }

  public request(): Promise<HttpResponse<this['datagram']>> {
    return Axios.request({
      url: this.url,
      method: this.method,
    })
  }
}

export class NamespaceUserList {
  readonly url = '/namespace/getBindUserList'
  readonly method = 'get'
  public param: ID
  public datagram!: {
    list: NamespaceUserData['datagram']['detail'][]
  }
  constructor(param: NamespaceUserList['param']) {
    this.param = param
  }
  public request(): Promise<HttpResponse<this['datagram']>> {
    return Axios.request({
      url: this.url,
      method: this.method,
      params: { ...this.param },
    })
  }
}

export class NamespaceAdd {
  readonly url = '/namespace/add'
  readonly method = 'post'
  public param: {
    name: string
  }
  public datagram!: ID
  constructor(param: NamespaceAdd['param']) {
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

export class NamespaceEdit {
  readonly url = '/namespace/edit'
  readonly method = 'put'
  public param: {
    id: number
    name: string
  }
  constructor(param: NamespaceEdit['param']) {
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

export class NamespaceUserAdd {
  readonly url = '/namespace/addUser'
  readonly method = 'post'
  public param: {
    namespaceId: number
    userIds: number[]
    role: string
  }
  constructor(param: NamespaceUserAdd['param']) {
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

export class NamespaceUserRemove {
  readonly url = '/namespace/removeUser'
  readonly method = 'delete'
  public param: {
    namespaceUserId: number
  }
  constructor(param: NamespaceUserRemove['param']) {
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
