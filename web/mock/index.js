import Mock from 'mockjs'
import { param2Obj } from '../src/utils'

const tokens = {
  admin: {
    token: 'admin'
  }
}

const users = {
  'admin': {
    id: 1,
    name: 'admin',
    account: 'admin',
    role: 'admin'
  }
}

const user = [
  // user login
  {
    url: '/user/login',
    type: 'post',
    response: config => {
      const { account } = config.body
      const token = tokens[account]
      // mock error
      if (!token) {
        return {
          code: 1,
          message: 'Account and password are incorrect.'
        }
      }

      return {
        code: 0,
        data: token
      }
    }
  },

  // get user info
  {
    url: '/user/info',
    type: 'get',
    response: _ => {
      const info = users.admin
      // mock error
      if (!info) {
        return {
          code: 1,
          message: 'Login failed, unable to get user details.'
        }
      }
      return {
        code: 0,
        data: {
          userInfo: info
        }
      }
    }
  },

  // user logout
  {
    url: '/user/logout',
    type: 'post',
    response: _ => {
      return {
        code: 0,
        data: 'success'
      }
    }
  },
  {
    url: '\.*',
    type: 'get',
    response: _ => {
      return {
        code: 0,
        data: 'success'
      }
    }
  },
  {
    url: '\.*',
    type: 'post',
    response: _ => {
      return {
        code: 0,
        data: 'success'
      }
    }
  }
]

const mocks = [
  ...user
]

// for front mock
// please use it cautiously, it will redefine XMLHttpRequest,
// which will cause many of your third-party libraries to be invalidated(like progress event).
export function mockXHR() {
  // mock patch
  // https://github.com/nuysoft/Mock/issues/300
  Mock.XHR.prototype.proxy_send = Mock.XHR.prototype.send
  Mock.XHR.prototype.send = function() {
    if (this.custom.xhr) {
      this.custom.xhr.withCredentials = this.withCredentials || false

      if (this.responseType) {
        this.custom.xhr.responseType = this.responseType
      }
    }
    this.proxy_send(...arguments)
  }

  function XHR2ExpressReqWrap(respond) {
    return function(options) {
      let result = null
      if (respond instanceof Function) {
        const { body, type, url } = options
        // https://expressjs.com/en/4x/api.html#req
        result = respond({
          method: type,
          body: JSON.parse(body),
          query: param2Obj(url)
        })
      } else {
        result = respond
      }
      return Mock.mock(result)
    }
  }

  for (const i of mocks) {
    Mock.mock(new RegExp(i.url), i.type || 'get', XHR2ExpressReqWrap(i.response))
  }
}

// for mock server
const responseFake = (url, type, respond) => {
  return {
    url: new RegExp(`/mock${url}`),
    type: type || 'get',
    response(req, res) {
      res.json(Mock.mock(respond instanceof Function ? respond(req, res) : respond))
    }
  }
}

export default mocks.map(route => {
  return responseFake(route.url, route.type, route.response)
})
