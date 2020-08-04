import { login, getInfo } from '@/api/user'
import { setLogin, logout } from '@/utils/auth'
import { getNamespace, setNamespace, setNamespaceList } from '@/utils/namespace'
import { resetRouter } from '@/router'

const state = {
  id: 0,
  account: '',
  name: '',
  superManager: 0
}

const mutations = {
  SET_ID: (state, id) => {
    state.id = id
  },
  SET_ACCOUNT: (state, account) => {
    state.account = account
  },
  SET_NAME: (state, name) => {
    state.name = name
  },
  SET_SUPER_MANAGER: (state, superManager) => {
    state.superManager = superManager
  }
}

const actions = {
  // user login
  login({ commit }, userInfo) {
    const { account, password } = userInfo
    return new Promise((resolve, reject) => {
      login({ account: account.trim(), password: password }).then(response => {
        const { data } = response

        let namespace = getNamespace()
        if (!namespace) {
          namespace = data.namespaceList[data.namespaceList.length - 1]
          setNamespace(namespace)
        }
        setNamespaceList(data.namespaceList)

        setLogin('ok')
        resolve()
      }).catch(error => {
        reject(error)
      })
    })
  },

  // get user info
  getInfo({ commit, state }) {
    return new Promise((resolve, reject) => {
      getInfo(state.token).then(response => {
        const { data } = response
        if (!data) {
          reject('Verification failed, please Login again.')
        }
        const { id, account, name, superManager } = data.userInfo
        commit('SET_ID', id)
        commit('SET_ACCOUNT', account)
        commit('SET_NAME', name)
        commit('SET_SUPER_MANAGER', superManager)
        resolve(data)
      }).catch(error => {
        reject(error)
      })
    })
  },

  // user logout
  logout({ commit }) {
    return new Promise((resolve) => {
      commit('SET_ID', 0)
      logout()
      resetRouter()
      resolve()
    })
  }

}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}

