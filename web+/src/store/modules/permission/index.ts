import { Module, MutationTree, ActionTree } from 'vuex'
import { homeRoutes, asyncRoutes, constantRoutes } from '@/router'
import { RouteRecordRaw } from 'vue-router'
import { getNamespace } from '@/utils/namespace'
import { RootState } from '../../types'
import { PermissionState } from './types'

/**
 * Use meta.role to determine if the current user has permission
 * @param role
 * @param route
 */
function hasPermission(role: string, route: RouteRecordRaw) {
  if (route.meta && route.meta.roles) {
    return route.meta['roles'].includes(role)
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
  role: string
): RouteRecordRaw[] {
  const res: RouteRecordRaw[] = []

  routes.forEach((route) => {
    const tmp = { ...route }
    if (hasPermission(role, tmp)) {
      if (tmp.children) {
        tmp.children = filterAsyncRoutes(tmp.children, role)
      }
      res.push(tmp)
    }
  })

  return res
}

const state: PermissionState = {
  routes: [],
}

const mutations: MutationTree<PermissionState> = {
  SET_ROUTES: (state, routes) => {
    state.routes = constantRoutes.concat(routes)
  },
}

const actions: ActionTree<PermissionState, RootState> = {
  generateRoutes({ commit }) {
    return new Promise((resolve) => {
      const namespace = getNamespace()
      let accessRoutes = filterAsyncRoutes(asyncRoutes, namespace.role)
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
      resolve(accessRoutes)
    })
  },
}

const permission: Module<PermissionState, RootState> = {
  namespaced: true,
  state,
  mutations,
  actions,
}

export default permission
