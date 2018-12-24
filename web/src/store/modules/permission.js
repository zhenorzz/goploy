import {asyncRouterMap, constantRouterMap} from '@/route';

/**
 * 通过meta.role判断是否与当前用户权限匹配
 * @param {string} role
 * @param {object} route
 * @return {boolean}
 */
function hasPermission(role, route) {
  if (route.meta && route.meta.roles) {
    return route.meta.roles.includes(role);
  } else {
    return true;
  }
}

/**
 * 递归过滤异步路由表，返回符合用户角色权限的路由表
 * @param {object} routes asyncRouterMap
 * @param {string} role
 * @return {object}
 */
function filterAsyncRouter(routes, role) {
  const res = [];
  routes.forEach((route) => {
    const tmp = {...route};
    if (hasPermission(role, tmp)) {
      if (tmp.children) {
        tmp.children = filterAsyncRouter(tmp.children, role);
      }
      res.push(tmp);
    }
  });

  return res;
}

const permission = {
  state: {
    routers: constantRouterMap,
    addRouters: [],
  },
  mutations: {
    SET_ROUTERS: (state, routers) => {
      state.addRouters = routers;
      state.routers = constantRouterMap.concat(routers);
    },
  },
  actions: {
    GenerateRoutes({commit}, role) {
      return new Promise((resolve) => {
        let accessedRouters;
        // if (role === 'admin') {
        //   accessedRouters = asyncRouterMap;
        // } else {
        accessedRouters = filterAsyncRouter(asyncRouterMap, role);
        // }
        commit('SET_ROUTERS', accessedRouters);
        resolve();
      });
    },
  },
};

export default permission;
