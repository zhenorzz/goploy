import Vue from 'vue';
import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';
import App from './App';
import router from './route';
import store from './store';
import '@/icons'; // permission control
import '@/permission'; // permission control
Vue.config.productionTip = false;
Vue.use(ElementUI);
new Vue({
  store,
  router,
  render: (h) => h(App),
}).$mount('#app');
