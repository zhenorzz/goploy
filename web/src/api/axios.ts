import axios, { AxiosResponse, AxiosRequestConfig, AxiosError } from 'axios'
import {
  NamespaceKey,
  getNamespaceId,
  removeNamespace,
} from '@/utils/namespace'
import { logout } from '@/utils/auth'

// create an axios instance
const service = axios.create({
  baseURL: import.meta.env.VITE_APP_BASE_API, // url = base url + request url
  withCredentials: true, // send cookies when cross-domain requests
  timeout: 60000, // request timeout
})

// request interceptor
service.interceptors.request.use(
  (config: AxiosRequestConfig) => {
    // do something before request is sent
    config.headers[NamespaceKey] = getNamespaceId()
    let href = window.location.hash.slice(1)
    if (href.indexOf('?') > -1) {
      href = href.substring(0, href.indexOf('?'))
    }
    config.headers['Router'] = href
    return config
  },
  (error: AxiosError) => {
    // do something with request error
    console.log(error) // for debug
    return Promise.reject(error)
  }
)

// response interceptor
service.interceptors.response.use(
  /**
   * Determine the request status by custom code
   * Here is just an example
   * You can also judge the status by HTTP Status Code
   */
  (response: AxiosResponse) => {
    const res = response.data
    if (res.code !== 0) {
      // 10000:account disabled;
      // 10001:invaild Token;
      // 10002:namespace invalid;
      // 10086:Token expired;
      if ([10000, 10001, 10086].includes(res.code)) {
        ElMessageBox.confirm(res.message, 'Confirm logout', {
          confirmButtonText: 'Re-login',
          type: 'warning',
        }).then(() => {
          logout()
          location.reload()
        })
        return Promise.reject(res.message)
      } else if (10002 === res.code) {
        removeNamespace()
        return Promise.reject(res.message)
      } else {
        ElMessage({
          message: res.message,
          type: 'error',
          duration: 5 * 1000,
        })
        return Promise.reject(response)
      }
    } else {
      return res
    }
  },
  (error: AxiosError) => {
    console.log('err' + error) // for debug
    ElMessage({
      message: error.message,
      type: 'error',
      duration: 5 * 1000,
    })
    return Promise.reject(error)
  }
)

export default service
