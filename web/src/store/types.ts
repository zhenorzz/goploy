import { AppState } from './modules/app/types'
import { SettingState } from './modules/setting/types'
import { PermissionState } from './modules/permission/types'
import { UserState } from './modules/user/types'
import { TagsViewState } from './modules/tagsView/types'
import { WebsocketState } from './modules/websocket/types'

export interface RootState {
  app: AppState
  setting: SettingState
  permission: PermissionState
  user: UserState
  tagsView: TagsViewState
  websocket: WebsocketState
}
