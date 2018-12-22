import Vue from 'vue';
import Vuex from 'vuex';
import user from './modules/user';
import permission from './modules/permission';
import getters from './getters';
Vue.use(Vuex);
const store = new Vuex.Store({
  modules: {
    permission,
    user,
  },
  getters,
});

export default store;

