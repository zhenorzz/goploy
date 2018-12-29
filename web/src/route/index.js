import Vue from 'vue';
import Router from 'vue-router';
import Layout from '@/components/layout';
import Login from '@/components/Login';
Vue.use(Router);
export const constantRouterMap = [
  {
    path: '/',
    name: 'Home',
    redirect: 'home',
    component: Layout,
    meta: {
      title: '主页',
      icon: 'home',
    },
    children: [
      {
        path: 'home',
        name: '主页',
        component: () => import('@/components/home'),
        meta: {
          title: '主页',
          icon: 'home',
        },
      },
    ],
  },
  {
    path: '/login',
    component: Login,
    hidden: true,
  },
];
export const asyncRouterMap = [
  {
    path: '/project',
    name: '项目管理',
    component: Layout,
    meta: {
      title: '项目管理',
      icon: 'project',
      roles: ['admin'],
    },
    children: [
      {
        path: '/project/deploy',
        name: '项目部署',
        component: () => import('@/components/project/deploy'),
        meta: {
          title: '项目部署',
          icon: 'deploy',
          roles: ['admin'],
        },
      },
      {
        path: '/project/setting',
        name: '项目设置',
        component: () => import('@/components/project/setting'),
        meta: {
          title: '项目设置',
          icon: 'setting',
          roles: ['admin'],
        },
      },
    ],
  },
  {
    path: '/server',
    name: '服务器管理',
    component: Layout,
    meta: {
      title: '服务器管理',
      icon: 'server',
      roles: ['admin'],
    },
    children: [
      {
        path: '/server/list',
        name: '服务器列表',
        component: () => import('@/components/server/list'),
        meta: {
          title: '服务器列表',
          icon: 'list',
          roles: ['admin'],
        },
      },
      {
        path: '/server/setting',
        name: '服务器设置',
        component: () => import('@/components/HelloWorld'),
        meta: {
          title: '服务器设置',
          icon: 'setting',
          roles: ['admin'],
        },
      },
    ],
  },
  {
    path: '/user',
    name: '成员管理',
    component: Layout,
    meta: {
      title: '成员管理',
      icon: 'user',
      roles: ['admin'],
    },
    children: [
      {
        path: '/user/list',
        name: '成员列表',
        component: () => import('@/components/BaseAdmin'),
        meta: {
          title: '成员列表',
          icon: 'list',
          roles: ['admin'],
        },
      },
    ],
  },
];
export default new Router({
  // mode: 'history', // 后端支持可开
  scrollBehavior: () => ({y: 0}),
  routes: constantRouterMap,
});

