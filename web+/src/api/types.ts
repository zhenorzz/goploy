// 接口响应通过格式
export type HttpResponse<T> = {
  code: number
  message: string
  data: T
}

export type Pagination = {
  page: number
  rows: number
}

export type Total = {
  total: number
}

export type ID = {
  id: number
}
