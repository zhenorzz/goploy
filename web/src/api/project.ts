import { Request, Pagination, ID, Total } from './types'

export class ProjectData {
  public datagram!: {
    id: number
    namespaceId: number
    userId: number
    name: string
    repoType: string
    url: string
    path: string
    environment: number
    branch: string
    symlinkPath: string
    symlinkBackupNumber: number
    review: number
    reviewURL: string
    afterPullScriptMode: string
    afterPullScript: string
    afterDeployScriptMode: string
    afterDeployScript: string
    rsyncOption: string
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
}

export class ProjectList extends Request {
  readonly url = '/project/getList'
  readonly method = 'get'
  public param: {
    projectName: string
  }
  public pagination: Pagination

  public datagram!: {
    list: ProjectData['datagram'][]
  }
  constructor(param: ProjectList['param'], pagination: Pagination) {
    super()
    this.pagination = pagination
    this.param = { ...param, ...pagination }
  }
}

export class ProjectTotal extends Request {
  readonly url = '/project/getTotal'
  readonly method = 'get'

  public param: {
    projectName: string
  }

  public datagram!: Total
  constructor(param: ProjectTotal['param']) {
    super()
    this.param = { ...param }
  }
}

export class ProjectPingRepos extends Request {
  readonly url = '/project/pingRepos'
  readonly method = 'get'

  public param: {
    repoType: string
    url: string
  }

  public datagram!: {
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

  public datagram!: {
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
    path: string
    environment: number
    branch: string
    symlinkPath: string
    review: number
    reviewURL: string
    afterPullScriptMode: string
    afterPullScript: string
    afterDeployScriptMode: string
    afterDeployScript: string
    rsyncOption: string
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
    path: string
    symlinkPath: string
    review: number
    reviewURL: string
    environment: number
    branch: string
    afterPullScriptMode: string
    afterPullScript: string
    afterDeployScriptMode: string
    afterDeployScript: string
    rsyncOption: string
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

export class ProjectServerData {
  public datagram!: {
    id: number
    projectId: number
    serverId: number
    serverName: string
    serverIP: string
    serverPort: number
    serverOwner: string
    serverPassword: string
    serverPath: string
    serverDescription: string
    insertTime: string
    updateTime: string
  }
}
export class ProjectServerList extends Request {
  readonly url = '/project/getBindServerList'
  readonly method = 'get'
  public param: ID
  public datagram!: {
    list: ProjectServerData['datagram'][]
  }
  constructor(param: ProjectServerList['param']) {
    super()
    this.param = param
  }
}

export class ProjectServerAdd extends Request {
  readonly url = '/project/addServer'
  readonly method = 'post'
  public param: {
    projectId: number
    serverIds: number[]
  }
  constructor(param: ProjectServerAdd['param']) {
    super()
    this.param = param
  }
}

export class ProjectServerRemove extends Request {
  readonly url = '/project/removeServer'
  readonly method = 'delete'
  public param: {
    projectServerId: number
  }
  constructor(param: ProjectServerRemove['param']) {
    super()
    this.param = param
  }
}

export class ProjectUserData {
  public datagram!: {
    id: number
    namespaceId: number
    projectId: number
    projectName: string
    userId: number
    userName: string
    role: string
    insertTime: string
    updateTime: string
  }
}
export class ProjectUserList extends Request {
  readonly url = '/project/getBindUserList'
  readonly method = 'get'
  public param: ID
  public datagram!: {
    list: ProjectUserData['datagram'][]
  }
  constructor(param: ProjectUserList['param']) {
    super()
    this.param = param
  }
}

export class ProjectUserAdd extends Request {
  readonly url = '/project/addUser'
  readonly method = 'post'
  public param: {
    projectId: number
    userIds: number[]
  }
  constructor(param: ProjectUserAdd['param']) {
    super()
    this.param = param
  }
}

export class ProjectUserRemove extends Request {
  readonly url = '/project/removeUser'
  readonly method = 'delete'
  public param: {
    projectUserId: number
  }
  constructor(param: ProjectUserRemove['param']) {
    super()
    this.param = param
  }
}

export class ProjectFileData {
  public datagram!: {
    detail: {
      id: number
      projectId: number
      filename: string
      insertTime: string
      updateTime: string
    }
  }
}
export class ProjectFileList extends Request {
  readonly url = '/project/getProjectFileList'
  readonly method = 'get'
  public param: ID
  public datagram!: {
    list: ProjectFileData['datagram'][]
  }
  constructor(param: ProjectUserList['param']) {
    super()
    this.param = param
  }
}

export class ProjectFileContent extends Request {
  readonly url = '/project/getProjectFileContent'
  readonly method = 'get'
  public param: ID
  public datagram!: {
    content: string
  }
  constructor(param: ProjectUserList['param']) {
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
  public datagram!: ID
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

export class ProjectTaskData {
  public datagram!: {
    detail: {
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
  }
}

export class ProjectTaskList extends Request {
  readonly url = '/project/getTaskList'
  readonly method = 'get'
  public param: ID
  public pagination: Pagination

  public datagram!: {
    list: ProjectTaskData['datagram'][]
    pagination: Pagination
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
  public datagram!: ID
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

  public datagram!: {
    list: ProjectTaskData['datagram'][]
    pagination: Pagination
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
