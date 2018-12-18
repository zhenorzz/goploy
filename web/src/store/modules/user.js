import {getToken, setToken} from '@/utils/auth';
import {getInfo} from '@/api/user';

const user = {
  state: {
    token: getToken(),
    id: 0,
  },

  mutations: {
    SET_TOKEN: (state, token) => {
      state.token = token;
    },
    SET_ID: (state, id) => {
      state.id = id;
    },
  },

  actions: {
    // 登录
    Login({commit}) {
      setToken(1);
      commit('SET_TOKEN', 1);
    },
    // 获取用户信息
    GetInfo({commit, state}) {
      return new Promise((resolve, reject) => {
        getInfo(state.token).then((response) => {
          const responseData = response.data;
          const userInfo = responseData.data.user_info;
          commit('SET_ID', userInfo.id);
          resolve(responseData.data);
        }).catch((error) => {
          reject(error);
        });
      });
    },
  },

};

export default user;
