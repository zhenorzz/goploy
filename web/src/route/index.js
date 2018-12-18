import Vue from 'vue';
import Router from 'vue-router';
import Layout from '@/components/layout';
import HelloWorld from '@/components/HelloWorld';
import BaseAdmin from '@/components/BaseAdmin';
import Login from '@/components/Login';
Vue.use(Router);
export const constantRouterMap = [
  {
    path: '/',
    name: 'Dashboard',
    redirect: 'dashboard',
    component: Layout,
    children: [
      {
        path: 'dashboard',
        name: 'Admin',
        component: BaseAdmin,
      },
    ],
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
  },
  {
    path: '/hello',
    name: 'HelloWorld',
    component: HelloWorld,
  },
];
export const asyncRouterMap = [
  {
    path: '/deploy',
    name: 'Deploy',
    component: Layout,
    children: [
      {
        path: '/deploy/test',
        name: 'Test',
        component: HelloWorld,
      },
    ],
  },
];
export default new Router({
  // mode: 'history', // 后端支持可开
  scrollBehavior: () => ({y: 0}),
  routes: constantRouterMap,
});

