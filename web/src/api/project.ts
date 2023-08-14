import { ServerData } from './server'
import { Request, Pagination, ID, Total } from './types'
import { UserData } from './user'

export interface ProjectScript {
  afterPull: { mode: string; content: string }
  afterDeploy: { mode: string; content: string }
  deployFinish: { mode: string; content: string }
}

export interface ProjectData {
  [key: string]: any
  id: number
  namespaceId: number
  userId: number
  name: string
  repoType: string
  url: string
  path: string
  environment: number
  branch: string
  label: string
  symlinkPath: string
  symlinkBackupNumber: number
  review: number
  reviewURL: string
  script: ProjectScript
  afterPullScriptMode: string
  afterPullScript: string
  afterDeployScriptMode: string
  afterDeployScript: string
  transferType: string
  transferOption: string
  deployServerMode: string
  autoDeploy: number
  publisherId: number
  publisherName: string
  publishExt: string
  deployState: number
  lastPublishToken: string
  notifyType: number
  notifyTarget: string
  serverIds: number[]
  userIds: number[]
  state: number
  insertTime: string
  updateTime: string
}

export interface ProjectServerData {
  id: number
  projectId: number
  serverId: number
  project: ProjectData
  server: ServerData
  insertTime: string
  updateTime: string
}

export interface ProjectUserData {
  id: number
  projectId: number
  userId: number
  project: ProjectData
  user: UserData
  role: string
  insertTime: string
  updateTime: string
}

export interface ProjectFileData {
  id: number
  projectId: number
  filename: string
  insertTime: string
  updateTime: string
}

export interface ProjectTaskData {
  id: number
  projectId: number
  branch: string
  commit: string
  date: string
  state: number
  isRun: number
  creator: string
  creatorId: number
  editor: string
  editorId: number
  insertTime: string
  updateTime: string
}

export interface ProjectProcessData {
  id: number
  projectId: number
  name: string
  status: string
  start: string
  stop: string
  restart: string
  insertTime: string
  updateTime: string
}
export class ProjectList extends Request {
  readonly url = '/project/getList'
  readonly method = 'get'

  public declare datagram: {
    list: ProjectData[]
  }
}

export class LabelList extends Request {
  readonly url = '/project/getLabelList'
  readonly method = 'get'

  public declare datagram: {
    list: string[]
  }
}

export class ProjectPingRepos extends Request {
  readonly url = '/project/pingRepos'
  readonly method = 'get'

  public param: {
    repoType: string
    url: string
  }

  public declare datagram: {
    branch: string[]
  }
  constructor(param: ProjectPingRepos['param']) {
    super()
    this.param = { ...param }
  }
}

export class ProjectRemoteBranchList extends Request {
  readonly url = '/project/getRemoteBranchList'
  readonly method = 'get'

  public param: {
    repoType: string
    url: string
  }

  public declare datagram: {
    branch: string[]
  }
  constructor(param: ProjectRemoteBranchList['param']) {
    super()
    this.param = { ...param }
  }
}

export class ProjectAdd extends Request {
  readonly url = '/project/add'
  readonly method = 'post'
  public param: {
    name: string
    repoType: string
    url: string
    label: string
    path: string
    environment: number
    branch: string
    symlinkPath: string
    review: number
    reviewURL: string
    script: ProjectScript
    transferType: string
    transferOption: string
    deployServerMode: string
    serverIds: number[]
    userIds: number[]
    notifyType: number
    notifyTarget: string
  }
  constructor(param: ProjectAdd['param']) {
    super()
    this.param = param
  }
}

export class ProjectEdit extends Request {
  readonly url = '/project/edit'
  readonly method = 'put'
  public param: {
    id: number
    name: string
    repoType: string
    url: string
    label: string
    path: string
    symlinkPath: string
    review: number
    reviewURL: string
    environment: number
    branch: string
    script: ProjectScript
    transferType: string
    transferOption: string
    deployServerMode: string
    notifyType: number
    notifyTarget: string
  }
  constructor(param: ProjectEdit['param']) {
    super()
    this.param = param
  }
}

export class ProjectRemove extends Request {
  readonly url = '/project/remove'
  readonly method = 'delete'
  public param: ID
  constructor(param: ID) {
    super()
    this.param = param
  }
}

export class ProjectServerList extends Request {
  readonly url = '/project/getBindServerList'
  readonly method = 'get'
  public param: ID
  public declare datagram: {
    list: ProjectServerData[]
  }
  constructor(param: ProjectServerList['param']) {
    super()
    this.param = param
  }
}

export class ProjectUserList extends Request {
  readonly url = '/project/getBindUserList'
  readonly method = 'get'
  public param: ID
  public declare datagram: {
    list: ProjectUserData[]
  }
  constructor(param: ProjectUserList['param']) {
    super()
    this.param = param
  }
}

export class ProjectFileList extends Request {
  readonly url = '/project/getProjectFileList'
  readonly method = 'get'
  public param: ID
  public declare datagram: {
    list: ProjectFileData[]
  }
  constructor(param: ProjectFileList['param']) {
    super()
    this.param = param
  }
}

export class ProjectFileContent extends Request {
  readonly url = '/project/getProjectFileContent'
  readonly method = 'get'
  public param: ID
  public declare datagram: {
    content: string
  }
  constructor(param: ProjectFileContent['param']) {
    super()
    this.param = param
  }
}

export class ProjectFileAdd extends Request {
  readonly url = '/project/addFile'
  readonly method = 'post'
  public param: {
    projectId: number
    content: string
    filename: string
  }
  public declare datagram: ID
  constructor(param: ProjectFileAdd['param']) {
    super()
    this.param = param
  }
}

export class ProjectFileEdit extends Request {
  readonly url = '/project/editFile'
  readonly method = 'put'
  public param: {
    id: number
    content: string
  }
  constructor(param: ProjectFileEdit['param']) {
    super()
    this.param = param
  }
}

export class ProjectFileRemove extends Request {
  readonly url = '/project/removeFile'
  readonly method = 'delete'
  public param: {
    projectFileId: number
  }
  constructor(param: ProjectFileRemove['param']) {
    super()
    this.param = param
  }
}

export class ProjectTaskList extends Request {
  readonly url = '/project/getTaskList'
  readonly method = 'get'
  public param: ID
  public pagination: Pagination

  public declare datagram: {
    list: ProjectTaskData[]
    pagination: Pagination & Total
  }
  constructor(param: ProjectTaskList['param'], pagination: Pagination) {
    super()
    this.pagination = pagination
    this.param = { ...param, ...pagination }
  }
}

export class ProjectTaskAdd extends Request {
  readonly url = '/project/addTask'
  readonly method = 'post'
  public param: {
    projectId: number
    branch: string
    commit: string
    date: string
  }
  public declare datagram: ID
  constructor(param: ProjectTaskAdd['param']) {
    super()
    this.param = param
  }
}

export class ProjectTaskRemove extends Request {
  readonly url = '/project/removeTask'
  readonly method = 'delete'
  public param: ID
  constructor(param: ProjectTaskRemove['param']) {
    super()
    this.param = param
  }
}

export class ProjectAutoDeploy extends Request {
  readonly url = '/project/setAutoDeploy'
  readonly method = 'put'
  public param: {
    id: number
    autoDeploy: number
  }
  constructor(param: ProjectAutoDeploy['param']) {
    super()
    this.param = param
  }
}

export class ProjectReviewList extends Request {
  readonly url = '/project/getReviewList'
  readonly method = 'get'
  public param: ID
  public pagination: Pagination

  public declare datagram: {
    list: ProjectTaskData[]
    pagination: Pagination & Total
  }
  constructor(param: ProjectReviewList['param'], pagination: Pagination) {
    super()
    this.pagination = pagination
    this.param = { ...param, ...pagination }
  }
}

export class ReposFileList extends Request {
  readonly url = '/project/getReposFileList'
  readonly method = 'get'
  public param: {
    id: number
    path: string
  }

  constructor(param: ReposFileList['param']) {
    super()
    this.param = param
  }
}

export class ProjectProcessList extends Request {
  readonly url = '/project/getProcessList'
  readonly method = 'get'
  public param: {
    projectId: number
  }
  public pagination: Pagination

  public declare datagram: {
    list: ProjectProcessData[]
  }
  constructor(param: ProjectProcessList['param'], pagination: Pagination) {
    super()
    this.pagination = pagination
    this.param = { ...param, ...pagination }
  }
}

export class ProjectProcessAdd extends Request {
  readonly url = '/project/addProcess'
  readonly method = 'post'
  public param: {
    projectId: number
    name: string
    start: string
    stop: string
    status: string
    restart: string
  }
  public declare datagram: ID
  constructor(param: ProjectProcessAdd['param']) {
    super()
    this.param = param
  }
}

export class ProjectProcessEdit extends Request {
  readonly url = '/project/editProcess'
  readonly method = 'put'
  public param: {
    id: number
    name: string
    start: string
    stop: string
    status: string
    restart: string
  }
  constructor(param: ProjectProcessEdit['param']) {
    super()
    this.param = param
  }
}

export class ProjectProcessDelete extends Request {
  readonly url = '/project/deleteProcess'
  readonly method = 'delete'
  public param: ID
  constructor(param: ProjectTaskRemove['param']) {
    super()
    this.param = param
  }
}
