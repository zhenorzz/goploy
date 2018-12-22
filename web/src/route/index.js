import Vue from 'vue';
import Router from 'vue-router';
import Layout from '@/components/layout';
import BaseAdmin from '@/components/BaseAdmin';
import Login from '@/components/Login';
Vue.use(Router);
export const constantRouterMap = [
  {
    path: '/',
    name: 'Dashboard',
    redirect: 'dashboard',
    component: Layout,
    meta: {
      title: 'eye',
      icon: 'eye',
    },
    children: [
      {
        path: 'dashboard',
        name: 'Admin',
        component: BaseAdmin,
        meta: {
          title: 'dashboard',
          icon: 'eye',
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
    path: '/deploy',
    name: 'Deploy',
    component: Layout,
    meta: {
      title: 'father',
      icon: 'eye',
      roles: ['admin'],
    },
    children: [
      {
        path: '/deploy/test',
        name: 'Test',
        component: () => import('@/components/HelloWorld'),
        meta: {
          title: 'son1',
          icon: 'eye',
          roles: ['admin'],
        },
      },
      {
        path: '/deploy/my',
        name: 'test1',
        component: () => import('@/components/BaseAdmin'),
        meta: {
          title: 'son2',
          icon: 'eye',
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

