import { Request, ID } from './types'

export interface UserData {
  account: string
  contact: string
  id: number
  insertTime: string
  lastLoginTime: string
  name: string
  password: string
  state: number
  superManager: number
  updateTime: string
}

export class Login extends Request {
  readonly url = '/user/login'
  readonly method = 'post'
  public param: {
    account: string
    password: string
    newPassword: string
    captchaKey: string
  }
  public declare datagram: {
    namespaceList: { id: number; name: string; roleId: number }[]
  }
  constructor(param: Login['param']) {
    super()
    this.param = param
  }
}

export class extLogin extends Request {
  readonly url = '/user/extLogin'
  readonly method = 'post'
  public param: {
    account: string
    token: string
    time: number
  }
  public declare datagram: {
    namespaceList: { id: number; name: string; roleId: number }[]
  }
  constructor(param: extLogin['param']) {
    super()
    this.param = param
  }
}

export class MediaLogin extends Request {
  readonly url = '/user/mediaLogin'
  readonly method = 'post'
  public param: {
    authCode: string
    state: string
    redirectUri: string
  }
  public declare datagram: {
    namespaceList: { id: number; name: string; roleId: number }[]
  }
  constructor(param: MediaLogin['param']) {
    super()
    this.param = param
  }
}

export class Info extends Request {
  readonly url = '/user/info'
  readonly method = 'get'
  public declare datagram: {
    userInfo: {
      account: string
      id: number
      name: string
      superManager: number
    }
    namespace: {
      id: number
      permissionIds: number[]
    }
  }
}

export class UserList extends Request {
  readonly url = '/user/getList'
  readonly method = 'get'
  public declare datagram: {
    list: UserData[]
  }
}

export class UserOption extends Request {
  readonly url = '/user/getList'
  readonly method = 'get'
  public declare datagram: {
    list: UserData[]
  }
}

export class UserAdd extends Request {
  readonly url = '/user/add'
  readonly method = 'post'
  public param: {
    account: string
    password: string
    name: string
    contact: string
    superManager: number
  }
  constructor(param: UserAdd['param']) {
    super()
    this.param = param
  }
}

export class UserEdit extends Request {
  readonly url = '/user/edit'
  readonly method = 'put'
  public param: {
    id: number
    password: string
    name: string
    contact: string
    superManager: number
  }
  constructor(param: UserEdit['param']) {
    super()
    this.param = param
  }
}

export class UserRemove extends Request {
  readonly url = '/user/remove'
  readonly method = 'delete'
  public param: ID
  constructor(param: ID) {
    super()
    this.param = param
  }
}

export class UserGetApiKey extends Request {
  readonly url = '/user/getApiKey'
  readonly method = 'get'
}

export class UserGenerateApiKey extends Request {
  readonly url = '/user/generateApiKey'
  readonly method = 'put'
}

export class UserChangePassword extends Request {
  readonly url = '/user/changePassword'
  readonly method = 'put'
  public param: {
    oldPwd: string
    newPwd: string
  }
  constructor(param: UserChangePassword['param']) {
    super()
    this.param = param
  }
}

export class GetConfig extends Request {
  readonly url = '/user/getConfig'
  readonly method = 'get'
  public declare datagram: {
    ldap: {
      enabled: boolean
    }
    captcha: {
      enabled: boolean
    }
    mediaURL: {
      dingtalk: string
      feishu: string
    }
  }
}

export class GetCaptcha extends Request {
  readonly url = '/user/getCaptcha'
  readonly method = 'get'
  public param: {
    language: any
  }
  public declare datagram: {
    base64: string
    thumbBase64: string
    key: string
  }
  constructor(param: GetCaptcha['param']) {
    super()
    this.param = param
  }
}

export class CheckCaptcha extends Request {
  readonly url = '/user/checkCaptcha'
  readonly method = 'post'
  public param: {
    dots: number[]
    key: string
  }
  constructor(param: CheckCaptcha['param']) {
    super()
    this.param = param
  }
}
