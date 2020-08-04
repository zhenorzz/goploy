const TokenKey = 'LOGIN'

export function isLogin() {
  return localStorage.getItem(TokenKey)
}

export function setLogin(status) {
  return localStorage.setItem(TokenKey, status)
}

export function logout() {
  return localStorage.removeItem(TokenKey)
}
