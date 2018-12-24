const TokenKey = 'Admin-Token';

/**
 * @return {string}
 */
export function getToken() {
  return localStorage.getItem(TokenKey);
}

/**
 * @param  {string} token
 * @return {string}
 */
export function setToken(token) {
  return localStorage.setItem(TokenKey, token);
}

/**
 * @return {string}
 */
export function removeToken() {
  return localStorage.removeItem(TokenKey);
}
