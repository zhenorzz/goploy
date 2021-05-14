import Cookies from 'js-cookie'
const NamespaceKey = 'goploy_namespace'
const NamespaceListKey = 'goploy_namespace_list'

interface Namespace {
  id: number
  name: string
  role: string
}

export const role = Object.freeze({
  Admin: 'admin',
  Manager: 'manager',
  GroupManager: 'group-manager',
  Member: 'member',
  toString(): string {
    return getNamespace()['role']
  },
  isAdmin(): boolean {
    return role.Admin === getNamespace()['role']
  },
  isManager(): boolean {
    return role.Manager === getNamespace()['role']
  },
  isGroupManager(): boolean {
    return role.GroupManager === getNamespace()['role']
  },
  isMember(): boolean {
    return role.Member === getNamespace()['role']
  },
  hasAdminPermission(): boolean {
    return role.isAdmin()
  },
  hasManagerPermission(): boolean {
    return role.isAdmin() || role.isManager()
  },
  hasGroupManagerPermission(): boolean {
    return role.isAdmin() || role.isManager() || role.isGroupManager()
  },
})

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
