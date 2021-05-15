// 接口响应通过格式
export interface HttpResponse<T> {
  code: number
  message: string
  data: T
}

export interface Pagination {
  page: number
  rows: number
}

export interface Total {
  total: number
}

export interface ID {
  id: number
}
