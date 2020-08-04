import Cookies from 'js-cookie'
const NamespaceKey = 'goploy_namespace'
const NamespaceListKey = 'goploy_namespace_list'

export function getNamespace() {
  const namespace = localStorage.getItem(NamespaceKey)
  return namespace ? JSON.parse(namespace) : false
}

export function setNamespace(namespace) {
  Cookies.set(NamespaceKey, namespace.id)
  return localStorage.setItem(NamespaceKey, JSON.stringify(namespace))
}

export function removeNamespace() {
  Cookies.remove(NamespaceKey)
  return localStorage.removeItem(NamespaceKey)
}

export function getNamespaceList() {
  const namespace = localStorage.getItem(NamespaceListKey)
  return namespace ? JSON.parse(namespace) : []
}

export function setNamespaceList(namespace) {
  return localStorage.setItem(NamespaceListKey, JSON.stringify(namespace))
}
