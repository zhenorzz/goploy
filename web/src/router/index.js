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
    path: '/redirect',
    component: Layout,
    hidden: true,
    children: [
      {
        path: '/redirect/:path(.*)',
        component: () => import('@/views/redirect/index')
      }
    ]
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
    redirect: '/user/profile',
    children: [{
      path: 'profile',
      name: 'UserProfile',
      component: () => import('@/views/user/profile'),
      meta: { title: 'userProfile' }
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
    redirect: '/deploy/index',
    meta: {
      title: 'deploy',
      icon: 'deploy'
    },
    children: [{
      path: 'index',
      name: 'Deploy',
      component: () => import('@/views/deploy/index'),
      meta: {
        title: 'deploy',
        icon: 'deploy',
        affix: true
      }
    }]
  },
  {
    path: '/toolbox',
    component: Layout,
    redirect: '/toolbox/json',
    meta: {
      title: 'toolbox',
      icon: 'toolbox'
    },
    children: [{
      path: 'json',
      name: 'JSONFormatter',
      component: () => import('@/views/toolbox/json'),
      meta: {
        title: 'json',
        icon: 'json'
      }
    }]
  },
  {
    path: '/monitor',
    component: Layout,
    redirect: '/monitor/index',
    meta: {
      title: 'monitor',
      icon: 'monitor',
      roles: ['admin', 'manager', 'group-manager']
    },
    children: [
      {
        path: 'index',
        name: 'Monitor',
        component: () => import('@/views/monitor/index'),
        meta: {
          title: 'monitor',
          icon: 'monitor',
          roles: ['admin', 'manager', 'group-manager']
        }
      }
    ]
  },
  {
    path: '/project',
    component: Layout,
    redirect: '/project/index',
    meta: {
      title: 'project',
      icon: 'project',
      roles: ['admin', 'manager', 'group-manager']
    },
    children: [
      {
        path: 'index',
        name: 'Project',
        component: () => import('@/views/project/index'),
        meta: {
          title: 'project',
          icon: 'project',
          roles: ['admin', 'manager', 'group-manager']
        }
      }
    ]
  },
  {
    path: '/server',
    component: Layout,
    redirect: '/server/index',
    meta: {
      title: 'server',
      icon: 'server',
      roles: ['admin', 'manager']
    },
    children: [
      {
        path: 'index',
        name: 'Server',
        component: () => import('@/views/server/index'),
        meta: {
          title: 'serverSetting',
          icon: 'setting',
          roles: ['admin', 'manager']
        }
      },
      {
        path: 'template',
        name: 'Template',
        component: () => import('@/views/server/template'),
        meta: {
          title: 'template',
          icon: 'template',
          roles: ['admin', 'manager']
        }
      },
      {
        path: 'crontab',
        name: 'Crontab',
        component: () => import('@/views/server/crontab'),
        meta: {
          title: 'crontab',
          icon: 'crontab',
          roles: ['admin', 'manager']
        }
      }
    ]
  },
  {
    path: '/namespace',
    component: Layout,
    redirect: '/namespace/index',
    meta: {
      title: 'namespace',
      icon: 'namespace',
      roles: ['admin', 'manager']
    },
    children: [
      {
        path: 'index',
        name: 'Namespace',
        component: () => import('@/views/namespace/index'),
        meta: {
          title: 'namespace',
          icon: 'namespace',
          roles: ['admin', 'manager']
        }
      }
    ]
  },
  {
    path: '/member',
    component: Layout,
    redirect: '/member/index',
    meta: {
      title: 'member',
      icon: 'user',
      roles: ['admin']
    },
    children: [{
      path: 'index',
      name: 'Member',
      component: () => import('@/views/member/index'),
      meta: {
        title: 'member',
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
