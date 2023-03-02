import { Request, ID } from './types'

export enum TemplateType {
  Script = 1,
}

export interface TemplateData {
  id: number
  name: string
  content: string
  description: string
  insertTime: string
  updateTime: string
}

export class TemplateOption extends Request {
  readonly url = '/template/getOption'
  readonly method = 'get'
  public param: {
    type: TemplateType
  }
  public declare datagram: {
    list: TemplateData[]
  }
  constructor(param: TemplateOption['param']) {
    super()
    this.param = param
  }
}

export class TemplateAdd extends Request {
  readonly url = '/template/add'
  readonly method = 'post'
  public param: {
    type: TemplateType
    name: string
    content: string
    description: string
  }
  public declare datagram: ID
  constructor(param: TemplateAdd['param']) {
    super()
    this.param = param
  }
}

export class TemplateRemove extends Request {
  readonly url = '/template/remove'
  readonly method = 'delete'
  public param: {
    id: number
  }
  constructor(param: TemplateRemove['param']) {
    super()
    this.param = param
  }
}
