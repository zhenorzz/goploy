import { RouteRecordRaw } from 'vue-router'

/* Layout */
import Layout from '@/layout/index.vue'

/**
 * asyncRoutes
 * the routes that need to be dynamically loaded based on user permission_uri
 */
export default <RouteRecordRaw[]>[
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
        component: () => import('@/views/project/manage/index.vue'),
        meta: {
          title: 'project',
          icon: 'project',
          roles: ['admin', 'manager', 'group-manager'],
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
          icon: 'server',
          roles: ['admin', 'manager'],
        },
      },
      {
        path: 'terminal',
        name: 'ServerTerminal',
        component: () => import('@/views/server/terminal/index.vue'),
        meta: {
          title: 'serverTerminal',
          icon: 'terminal',
          roles: ['admin', 'manager'],
        },
      },
      {
        path: 'sftp',
        name: 'ServerSFTP',
        component: () => import('@/views/server/sftp/index.vue'),
        meta: {
          title: 'serverSFTP',
          icon: 'ftp',
          roles: ['admin', 'manager'],
        },
      },
      {
        path: 'agent',
        name: 'ServerAgent',
        component: () => import('@/views/server/agent/index.vue'),
        meta: {
          hidden: true,
          title: 'serverAgent',
          icon: 'monitor',
          roles: ['admin', 'manager'],
        },
      },
      {
        path: 'cron',
        name: 'ServerCron',
        component: () => import('@/views/server/cron/index.vue'),
        meta: {
          title: 'serverCron',
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
          title: 'namespaceSetting',
          icon: 'namespaceSetting',
          roles: ['admin', 'manager'],
        },
      },
      {
        path: 'role',
        name: 'NamespaceRole',
        component: () => import('@/views/namespace/role.vue'),
        meta: {
          title: 'roleSetting',
          icon: 'roleSetting',
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
  {
    path: '/log',
    component: Layout,
    redirect: '/log/loginLog',
    name: 'log',
    meta: {
      title: 'log',
      icon: 'log',
      roles: ['admin'],
    },
    children: [
      {
        path: 'loginLog',
        name: 'LoginLog',
        component: () => import('@/views/log/loginLog.vue'),
        meta: {
          title: 'loginLog',
          icon: 'log',
          roles: ['admin'],
        },
      },
      {
        path: 'publishLog',
        name: 'PublishLog',
        component: () => import('@/views/log/publishLog.vue'),
        meta: {
          title: 'publishLog',
          icon: 'log',
          roles: ['admin'],
        },
      },
      {
        path: 'sftpLog',
        name: 'SftpLog',
        component: () => import('@/views/log/sftpLog.vue'),
        meta: {
          title: 'sftpLog',
          icon: 'log',
          roles: ['admin'],
        },
      },
      {
        path: 'terminalLog',
        name: 'TerminalLog',
        component: () => import('@/views/log/terminalLog.vue'),
        meta: {
          title: 'terminalLog',
          icon: 'log',
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
