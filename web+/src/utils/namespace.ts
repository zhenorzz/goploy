import Cookies from 'js-cookie'
const NamespaceKey = 'goploy_namespace'
const NamespaceListKey = 'goploy_namespace_list'

interface Namespace {
  id: number
  name: string
  role: string
}

export function getNamespace(): Namespace {
  const namespace = localStorage.getItem(NamespaceKey)
  return namespace ? JSON.parse(namespace) : undefined
}

export function setNamespace(namespace: Namespace): void {
  Cookies.set(NamespaceKey, namespace.id.toString(), { expires: 365 })
  localStorage.setItem(NamespaceKey, JSON.stringify(namespace))
}

export function removeNamespace(): void {
  Cookies.remove(NamespaceKey)
  localStorage.removeItem(NamespaceKey)
}

export function getNamespaceList(): Array<Namespace> {
  const namespaceList = localStorage.getItem(NamespaceListKey)
  return namespaceList ? JSON.parse(namespaceList) : []
}

export function setNamespaceList(namespaceList: Array<Namespace>): void {
  localStorage.setItem(NamespaceListKey, JSON.stringify(namespaceList))
}
