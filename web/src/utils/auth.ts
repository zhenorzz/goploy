const TokenKey = 'LOGIN'

export function isLogin(): string | null {
  return localStorage.getItem(TokenKey)
}

export function setLogin(status: string): void {
  localStorage.setItem(TokenKey, status)
}

export function logout(): void {
  localStorage.removeItem(TokenKey)
}
