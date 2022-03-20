import { Request, Pagination, ID } from './types'

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

export class CronList extends Request {
  readonly url = '/cron/getList'
  readonly method = 'post'
  public param: {
    serverId: number
  }
  public pagination: Pagination

  public declare datagram: {
    list: CronData[]
  }
  constructor(param: CronList['param'], pagination: Pagination) {
    super()
    this.pagination = pagination
    this.param = { ...param, ...pagination }
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
