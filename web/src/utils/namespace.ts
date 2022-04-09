export const NamespaceKey = 'G-N-ID'

export interface Namespace {
  id: number
  name: string
  roleId: number
}

export function getNamespace(): Namespace {
  const n =
    sessionStorage.getItem(NamespaceKey) || localStorage.getItem(NamespaceKey)

  try {
    return n ? JSON.parse(n) : { id: 0, name: '', roleId: -1 }
  } catch (e) {
    return { id: 0, name: '', roleId: -1 }
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
