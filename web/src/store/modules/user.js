import {login} from '@/api/login';
import {getToken, setToken, removeToken} from '@/utils/auth';
import {getInfo} from '@/api/user';

const user = {
  state: {
    token: getToken(),
    id: 0,
    name: '',
    account: '',
    role: '',
  },

  mutations: {
    SET_TOKEN: (state, token) => {
      state.token = token;
    },
    SET_ID: (state, id) => {
      state.id = id;
    },
    SET_ACCOUNT: (state, account) => {
      state.account = account;
    },
    SET_NAME: (state, name) => {
      state.name = name;
    },
    SET_ROLE: (state, role) => {
      state.role = role;
    },
  },

  actions: {
    // 登录
    Login({commit}, userInfo) {
      const account = userInfo.account.trim();
      return new Promise((resolve, reject) => {
        login(account, userInfo.password).then((response) => {
          const data = response.data;
          const token = data.data.token;
          setToken(token);
          commit('SET_TOKEN', token);
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
          const userInfo = responseData.data.userInfo;
          commit('SET_ID', userInfo.id);
          commit('SET_NAME', userInfo.name);
          commit('SET_ROLE', userInfo.role);
          commit('SET_ACCOUNT', userInfo.account);
          resolve(responseData.data);
        }).catch((error) => {
          reject(error);
        });
      });
    },
    // 前端 登出
    FedLogOut({commit}) {
      return new Promise((resolve) => {
        commit('SET_TOKEN', '');
        removeToken();
        resolve();
      });
    },
  },

};

export default user;
