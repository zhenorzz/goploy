export interface ExtLoginCallbackParams {
  account: string
  time: number
  token: string
}

export interface MediaLoginCallbackParams {
  authCode: string
  state: string
}

export interface LoginCallbackParams {
  extLogin?: ExtLoginCallbackParams
  mediaLogin?: MediaLoginCallbackParams
}

function normalizeQueryValue(value: unknown): string {
  if (typeof value === 'string') {
    return value
  }
  if (Array.isArray(value) && typeof value[0] === 'string') {
    return value[0]
  }
  return ''
}

function getCurrentQueryParams(
  toQuery: Record<string, unknown> = {}
): Record<string, string> {
  const params: Record<string, string> = {}

  const searchParams = new URLSearchParams(window.location.search)
  for (const [key, value] of searchParams.entries()) {
    params[key] = value
  }

  const hash = window.location.hash
  const hashQueryIndex = hash.indexOf('?')
  if (hashQueryIndex > -1) {
    const hashParams = new URLSearchParams(hash.slice(hashQueryIndex + 1))
    for (const [key, value] of hashParams.entries()) {
      params[key] = value
    }
  }

  for (const [key, value] of Object.entries(toQuery)) {
    const normalizedValue = normalizeQueryValue(value)
    if (normalizedValue !== '') {
      params[key] = normalizedValue
    }
  }

  return params
}

export function parseLoginCallbackParams(
  toQuery: Record<string, unknown> = {}
): LoginCallbackParams | null {
  const params = getCurrentQueryParams(toQuery)

  if (params.account && params.time && params.token) {
    return {
      extLogin: {
        account: params.account,
        time: Number(params.time),
        token: params.token,
      },
    }
  }

  const authCode = params.code || params.authCode
  if (authCode && params.state) {
    return {
      mediaLogin: {
        authCode,
        state: params.state,
      },
    }
  }

  return null
}

export function hasLoginCallbackQuery(
  toQuery: Record<string, unknown> = {}
): boolean {
  return parseLoginCallbackParams(toQuery) !== null
}
