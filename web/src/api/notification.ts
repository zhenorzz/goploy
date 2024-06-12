import { Request, ID } from './types'

export interface NotificationData {
  id: number
  useBy: string
  type: number
  title: string
  template: string
  insertTime: string
  updateTime: string
}

export class NotificationList extends Request {
  readonly url = '/notification/getList'
  readonly method = 'get'

  public declare datagram: {
    list: NotificationData[]
  }
}

export class NotificationEdit extends Request {
  readonly url = '/notification/edit'
  readonly method = 'put'
  public param: {
    id: number
    template: string
  }
  constructor(param: NotificationEdit['param']) {
    super()
    this.param = param
  }
}
