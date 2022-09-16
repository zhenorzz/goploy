import { Request, Pagination, Total } from './types'

export interface LoginLogData {
  id: number
  account: string
  remoteAddr: string
  userAgent: string
  referer: string
  reason: string
  loginTime: string
}

export interface OperationLogData {
  [key: string]: any
  id: number
  namespaceId: number
  userId: number
  router: string
  api: string
  requestTime: string
  requestData: string
  responseTime: string
  responseData: string
}

export interface SftpLogData {
  id: number
  namespaceId: number
  userId: number
  username: string
  serverId: number
  serverName: string
  remoteAddr: string
  userAgent: string
  type: string
  path: string
  reason: string
}

export interface TerminalLogData {
  id: number
  namespaceId: number
  userId: number
  username: string
  serverId: number
  serverName: string
  remoteAddr: string
  userAgent: string
  startTime: string
  endTime: string
}

export interface PublishLogData {
  token: string
  publisherId: number
  publisherName: string
  projectId: number
  projectName: string
  state: number
  insertTime: string
}

export class LoginLogList extends Request {
  readonly url = '/log/getLoginLogList'
  readonly method = 'get'

  public pagination: Pagination

  public param: {
    account: string
  }

  public declare datagram: {
    list: LoginLogData[]
  }
  constructor(param: LoginLogList['param'], pagination: Pagination) {
    super()
    this.pagination = pagination
    this.param = { ...param, ...pagination }
  }
}

export class LoginLogTotal extends Request {
  readonly url = '/log/getLoginLogTotal'
  readonly method = 'get'

  public param: {
    account: string
  }

  public declare datagram: Total

  constructor(param: LoginLogTotal['param']) {
    super()
    this.param = param
  }
}

export class OperationLogList extends Request {
  readonly url = '/log/getOperationLogList'
  readonly method = 'get'

  public pagination: Pagination

  public param: {
    username: string
    router: string
    api: string
  }

  public declare datagram: {
    list: OperationLogData[]
  }
  constructor(param: OperationLogList['param'], pagination: Pagination) {
    super()
    this.pagination = pagination
    this.param = { ...param, ...pagination }
  }
}

export class OperationLogTotal extends Request {
  readonly url = '/log/getOperationLogTotal'
  readonly method = 'get'

  public param: {
    username: string
    router: string
    api: string
  }

  public declare datagram: Total

  constructor(param: OperationLogTotal['param']) {
    super()
    this.param = param
  }
}

export class SftpLogList extends Request {
  readonly url = '/log/getSftpLogList'
  readonly method = 'get'

  public pagination: Pagination

  public param: {
    username: string
    serverName: string
  }

  public declare datagram: {
    list: SftpLogData[]
  }
  constructor(param: SftpLogList['param'], pagination: Pagination) {
    super()
    this.pagination = pagination
    this.param = { ...param, ...pagination }
  }
}

export class SftpLogTotal extends Request {
  readonly url = '/log/getSftpLogTotal'
  readonly method = 'get'

  public param: {
    username: string
    serverName: string
  }

  public declare datagram: Total

  constructor(param: SftpLogTotal['param']) {
    super()
    this.param = param
  }
}

export class TerminalLogList extends Request {
  readonly url = '/log/getTerminalLogList'
  readonly method = 'get'

  public pagination: Pagination

  public param: {
    username: string
    serverName: string
  }

  public declare datagram: {
    list: TerminalLogData[]
  }
  constructor(param: TerminalLogList['param'], pagination: Pagination) {
    super()
    this.pagination = pagination
    this.param = { ...param, ...pagination }
  }
}

export class TerminalLogTotal extends Request {
  readonly url = '/log/getTerminalLogTotal'
  readonly method = 'get'

  public param: {
    username: string
    serverName: string
  }

  public declare datagram: Total

  constructor(param: TerminalLogTotal['param']) {
    super()
    this.param = param
  }
}

export class PublishLogList extends Request {
  readonly url = '/log/getPublishLogList'
  readonly method = 'get'

  public pagination: Pagination

  public param: {
    username: string
    projectName: string
  }

  public declare datagram: {
    list: PublishLogData[]
  }
  constructor(param: PublishLogList['param'], pagination: Pagination) {
    super()
    this.pagination = pagination
    this.param = { ...param, ...pagination }
  }
}

export class PublishLogTotal extends Request {
  readonly url = '/log/getPublishLogTotal'
  readonly method = 'get'

  public param: {
    username: string
    projectName: string
  }

  public declare datagram: Total

  constructor(param: PublishLogTotal['param']) {
    super()
    this.param = param
  }
}
