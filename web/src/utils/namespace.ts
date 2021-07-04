import Cookies from 'js-cookie'
const NamespaceKey = 'goploy_namespace'
const NamespaceListKey = 'goploy_namespace_list'

export interface Namespace {
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
  localStorage.setItem(NamespaceKey, JSON.stringify(namespace))
}

export function getNamespaceIdCookie(): string | undefined {
  return Cookies.get(NamespaceKey)
}

export function setNamespaceIdCookie(namespaceId: string): void {
  Cookies.set(NamespaceKey, namespaceId, { expires: 365 })
}

export function removeNamespaceIdCookie(): void {
  Cookies.remove(NamespaceKey)
}

export function getNamespaceList(): Array<Namespace> {
  const namespaceList = localStorage.getItem(NamespaceListKey)
  return namespaceList ? JSON.parse(namespaceList) : []
}

export function setNamespaceList(namespaceList: Array<Namespace>): void {
  localStorage.setItem(NamespaceListKey, JSON.stringify(namespaceList))
}
