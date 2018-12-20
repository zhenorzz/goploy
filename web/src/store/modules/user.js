import {login, logout} from '@/api/login';
import {getToken, setToken, removeToken} from '@/utils/auth';
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
    Login({commit}, userInfo) {
      const account = userInfo.account.trim();
      return new Promise((resolve, reject) => {
        login(account, userInfo.password).then((response) => {
          const data = response.data;
          setToken(data.token);
          commit('SET_TOKEN', data.token);
          resolve();
        }).catch((error) => {
          reject(error);
        });
      });
    },
    // 登出
    LogOut({commit, state}) {
      return new Promise((resolve, reject) => {
        logout(state.token).then(() => {
          commit('SET_TOKEN', '');
          removeToken();
          resolve();
        }).catch((error) => {
          reject(error);
        });
      });
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
