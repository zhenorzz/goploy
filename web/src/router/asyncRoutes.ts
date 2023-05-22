import { RouteRecordRaw } from 'vue-router'
import permission from '@/permission'
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
          permissions: [permission.ShowDeployPage],
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
    },
    children: [
      {
        path: 'index',
        name: 'ProjectIndex',
        component: () => import('@/views/project/manage/index.vue'),
        meta: {
          title: 'project',
          icon: 'project',
          permissions: [permission.ShowProjectPage],
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
    },
    children: [
      {
        path: 'index',
        name: 'MonitorIndex',
        component: () => import('@/views/monitor/index.vue'),
        meta: {
          title: 'monitor',
          icon: 'monitor',
          permissions: [permission.ShowMonitorPage],
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
    },
    children: [
      {
        path: 'index',
        name: 'ServerIndex',
        component: () => import('@/views/server/index.vue'),
        meta: {
          title: 'serverSetting',
          icon: 'server',
          permissions: [permission.ShowServerPage],
        },
      },
      {
        path: 'terminal',
        name: 'ServerTerminal',
        component: () => import('@/views/server/terminal/index.vue'),
        meta: {
          title: 'serverTerminal',
          icon: 'terminal',
          permissions: [permission.ShowTerminalPage],
        },
      },
      {
        path: 'sftp',
        name: 'ServerSFTP',
        component: () => import('@/views/server/sftp/index.vue'),
        meta: {
          title: 'serverSFTP',
          icon: 'sftpManage',
          permissions: [permission.ShowSftpFilePage],
        },
      },
      {
        path: 'script',
        name: 'ServerScript',
        component: () => import('@/views/server/script.vue'),
        meta: {
          title: 'serverScript',
          icon: 'script',
          permissions: [permission.ShowServerScriptPage],
        },
      },
      {
        path: 'process',
        name: 'ServerProcess',
        component: () => import('@/views/server/process.vue'),
        meta: {
          title: 'serverProcess',
          icon: 'processManage',
          permissions: [permission.ShowServerProcessPage],
        },
      },
      {
        path: 'cron',
        name: 'ServerCron',
        component: () => import('@/views/server/cron.vue'),
        meta: {
          title: 'serverCron',
          icon: 'crontab',
          permissions: [permission.ShowCronPage],
        },
      },
      {
        path: 'agent',
        name: 'ServerAgent',
        component: () => import('@/views/server/agent.vue'),
        meta: {
          hidden: true,
          title: 'serverAgent',
          icon: 'monitor',
          permissions: [permission.ShowServerMonitorPage],
        },
      },
      {
        path: 'nginx',
        name: 'ServerNginx',
        component: () => import('@/views/server/nginx/index.vue'),
        meta: {
          title: 'serverNginx',
          icon: 'nginx',
          permissions: [permission.ShowServerNginxPage],
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
    },
    children: [
      {
        path: 'index',
        name: 'NamespaceIndex',
        component: () => import('@/views/namespace/index.vue'),
        meta: {
          title: 'namespaceSetting',
          icon: 'namespaceSetting',
          permissions: [permission.ShowNamespacePage],
        },
      },
      {
        path: 'role',
        name: 'NamespaceRole',
        component: () => import('@/views/namespace/role.vue'),
        meta: {
          title: 'roleSetting',
          icon: 'roleSetting',
          permissions: [permission.ShowRolePage],
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
    },
    children: [
      {
        path: 'index',
        name: 'MemberIndex',
        component: () => import('@/views/member/index.vue'),
        meta: {
          title: 'member',
          icon: 'user',
          permissions: [permission.ShowMemberPage],
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
    },
    children: [
      {
        path: 'loginLog',
        name: 'LoginLog',
        component: () => import('@/views/log/loginLog.vue'),
        meta: {
          title: 'loginLog',
          icon: 'log',
          permissions: [permission.ShowLoginLogPage],
        },
      },
      {
        path: 'operationLog',
        name: 'OperationLog',
        component: () => import('@/views/log/operationLog.vue'),
        meta: {
          title: 'operationLog',
          icon: 'log',
          permissions: [permission.ShowOperationLogPage],
        },
      },
      {
        path: 'publishLog',
        name: 'PublishLog',
        component: () => import('@/views/log/publishLog.vue'),
        meta: {
          title: 'publishLog',
          icon: 'log',
          permissions: [permission.ShowPublishLogPage],
        },
      },
      {
        path: 'sftpLog',
        name: 'SftpLog',
        component: () => import('@/views/log/sftpLog.vue'),
        meta: {
          title: 'sftpLog',
          icon: 'log',
          permissions: [permission.ShowSFTPLogPage],
        },
      },
      {
        path: 'terminalLog',
        name: 'TerminalLog',
        component: () => import('@/views/log/terminalLog.vue'),
        meta: {
          title: 'terminalLog',
          icon: 'log',
          permissions: [permission.ShowTerminalLogPage],
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
