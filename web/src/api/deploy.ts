import { Request, Pagination, Total, ID } from './types'
import { ProjectData } from './project'

export enum DeployState {
  Uninitialized = 0,
  Deploying = 1,
  Success = 2,
  Fail = 3,
}

export enum PublishTraceType {
  Queue,
  BeforePull,
  Pull,
  AfterPull,
  BeforeDeploy,
  Deploy,
  AfterDeploy,
  DeployFinish,
  PublishFinish,
}

export interface PublishTraceData {
  id: number
  token: string
  projectId: number
  projectName: string
  detail: string
  state: number
  publisherId: number
  publisherName: string
  type: number
  ext: string
  serverName: string
  insertTime: string
  updateTime: string
}

export interface PublishTraceExt {
  branch: string
  commit: string
  message: string
  author: string
  timestamp: number
  diff: string
  script: string
  command: string
  step: string
}

export class DeployList extends Request {
  readonly url = '/deploy/getList'
  readonly method = 'get'
  public declare datagram: {
    list: ProjectData[]
  }
}

export class DeployRebuild extends Request {
  readonly url = '/deploy/rebuild'
  readonly method = 'post'
  public param: {
    projectId: number
    token: string
  }
  public declare datagram: {
    type: string
    token: string
  }
  constructor(param: DeployRebuild['param']) {
    super()
    this.param = param
  }
}

export class DeployPreviewList extends Request {
  readonly url = '/deploy/getPreview'
  readonly method = 'get'
  public param: {
    projectId: number
    userId: number
    state: number
    commitDate: string
    branch: string
    commit: string
    filename: string
    deployDate: string
    token: string
  }

  public pagination: Pagination

  public declare datagram: {
    list: PublishTraceData[]
    pagination: Pagination & Total
  }

  constructor(param: DeployPreviewList['param'], pagination: Pagination) {
    super()
    this.pagination = pagination
    this.param = { ...param, ...pagination }
  }
}

export class DeployTrace extends Request {
  readonly url = '/deploy/getPublishTrace'
  readonly method = 'get'
  public param: {
    lastPublishToken: string
  }
  public declare datagram: {
    list: PublishTraceData[]
  }
  constructor(param: DeployTrace['param']) {
    super()
    this.param = param
  }
}

export class DeployTraceDetail extends Request {
  readonly url = '/deploy/getPublishTraceDetail'
  readonly method = 'get'
  readonly timeout = 0
  public param: ID
  public declare datagram: {
    detail: string
  }
  constructor(param: DeployTraceDetail['param']) {
    super()
    this.param = param
  }
}

export class DeployPublish extends Request {
  readonly url = '/deploy/publish'
  readonly method = 'post'
  public param: {
    projectId: number
    branch: string
    commit: string
  }
  constructor(param: DeployPublish['param']) {
    super()
    this.param = param
  }
}

export class DeployGreyPublish extends Request {
  readonly url = '/deploy/greyPublish'
  readonly method = 'post'
  public param: {
    projectId: number
    commit: string
    serverIds: number[]
  }
  constructor(param: DeployGreyPublish['param']) {
    super()
    this.param = param
  }
}

export class DeployResetState extends Request {
  readonly url = '/deploy/resetState'
  readonly method = 'put'
  public param: {
    projectId: number
  }
  constructor(param: DeployResetState['param']) {
    super()
    this.param = param
  }
}

export class DeployReview extends Request {
  readonly url = '/deploy/review'
  readonly method = 'put'
  readonly timeout = 0
  public param: {
    projectReviewId: number
    state: number
  }
  constructor(param: DeployReview['param']) {
    super()
    this.param = param
  }
}

export class FileCompare extends Request {
  readonly url = '/deploy/fileCompare'
  readonly method = 'post'
  readonly timeout = 0
  public param: {
    projectId: number
    filePath: string
  }
  constructor(param: FileCompare['param']) {
    super()
    this.param = param
  }
}

export class FileDiff extends Request {
  readonly url = '/deploy/fileDiff'
  readonly method = 'post'
  readonly timeout = 0
  public param: {
    projectId: number
    serverId: number
    filePath: string
  }
  constructor(param: FileDiff['param']) {
    super()
    this.param = param
  }
}

export class ManageProcess extends Request {
  readonly url = '/deploy/manageProcess'
  readonly method = 'post'
  readonly timeout = 0
  public param: {
    serverId: number
    projectProcessId: number
    command: string
  }
  public declare datagram: {
    execRes: boolean
    stdout: string
    stderr: string
    startTime: string
    endTime: string
  }
  constructor(param: ManageProcess['param']) {
    super()
    this.param = param
  }
}
