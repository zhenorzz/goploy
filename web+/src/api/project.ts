import Axios from './axios'
import { Request, Pagination, ID, Total } from './types'

export class ProjectServerData {
  public datagram!: {
    detail: {
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
}
export class ProjectServerList extends Request {
  readonly url = '/project/getBindServerList'
  readonly method = 'get'
  public param: ID
  public datagram!: {
    list: ProjectServerData['datagram']['detail'][]
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
    detail: {
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
}
export class ProjectUserList extends Request {
  readonly url = '/project/getBindUserList'
  readonly method = 'get'
  public param: ID
  public datagram!: {
    list: ProjectUserData['datagram']['detail'][]
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
    list: ProjectFileData['datagram']['detail'][]
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

/**
 * @return {Promise}
 */
export function getList({ page, rows }, projectName) {
  return Axios.request({
    url: '/project/getList',
    method: 'get',
    params: { page, rows, projectName },
  })
}

/**
 * @return {Promise}
 */
export function getTotal(projectName) {
  return Axios.request({
    url: '/project/getTotal',
    method: 'get',
    params: { projectName },
  })
}

/**
 * @return {Promise}
 */
export function getRemoteBranchList(url) {
  return Axios.request({
    url: '/project/getRemoteBranchList',
    method: 'get',
    timeout: 0,
    params: { url },
  })
}

/**
 * @param  {string} project
 * @param  {string} owner
 * @param  {string} repository
 * @param  {string} serverIds
 * @param  {string} userIds
 * @return {Promise}
 */
export function add(data) {
  return Axios.request({
    url: '/project/add',
    method: 'post',
    data,
  })
}

/**
 * @param  {string} project
 * @param  {string} owner
 * @param  {string} repository
 * @return {Promise}
 */
export function edit(data) {
  return Axios.request({
    url: '/project/edit',
    method: 'put',
    data,
  })
}

export function setAutoDeploy(data) {
  return Axios.request({
    url: '/project/setAutoDeploy',
    method: 'put',
    data,
  })
}

export function remove(id) {
  return Axios.request({
    url: '/project/remove',
    method: 'delete',
    data: { id },
  })
}

export function addTask(data) {
  return Axios.request({
    url: '/project/addTask',
    method: 'post',
    data,
  })
}

export function removeTask(id) {
  return Axios.request({
    url: '/project/removeTask',
    method: 'delete',
    data: { id },
  })
}

/**
 * @param  {id} id
 * @return {Promise}
 */
export function getTaskList({ page, rows }, id) {
  return Axios.request({
    url: '/project/getTaskList',
    method: 'get',
    params: { page, rows, id },
  })
}

/**
 * @param  {id} id
 * @return {Promise}
 */
export function getReviewList({ page, rows }, id) {
  return Axios.request({
    url: '/project/getReviewList',
    method: 'get',
    params: { page, rows, id },
  })
}
