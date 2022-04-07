import { Module, MutationTree, ActionTree } from 'vuex'
import { homeRoutes, asyncRoutes, constantRoutes } from '@/router'
import { RouteRecordRaw } from 'vue-router'
import { RootState } from '../../types'
import { PermissionState } from './types'
import { Info } from '@/api/user'

/**
 * Use meta.role to determine if the current user has permission
 * @param role
 * @param route
 */
function hasPermission(permissionIds: number[], route: RouteRecordRaw) {
  if (route.meta && route.meta.permissions) {
    return route.meta['permissions'].some((permission) =>
      permissionIds.includes(permission)
    )
  } else {
    return true
  }
}

/**
 * Filter asynchronous routing tables by recursion
 * @param routes asyncRoutes
 * @param role
 */
export function filterAsyncRoutes(
  routes: RouteRecordRaw[],
  permissionIds: number[]
): RouteRecordRaw[] {
  const res: RouteRecordRaw[] = []

  routes.forEach((route) => {
    const tmp = { ...route }
    if (tmp.children) {
      tmp.children = filterAsyncRoutes(tmp.children, permissionIds)
      if (tmp.children.length > 0) {
        res.push(tmp)
      }
    } else {
      if (hasPermission(permissionIds, tmp)) {
        res.push(tmp)
      }
    }
  })

  return res
}

const state: PermissionState = {
  routes: [],
  permissionIds: [],
}

const mutations: MutationTree<PermissionState> = {
  SET_ROUTES: (state, routes) => {
    state.routes = constantRoutes.concat(routes)
  },
  SET_PERMISSION_IDS: (state, permissionIds) => {
    state.permissionIds = permissionIds
  },
}

const actions: ActionTree<PermissionState, RootState> = {
  generateRoutes({ commit }, userInfo: Info['datagram']) {
    return new Promise<RouteRecordRaw[]>((resolve) => {
      let accessRoutes = filterAsyncRoutes(
        asyncRoutes,
        userInfo.namespace.permissionIds
      )
      if (
        accessRoutes.length !== 0 &&
        accessRoutes[0]['children'] &&
        accessRoutes[0]['children'].length !== 0
      ) {
        homeRoutes[0].redirect =
          accessRoutes[0].path + '/' + accessRoutes[0].children[0].path
      }
      accessRoutes = homeRoutes.concat(accessRoutes)

      commit('SET_ROUTES', accessRoutes)
      commit('SET_PERMISSION_IDS', userInfo.namespace.permissionIds)
      resolve(accessRoutes)
    })
  },
}

export default <Module<PermissionState, RootState>>{
  namespaced: true,
  state,
  mutations,
  actions,
}
