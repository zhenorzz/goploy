import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

/* Layout */
import Layout from '@/layout'

/**
 * Note: sub-menu only appear when route children.length >= 1
 * Detail see: https://panjiachen.github.io/vue-element-admin-site/guide/essentials/router-and-nav.html
 *
 * hidden: true                   if set true, item will not show in the sidebar(default is false)
 * alwaysShow: true               if set true, will always show the root menu
 *                                if not set alwaysShow, when item has more than one children route,
 *                                it will becomes nested mode, otherwise not show the root menu
 * redirect: noRedirect           if set noRedirect will no redirect in the breadcrumb
 * name:'router-name'             the name is used by <keep-alive> (must set!!!)
 * meta : {
    roles: ['admin', 'manager', 'group-manager', 'member']   control the page roles (you can set multiple roles)
    title: 'title'               the name show in sidebar and breadcrumb (recommend set)
    icon: 'svg-name'             the icon show in the sidebar
    breadcrumb: false            if set false, the item will hidden in breadcrumb(default is true)
    activeMenu: '/example/list'  if set path, the sidebar will highlight the path you set
  }
 */
export const homeRoutes = [
  // 预留常量 permission.js 会修改权限的第一条
  { path: '/', redirect: '/user' }
]
/**
 * constantRoutes
 * a base page that does not have permission requirements
 * all roles can be accessed
 */
export const constantRoutes = [
  {
    path: '/login',
    component: () => import('@/views/login/index'),
    hidden: true
  },

  {
    path: '/404',
    component: () => import('@/views/404'),
    hidden: true
  },
  {
    path: '/user',
    component: Layout,
    hidden: true,
    redirect: '/user/info',
    children: [{
      path: 'info',
      name: '个人信息',
      component: () => import('@/views/user/info'),
      meta: { title: '个人信息' }
    }]
  }
]

/**
 * asyncRoutes
 * the routes that need to be dynamically loaded based on user permission_uri
 */
export const asyncRoutes = [
  {
    path: '/deploy',
    component: Layout,
    redirect: '/deploy/publish',
    children: [{
      path: 'publish',
      name: '构建发布',
      component: () => import('@/views/deploy/publish'),
      meta: {
        title: '构建发布',
        icon: 'deploy'
      }
    }]
  },
  {
    path: '/project',
    name: '项目管理',
    component: Layout,
    meta: {
      title: '项目管理',
      icon: 'project',
      roles: ['admin', 'manager', 'group-manager']
    },
    children: [
      {
        path: 'setting',
        name: '项目设置',
        component: () => import('@/views/project/setting'),
        meta: {
          title: '项目设置',
          icon: 'setting',
          roles: ['admin', 'manager', 'group-manager']
        }
      },
      {
        path: 'group',
        name: '分组设置',
        component: () => import('@/views/project/group'),
        meta: {
          title: '分组设置',
          icon: 'list',
          roles: ['admin', 'manager', 'group-manager']
        }
      },
      {
        path: 'template',
        name: '模板设置',
        component: () => import('@/views/project/template'),
        meta: {
          title: '模板设置',
          icon: 'template',
          roles: ['admin', 'manager', 'group-manager']
        }
      },
      {
        path: 'server',
        name: '服务器设置',
        component: () => import('@/views/project/server'),
        meta: {
          title: '服务器设置',
          icon: 'server',
          roles: ['admin', 'manager', 'group-manager']
        }
      }
    ]
  },
  {
    path: '/member',
    component: Layout,
    redirect: '/member/list',
    name: '成员管理',
    meta: {
      title: '成员管理',
      icon: 'user',
      roles: ['admin']
    },
    children: [{
      path: 'list',
      name: '成员列表',
      component: () => import('@/views/member/list'),
      meta: {
        title: '成员列表',
        icon: 'user',
        roles: ['admin']
      }
    }]
  },
  // 404 page must be placed at the end !!!
  { path: '*', redirect: '/404', hidden: true }
]

const createRouter = () => new Router({
  // mode: 'history', // require service support
  scrollBehavior: () => ({ y: 0 }),
  routes: constantRoutes
})

const router = createRouter()

// Detail see: https://github.com/vuejs/vue-router/issues/1234#issuecomment-357941465
export function resetRouter() {
  const newRouter = createRouter()
  router.matcher = newRouter.matcher // reset router
}

export default router
