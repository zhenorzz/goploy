export interface AppState {
  sidebar: Sidebar
  device: string
  language: string
}

type Sidebar = {
  opened: boolean
  withoutAnimation: boolean
}
