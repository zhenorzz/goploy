export const NamespaceKey = 'G-N-ID'

export interface Namespace {
  id: number
  name: string
  role: string
}

export function getRole() {
  return {
    Admin: 'admin',
    Manager: 'manager',
    GroupManager: 'group-manager',
    Member: 'member',
    Namespace: getNamespace(),
    toString(): string {
      return this.Namespace['role']
    },
    isAdmin(): boolean {
      return this.Admin === this.Namespace['role']
    },
    isManager(): boolean {
      return this.Manager === this.Namespace['role']
    },
    isGroupManager(): boolean {
      return this.GroupManager === this.Namespace['role']
    },
    isMember(): boolean {
      return this.Member === this.Namespace['role']
    },
    hasAdminPermission(): boolean {
      return this.isAdmin()
    },
    hasManagerPermission(): boolean {
      return this.isAdmin() || this.isManager()
    },
    hasGroupManagerPermission(): boolean {
      return this.isAdmin() || this.isManager() || this.isGroupManager()
    },
  }
}

export function getNamespace(): Namespace {
  const n =
    sessionStorage.getItem(NamespaceKey) || localStorage.getItem(NamespaceKey)

  try {
    return n ? JSON.parse(n) : { id: 0, name: '', role: '' }
  } catch (e) {
    return { id: 0, name: '', role: '' }
  }
}

export function getNamespaceId(): number {
  const namespaceId = getNamespace().id
  return namespaceId
}

export function setNamespace(namespace: Namespace): void {
  sessionStorage.setItem(NamespaceKey, JSON.stringify(namespace))
  localStorage.setItem(NamespaceKey, JSON.stringify(namespace))
}

export function removeNamespace(): void {
  sessionStorage.removeItem(NamespaceKey)
  localStorage.removeItem(NamespaceKey)
}
