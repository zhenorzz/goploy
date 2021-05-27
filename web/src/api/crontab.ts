import { Request, Pagination, ID, Total } from './types'

export class CrontabData {
  public datagram!: {
    detail: {
      id: number
      namespaceId: number
      command: string
      commandMD5: string
      date?: string
      dateLocale?: string
      script?: string
      description?: string
      creator: string
      creatorId: number
      editor: string
      editorId: number
      InsertTime: string
      UpdateTime: string
    }
  }
}

export class CrontabServerData {
  public datagram!: {
    detail: {
      id: number
      crontabId: number
      serverId: number
      serverName: string
      serverIP: string
      serverPort: number
      serverOwner: string
      serverDescription: string
      insertTime: string
      updateTime: string
    }
  }
}

export class CrontabList extends Request {
  readonly url = '/crontab/getList'
  readonly method = 'get'
  public param: {
    command: string
  }
  public pagination: Pagination

  public datagram!: {
    list: CrontabData['datagram']['detail'][]
  }
  constructor(param: CrontabList['param'], pagination: Pagination) {
    super()
    this.pagination = pagination
    this.param = { ...param, ...pagination }
  }
}

export class CrontabTotal extends Request {
  readonly url = '/crontab/getTotal'
  readonly method = 'get'

  public param: {
    command: string
  }

  public datagram!: Total
  constructor(param: CrontabTotal['param']) {
    super()
    this.param = { ...param }
  }
}

export class CrontabsInRemoteServer extends Request {
  readonly url = '/crontab/getCrontabsInRemoteServer'
  readonly method = 'get'
  public param: {
    serverId: number
  }
  public datagram!: {
    list: string[]
  }
  constructor(param: CrontabsInRemoteServer['param']) {
    super()
    this.param = param
  }
}

export class CrontabServerList extends Request {
  readonly url = '/crontab/getBindServerList'
  readonly method = 'get'
  public param: ID
  public datagram!: {
    list: CrontabServerData['datagram']['detail'][]
  }
  constructor(param: CrontabServerList['param']) {
    super()
    this.param = param
  }
}

export class CrontabAdd extends Request {
  readonly url = '/crontab/add'
  readonly method = 'post'
  public param: {
    command: string
    serverIds: string[]
  }
  public datagram!: ID
  constructor(param: CrontabAdd['param']) {
    super()
    this.param = param
  }
}

export class CrontabEdit extends Request {
  readonly url = '/crontab/edit'
  readonly method = 'put'
  public param: {
    id: number
    command: string
  }
  constructor(param: CrontabEdit['param']) {
    super()
    this.param = param
  }
}

export class CrontabRemove extends Request {
  readonly url = '/crontab/remove'
  readonly method = 'delete'
  public param: {
    id: number
    radio: number
  }
  public datagram!: ID
  constructor(param: CrontabRemove['param']) {
    super()
    this.param = param
  }
}

export class CrontabImport extends Request {
  readonly url = '/crontab/import'
  readonly method = 'post'
  public param: {
    serverId: number
    commands: string[]
  }
  public datagram!: {
    list: string[]
  }
  constructor(param: CrontabImport['param']) {
    super()
    this.param = param
  }
}

export class CrontabServerAdd extends Request {
  readonly url = '/crontab/addServer'
  readonly method = 'post'
  public param: {
    crontabId: number
    serverIds: number[]
  }
  constructor(param: CrontabServerAdd['param']) {
    super()
    this.param = param
  }
}

export class CrontabServerRemove extends Request {
  readonly url = '/crontab/removeCrontabServer'
  readonly method = 'delete'
  public param: {
    crontabServerId: number
    crontabId: number
    serverId: number
  }
  constructor(param: CrontabServerRemove['param']) {
    super()
    this.param = param
  }
}
