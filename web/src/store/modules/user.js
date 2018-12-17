import {getToken, setToken} from '@/utils/auth';
const user = {
  state: {
    token: getToken(),
  },

  mutations: {
    SET_TOKEN: (state, token) => {
      state.token = token;
    },
  },

  actions: {
    // 登录
    Login({commit}) {
      setToken(1);
      commit('SET_TOKEN', 1);
    },
  },

};

export default user;
