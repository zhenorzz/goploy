import { Request, ID } from './types'

export interface CommitData {
  branch: string
  commit: string
  author: string
  timestamp: number
  message: string
  tag: string
  diff: string
}

export class RepositoryBranchList extends Request {
  readonly url = '/repository/getBranchList'
  readonly method = 'get'
  readonly timeout = 0
  public param: ID

  public declare datagram: {
    list: string[]
  }
  constructor(param: RepositoryBranchList['param']) {
    super()
    this.param = { ...param }
  }
}

export class RepositoryCommitList extends Request {
  readonly url = '/repository/getCommitList'
  readonly method = 'get'
  readonly timeout = 0
  public param: {
    id: number
    branch: string
  }

  public declare datagram: {
    list: CommitData[]
  }
  constructor(param: RepositoryCommitList['param']) {
    super()
    this.param = { ...param }
  }
}

export class RepositoryTagList extends Request {
  readonly url = '/repository/getTagList'
  readonly method = 'get'
  readonly timeout = 0
  public param: ID

  public declare datagram: {
    list: CommitData[]
  }
  constructor(param: RepositoryTagList['param']) {
    super()
    this.param = { ...param }
  }
}

export interface RepositoryFile {
  [key: string]: any
  name: string
  mode: string
  isDir: boolean
}

export class RepositoryFileList extends Request {
  readonly url = '/repository/getFileList'
  readonly method = 'get'
  readonly timeout = 0
  public param: {
    id: number
    dir: string
  }

  public declare datagram: {
    list: RepositoryFile[]
  }
  constructor(param: RepositoryFileList['param']) {
    super()
    this.param = { ...param }
  }
}

export class RepositoryDeleteFile extends Request {
  readonly url = '/repository/deleteFile'
  readonly method = 'delete'
  readonly timeout = 0
  public param: {
    id: number
    file: string
  }
  constructor(param: RepositoryDeleteFile['param']) {
    super()
    this.param = param
  }
}
