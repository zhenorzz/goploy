import router from './route';
import store from './store';
import NProgress from 'nprogress'; // progress bar
import 'nprogress/nprogress.css';// progress bar style
import {getToken} from '@/utils/auth';
import {Message} from 'element-ui';

NProgress.configure({showSpinner: false});// NProgress Configuration
const whiteList = ['/login'];// no redirect whitelist

router.beforeEach((to, from, next) => {
  NProgress.start(); // start progress bar
  if (getToken()) {
    if (to.path === '/login') {
      next({path: '/'});
      NProgress.done();
    } else {
      if (store.getters.role.length === 0) { // 判断当前用户是否已拉取完user_info信息
        store.dispatch('GetInfo').then((res) => { // 拉取user_info
          const role = res.userInfo.role; // note: roles must be a array! such as: ['editor','develop']
          store.dispatch('GenerateRoutes', role).then(() => { // 根据roles权限生成可访问的路由表
            router.addRoutes(store.getters.addRouters); // 动态添加可访问路由表
            next({...to, replace: true}); // hack方法 确保addRoutes已完成 ,set the replace: true so the navigation will not leave a history record
          });
        }).catch((err) => {
          store.dispatch('FedLogOut').then(() => {
            Message.error(err || 'Verification failed, please login again');
            next({path: '/'});
          });
        });
      } else {
        next();
      }
    }
  } else {
    if (whiteList.indexOf(to.path) !== -1) { // 在免登录白名单，直接进入
      next();
    } else {
      next(`/login?redirect=${to.path}`); // 否则全部重定向到登录页
      NProgress.done();
    }
  }
});

router.afterEach(() => {
  NProgress.done(); // finish progress bar
});
