export const NamespaceKey = 'G-N-ID'
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
  Namespace: getNamespace(),
  toString(): string {
    return role.Namespace['role']
  },
  isAdmin(): boolean {
    return role.Admin === role.Namespace['role']
  },
  isManager(): boolean {
    return role.Manager === role.Namespace['role']
  },
  isGroupManager(): boolean {
    return role.GroupManager === role.Namespace['role']
  },
  isMember(): boolean {
    return role.Member === role.Namespace['role']
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
  const namespaceId = getNamespaceId()
  const namespaceList = getNamespaceList()
  if (namespaceId && namespaceList) {
    return namespaceList.find(
      (_) => _.id.toString() === namespaceId
    ) as Namespace
  }
  return { id: 0, name: '', role: '' }
}

export function getNamespaceId(): string | undefined {
  const namespaceId =
    sessionStorage.getItem(NamespaceKey) || localStorage.getItem(NamespaceKey)
  return namespaceId || undefined
}

export function setNamespaceId(namespaceId: string): void {
  sessionStorage.setItem(NamespaceKey, namespaceId)
  localStorage.setItem(NamespaceKey, namespaceId)
}

export function removeNamespaceId(): void {
  sessionStorage.removeItem(NamespaceKey)
  localStorage.removeItem(NamespaceKey)
}

export function getNamespaceList(): Array<Namespace> {
  const namespaceList = localStorage.getItem(NamespaceListKey)
  return namespaceList ? JSON.parse(namespaceList) : []
}

export function setNamespaceList(namespaceList: Array<Namespace>): void {
  localStorage.setItem(NamespaceListKey, JSON.stringify(namespaceList))
}
