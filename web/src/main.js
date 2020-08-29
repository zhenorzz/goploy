import Vue from 'vue'
import ElementUI from 'element-ui'

import 'normalize.css/normalize.css' // A modern alternative to CSS resets
import 'element-ui/lib/theme-chalk/index.css'
import '@/styles/index.scss' // global css
import '@/icons' // icon
import '@/permission' // permission control

import App from './App'
import i18n from './lang' // internationalization
import store from './store'
import router from './router'
import mixin from '@/mixin'
import global from '@/global' // global config

Vue.use(ElementUI, {
  i18n: (key, value) => i18n.t(key, value),
  size: 'mini'
})

Vue.config.productionTip = false
Vue.prototype.$global = global
Vue.mixin(mixin)

new Vue({
  el: '#app',
  router,
  store,
  i18n,
  render: h => h(App)
})
