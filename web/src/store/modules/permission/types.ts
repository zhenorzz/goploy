import { RouteRecordRaw } from 'vue-router'
export interface PermissionState {
  routes: RouteRecordRaw[]
  permissionIds: number[]
}
