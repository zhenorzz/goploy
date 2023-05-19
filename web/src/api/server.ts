import { HttpResponse, Request, ID } from './types'

export interface ServerData {
  id: number
  label: string
  name: string
  os: string
  ip: string
  port: number
  owner: string
  path: string
  password: string
  jumpIP: string
  jumpPort: number
  jumpOwner: string
  jumpPath: string
  jumpPassword: string
  namespaceId: number
  description: string
  state: number
  insertTime: string
  updateTime: string
}

export interface ServerMonitorData {
  id: number
  serverId: number
  item: string
  formula: string
  operator: string
  value: string
  groupCycle: number
  lastCycle: number
  silentCycle: number
  startTime: string
  endTime: string
  notifyType: number
  notifyTarget: string
  insertTime: string
  updateTime: string
}

export interface ServerSFTPFile {
  [key: string]: any
  uuid: number
  isDir: boolean
  modTime: string
  mode: string
  name: string
  size: number
  icon: string
  uploading: boolean
}

export class ServerList extends Request {
  readonly url = '/server/getList'
  readonly method = 'get'

  public declare datagram: {
    list: ServerData[]
  }

  public request(): Promise<HttpResponse<this['datagram']>> {
    return super.request().then((response) => {
      response.data.list = response.data.list.map((element) => {
        element.label =
          element.name +
          (element.description.length > 0
            ? '(' + element.description + ')'
            : '')
        return element
      })
      return response
    })
  }
}

export class ServerPublicKey extends Request {
  readonly url = '/server/getPublicKey'
  readonly method = 'get'

  public param: {
    path: string
  }

  public declare datagram: {
    key: string
  }

  constructor(param: ServerPublicKey['param']) {
    super()
    this.param = param
  }
}

export class ServerOption extends Request {
  readonly url = '/server/getOption'
  readonly method = 'get'

  public declare datagram: {
    list: ServerData[]
  }
  public request(): Promise<HttpResponse<this['datagram']>> {
    return super.request().then((response) => {
      response.data.list = response.data.list.map((element) => {
        element.label =
          element.name +
          (element.description.length > 0
            ? '(' + element.description + ')'
            : '')
        return element
      })
      return response
    })
  }
}

export class ServerAdd extends Request {
  readonly url = '/server/add'
  readonly method = 'post'
  public param: {
    namespaceId: number
    name: string
    os: string
    ip: string
    port: number
    owner: string
    path: string
    password: string
    jumpIP: string
    jumpPort: number
    jumpOwner: string
    jumpPath: string
    jumpPassword: string
    description: string
  }
  constructor(param: ServerAdd['param']) {
    super()
    this.param = param
  }
}

export class ServerEdit extends Request {
  readonly url = '/server/edit'
  readonly method = 'put'
  public param: {
    id: number
    namespaceId: number
    name: string
    os: string
    ip: string
    port: number
    owner: string
    path: string
    password: string
    jumpIP: string
    jumpPort: number
    jumpOwner: string
    jumpPath: string
    jumpPassword: string
    description: string
  }
  constructor(param: ServerEdit['param']) {
    super()
    this.param = param
  }
}

export class ServerCheck extends Request {
  readonly url = '/server/check'
  readonly method = 'post'
  public param: {
    ip: string
    port: number
    owner: string
    path: string
    password: string
    jumpIP: string
    jumpPort: number
    jumpOwner: string
    jumpPath: string
    jumpPassword: string
  }
  constructor(param: ServerCheck['param']) {
    super()
    this.param = param
  }
}

export class ServerToggle extends Request {
  readonly url = '/server/toggle'
  readonly method = 'put'
  public param: {
    id: number
    state: number
  }
  constructor(param: ServerToggle['param']) {
    super()
    this.param = param
  }
}

export class ServerInstallAgent extends Request {
  readonly url = '/server/installAgent'
  readonly method = 'post'
  public param: {
    ids: number[]
    installPath: string
    tool: string
    reportURL: string
    webPort: string
  }
  constructor(param: ServerInstallAgent['param']) {
    super()
    this.param = param
  }
}

export class ServerReport extends Request {
  readonly url = '/server/report'
  readonly method = 'get'
  public param: {
    serverId: number
    type: number
    datetimeRange: string
  }
  constructor(param: ServerReport['param']) {
    super()
    this.param = param
  }
}

export class ServerMonitorList extends Request {
  readonly url = '/server/getAllMonitor'
  readonly method = 'get'

  public param: {
    serverId: number
  }

  public declare datagram: {
    list: ServerMonitorData[]
  }

  constructor(param: ServerMonitorList['param']) {
    super()
    this.param = param
  }
}

export class ServerMonitorAdd extends Request {
  readonly url = '/server/addMonitor'
  readonly method = 'post'
  public param: {
    serverId: number
    item: string
    formula: string
    operator: string
    value: string
    groupCycle: number
    lastCycle: number
    silentCycle: number
    startTime: string
    endTime: string
    notifyType: number
    notifyTarget: string
  }
  constructor(param: ServerMonitorAdd['param']) {
    super()
    this.param = param
  }
}

export class ServerMonitorEdit extends Request {
  readonly url = '/server/editMonitor'
  readonly method = 'put'
  public param: {
    id: number
    serverId: number
    item: string
    formula: string
    operator: string
    value: string
    groupCycle: number
    lastCycle: number
    silentCycle: number
    startTime: string
    endTime: string
    notifyType: number
    notifyTarget: string
  }
  constructor(param: ServerMonitorEdit['param']) {
    super()
    this.param = param
  }
}

export class ServerMonitorDelete extends Request {
  readonly url = '/server/deleteMonitor'
  readonly method = 'delete'
  public param: ID
  constructor(param: ServerMonitorDelete['param']) {
    super()
    this.param = param
  }
}

export interface ServerProcessData {
  [key: string]: any
  id: number
  name: string
  items: { name: string; command: string }[]
  InsertTime: string
  UpdateTime: string
}

export class ServerProcessList extends Request {
  readonly url = '/server/getProcessList'
  readonly method = 'get'

  public declare datagram: {
    list: ServerProcessData[]
  }

  public request(): Promise<HttpResponse<this['datagram']>> {
    return super.request().then((response) => {
      response.data.list = response.data.list.map((element) => {
        element.items = JSON.parse(element.items as any)
        return element
      })
      return response
    })
  }
}

export class ServerProcessAdd extends Request {
  readonly url = '/server/addProcess'
  readonly method = 'post'
  public param: {
    name: string
    items: string
  }
  public declare datagram: ID
  constructor(param: ServerProcessAdd['param']) {
    super()
    this.param = param
  }
}

export class ServerProcessEdit extends Request {
  readonly url = '/server/editProcess'
  readonly method = 'put'
  public param: {
    id: number
    name: string
    items: string
  }
  constructor(param: ServerProcessEdit['param']) {
    super()
    this.param = param
  }
}

export class ServerProcessDelete extends Request {
  readonly url = '/server/deleteProcess'
  readonly method = 'delete'
  public param: {
    id: number
  }
  constructor(param: ServerProcessDelete['param']) {
    super()
    this.param = param
  }
}

export class ServerExecProcess extends Request {
  readonly url = '/server/execProcess'
  readonly method = 'post'
  readonly timeout = 0
  public param: {
    id: number
    serverId: number
    name: string
  }
  public declare datagram: {
    serverId: number
    execRes: boolean
    stdout: string
    stderr: string
  }
  constructor(param: ServerExecProcess['param']) {
    super()
    this.param = param
  }
}

export class ServerCopyFile extends Request {
  readonly url = '/server/copyFile'
  readonly method = 'put'
  readonly timeout = 0
  public param: {
    serverId: number
    dir: string
    isDir: boolean
    srcName: string
    dstName: string
  }
  constructor(param: ServerCopyFile['param']) {
    super()
    this.param = param
  }
}

export class ServerRenameFile extends Request {
  readonly url = '/server/renameFile'
  readonly method = 'put'
  readonly timeout = 0
  public param: {
    serverId: number
    dir: string
    currentName: string
    newName: string
  }
  constructor(param: ServerRenameFile['param']) {
    super()
    this.param = param
  }
}

export class ServerEditFile extends Request {
  readonly url = '/server/editFile'
  readonly method = 'put'
  readonly timeout = 0
  public param: {
    serverId: number
    file: string
    content: string
  }
  constructor(param: ServerEditFile['param']) {
    super()
    this.param = param
  }
}

export class ServerDeleteFile extends Request {
  readonly url = '/server/deleteFile'
  readonly method = 'delete'
  readonly timeout = 0
  public param: {
    serverId: number
    file: string
  }
  constructor(param: ServerDeleteFile['param']) {
    super()
    this.param = param
  }
}

export class ServerTransferFile extends Request {
  readonly url = '/server/transferFile'
  readonly method = 'post'
  readonly timeout = 0
  public param: {
    sourceServerId: number
    sourceFile: string
    sourceIsDir: boolean
    destServerIds: number[]
    destDir: string
  }
  constructor(param: ServerTransferFile['param']) {
    super()
    this.param = param
  }
}
export class ServerExecScript extends Request {
  readonly url = '/server/execScript'
  readonly method = 'post'
  readonly timeout = 0
  public param: {
    serverIds: number[]
    script: string
  }

  public declare datagram: {
    serverId: number
    execRes: boolean
    stderr: string
    stdout: string
  }[]

  constructor(param: ServerExecScript['param']) {
    super()
    this.param = param
  }
}

export class ServerRemoteCrontabList extends Request {
  readonly url = '/server/getRemoteCrontabList'
  readonly method = 'get'
  public param: {
    serverId: number
  }
  constructor(param: ServerRemoteCrontabList['param']) {
    super()
    this.param = param
  }
}
