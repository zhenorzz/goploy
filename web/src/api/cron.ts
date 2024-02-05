import { Request, ID } from './types'

export interface CronData {
  id: number
  serverId: number
  expression: string
  command: string
  singleMode: number
  logLevel: number
  description: string
  creator: string
  creatorId: number
  editor: string
  editorId: number
  state: number
  InsertTime: string
  UpdateTime: string
}

export interface CronLogData {
  id: number
  serverId: number
  cronId: number
  execCode: number
  message: string
  ReportTime: string
  InsertTime: string
}

export class CronLogs extends Request {
  readonly url = '/cron/getLogs'
  readonly method = 'get'
  public param: {
    cronId: number
    serverId: number
    page: number
    rows: number
  }

  public declare datagram: {
    list: CronLogData[]
  }
  constructor(param: CronLogs['param']) {
    super()
    this.param = param
  }
}

export class CronList extends Request {
  readonly url = '/cron/getList'
  readonly method = 'post'
  public param: {
    serverId: number
  }

  public declare datagram: {
    list: CronData[]
  }
  constructor(param: CronList['param']) {
    super()
    this.param = param
  }
}

export class CronAdd extends Request {
  readonly url = '/cron/add'
  readonly method = 'post'
  public param: {
    expression: string
    command: string
    singleMode: number
    logLevel: number
    description: string
  }
  public declare datagram: ID
  constructor(param: CronAdd['param']) {
    super()
    this.param = param
  }
}

export class CronEdit extends Request {
  readonly url = '/cron/edit'
  readonly method = 'put'
  public param: {
    id: number
    expression: string
    command: string
    singleMode: number
    logLevel: number
    description: string
  }
  constructor(param: CronEdit['param']) {
    super()
    this.param = param
  }
}

export class CronRemove extends Request {
  readonly url = '/cron/remove'
  readonly method = 'delete'
  public param: {
    id: number
  }
  constructor(param: CronRemove['param']) {
    super()
    this.param = param
  }
}
