import { homeRoutes, asyncRoutes, constantRoutes } from '@/router'
/**
 * 通过meta.permissionMap判断是否与当前用户权限匹配
 * @param router
 * @param permissionMap
 */
function hasPermission(router, permissionMap) {
  let hasPermission = false
  if (router.meta && router.meta.permission_uri) {
    const uri = router.meta.permission_uri
    let index
    for (index in permissionMap) {
      if (!permissionMap.hasOwnProperty(index)) continue
      if (permissionMap[index].uri === uri) {
        hasPermission = index
        break
      }
    }
  }
  return hasPermission
}

/**
 * 递归过滤异步路由表，返回符合用户角色权限的路由表
 * @param asyncRouterMap asyncRouterMap
 * @param resMap
 */
function filterAsyncRoutes(asyncRouterMap, resMap) {
  return asyncRouterMap.filter(route => {
    if (route.hasOwnProperty('meta') && route.meta.permission_uri === 'all') {
      if (route.children && route.children.length) {
        route.children = filterAsyncRoutes(route.children, resMap)
      }
      return true
    }
    const mapIndex = hasPermission(route, resMap)
    if (mapIndex !== false) {
      if (route.children && route.children.length) {
        route.children = filterAsyncRoutes(route.children, resMap[mapIndex].children)
      }
      return true
    } else {
      return false
    }
  })
}

const state = {
  routes: [],
  addRouters: [],
  permissionUri: []
}

const mutations = {
  SET_ROUTES: (state, routes) => {
    state.addRoutes = routes
    state.routes = constantRoutes.concat(routes)
  },
  SET_URI: (state, permissionUri) => {
    state.permissionUri = permissionUri
  }
}

const actions = {
  generateRoutes({ commit }, data) {
    return new Promise(resolve => {
      let accessRoutes = filterAsyncRoutes(asyncRoutes, data.permission)
      if (accessRoutes.length !== 0) {
        homeRoutes[0].redirect = accessRoutes[0].path + '/' + accessRoutes[0].children[0].path
      }
      accessRoutes = homeRoutes.concat(accessRoutes)

      commit('SET_ROUTES', accessRoutes)
      commit('SET_URI', data.permissionUri)

      resolve(accessRoutes)
    })
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
