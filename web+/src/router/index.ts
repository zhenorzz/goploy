import { RouteRecordRaw, createRouter, createWebHashHistory } from 'vue-router'

/* Layout */
import Layout from '@/layout/index.vue'

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
export const homeRoutes: RouteRecordRaw[] = [
  // 预留常量 permission.js 会修改权限的第一条
  { path: '/', redirect: '/user' },
]
/**
 * constantRoutes
 * a base page that does not have permission requirements
 * all roles can be accessed
 */
export const constantRoutes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/login/index.vue'),
    meta: { hidden: true },
  },
  {
    path: '/redirect',
    name: 'redirect',
    component: Layout,
    meta: { hidden: true },
    children: [
      {
        path: '/redirect/:path(.*)',
        component: () => import('@/views/redirect/index.vue'),
      },
    ],
  },
  {
    path: '/404',
    name: '404',
    component: () => import('@/views/404.vue'),
    meta: { hidden: true },
  },
  {
    path: '/user',
    name: 'user',
    component: Layout,
    redirect: '/user/profile',
    meta: { hidden: true },
    children: [
      {
        path: 'profile',
        name: 'UserProfile',
        component: () => import('@/views/user/profile.vue'),
        meta: { title: 'userProfile' },
      },
    ],
  },
]

/**
 * asyncRoutes
 * the routes that need to be dynamically loaded based on user permission_uri
 */
export const asyncRoutes: RouteRecordRaw[] = [
  {
    path: '/deploy',
    name: 'deploy',
    component: Layout,
    redirect: '/deploy/index',
    meta: {
      title: 'deploy',
      icon: 'deploy',
    },
    children: [
      {
        path: 'index',
        name: 'DeployIndex',
        component: () => import('@/views/deploy/index.vue'),
        meta: {
          title: 'deploy',
          icon: 'deploy',
          affix: true,
        },
      },
    ],
  },
  {
    path: '/monitor',
    name: 'monitor',
    component: Layout,
    redirect: '/monitor/index',
    meta: {
      title: 'monitor',
      icon: 'monitor',
      roles: ['admin', 'manager', 'group-manager'],
    },
    children: [
      {
        path: 'index',
        name: 'MonitorIndex',
        component: () => import('@/views/monitor/index.vue'),
        meta: {
          title: 'monitor',
          icon: 'monitor',
          roles: ['admin', 'manager', 'group-manager'],
        },
      },
    ],
  },
  {
    path: '/project',
    name: 'project',
    component: Layout,
    redirect: '/project/index',
    meta: {
      title: 'project',
      icon: 'project',
      roles: ['admin', 'manager', 'group-manager'],
    },
    children: [
      {
        path: 'index',
        name: 'ProjectIndex',
        component: () => import('@/views/project/index.vue'),
        meta: {
          title: 'project',
          icon: 'project',
          roles: ['admin', 'manager', 'group-manager'],
        },
      },
    ],
  },
  {
    path: '/server',
    name: 'server',
    component: Layout,
    redirect: '/server/index',
    meta: {
      title: 'server',
      icon: 'server',
      roles: ['admin', 'manager'],
    },
    children: [
      {
        path: 'index',
        name: 'ServerIndex',
        component: () => import('@/views/server/manage/index.vue'),
        meta: {
          title: 'serverSetting',
          icon: 'setting',
          roles: ['admin', 'manager'],
        },
      },
      {
        path: 'crontab',
        name: 'Crontab',
        component: () => import('@/views/server/crontab/index.vue'),
        meta: {
          title: 'crontab',
          icon: 'crontab',
          roles: ['admin', 'manager'],
        },
      },
    ],
  },
  {
    path: '/namespace',
    component: Layout,
    redirect: '/namespace/index',
    name: 'namespace',
    meta: {
      title: 'namespace',
      icon: 'namespace',
      roles: ['admin', 'manager'],
    },
    children: [
      {
        path: 'index',
        name: 'NamespaceIndex',
        component: () => import('@/views/namespace/index.vue'),
        meta: {
          title: 'namespace',
          icon: 'namespace',
          roles: ['admin', 'manager'],
        },
      },
    ],
  },
  {
    path: '/member',
    component: Layout,
    redirect: '/member/index',
    name: 'member',
    meta: {
      title: 'member',
      icon: 'user',
      roles: ['admin'],
    },
    children: [
      {
        path: 'index',
        name: 'MemberIndex',
        component: () => import('@/views/member/index.vue'),
        meta: {
          title: 'member',
          icon: 'user',
          roles: ['admin'],
        },
      },
    ],
  },
  // 404 page must be placed at the end !!!
  {
    path: '/:pathMatch(.*)*',
    name: '404*',
    redirect: '/404',
    meta: { hidden: true },
  },
]

const router = createRouter({
  history: createWebHashHistory(),
  scrollBehavior() {
    return {
      el: '#app',
      left: 0,
      behavior: 'smooth',
    }
  },
  routes: constantRoutes,
})

export function resetRouter(): void {
  router
    .getRoutes()
    .forEach((route) => route.name && router.removeRoute(route.name))
  constantRoutes.forEach((route: RouteRecordRaw) => router.addRoute(route))
}

export default router
