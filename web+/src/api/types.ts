// 接口响应通过格式
import Axios from './axios'
import { Method, AxiosRequestConfig } from 'axios'

export interface HttpResponse<T> {
  code: number
  message: string
  data: T
}

export interface Pagination {
  page: number
  rows: number
  total: number
}

export interface Total {
  total: number
}

export interface ID {
  id: number
}

export abstract class Request {
  abstract url: string
  abstract method: Method
  public timeout = 5000
  public param!: ID | Record<string, unknown>
  public datagram!: any

  public request(): Promise<HttpResponse<this['datagram']>> {
    const config: AxiosRequestConfig = {
      url: this.url,
      method: this.method,
      timeout: this.timeout,
    }
    if (this.method.toLowerCase() === 'get') {
      config.params = { ...this.param }
    } else {
      config.data = { ...this.param }
    }
    return Axios.request(config)
  }
}
