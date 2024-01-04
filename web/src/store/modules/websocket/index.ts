import { Module, MutationTree, ActionTree } from 'vuex'
import { WebsocketState } from './types'
import { RootState } from '../../types'
import { parseTime } from '@/utils'
import { NamespaceKey, getNamespaceId } from '@/utils/namespace'

const state: WebsocketState = {
  ws: null,
  message: {},
  againConnectTime: 0, // 规定时间重连
}

const mutations: MutationTree<WebsocketState> = {
  SET_MESSAGE: (state, message) => {
    state.message = message
  },
  SET_WS: (state, ws) => {
    state.ws = ws
  },
  CLOSE_WS: (state) => {
    if (state.ws) {
      state.ws.close()
    }
  },
  SET_AGAINCONNECTTIME: (state, time) => {
    state.againConnectTime = time
  },
}

const actions: ActionTree<WebsocketState, RootState> = {
  init({ dispatch, commit, state }) {
    return new Promise((resolve) => {
      const websocket = new WebSocket(
        `${location.protocol.replace('http', 'ws')}//${location.host}${
          import.meta.env.VITE_APP_BASE_API
        }/ws/connect?${NamespaceKey}=${getNamespaceId()}`
      )
      websocket.onopen = () => {
        console.log('websocket连接成功, 时间：' + parseTime(Date.now()))
        // 连接成功，当成重连次数0 置0
        dispatch('setAgainConnectTime', 0)
        resolve(websocket)
      }

      websocket.onerror = () => {
        console.log('websocket连接发生错误, 时间：' + parseTime(Date.now()))
      }

      websocket.onmessage = (event) => {
        const responseData = JSON.parse(event.data)
        import.meta.env.DEV && console.log(responseData)
        dispatch('setMessage', responseData)
      }
      websocket.onclose = (e) => {
        // 1005 主动断开
        // websocket close code https://developer.mozilla.org/en-US/docs/Web/API/CloseEvent
        // 顶号，后台发送关闭帧
        dispatch('destory')

        if (e.code !== 1005) {
          if (state.againConnectTime === 0) {
            // 第一次连接断开进来 有一次重连机会
            dispatch('setAgainConnectTime', new Date().getTime())
            setTimeout(() => {
              console.log(
                '首次断开，重新主动连接, 时间：' + parseTime(Date.now())
              )
              dispatch('init')
            }, 60000)
          } else {
            if (new Date().getTime() - state.againConnectTime >= 60000) {
              console.log(
                '主动连接失败，再次尝试, 时间：' + parseTime(Date.now())
              )
              // 一分钟后的连接 一次重连机会已用完 还是连接失败,就弹窗询问用户
              ElMessageBox.confirm(
                'Detected Websocket disconnect, please reconnect!',
                'Tips',
                {
                  confirmButtonText: 'Confirm',
                  cancelButtonText: 'Cancel',
                  type: 'warning',
                }
              )
                .then(() => {
                  dispatch('init')
                })
                .catch(() => {
                  ElMessage({
                    type: 'info',
                    message: 'Cancel reconnect',
                  })
                })
            }
          }
        }

        console.log(
          `connection closed (${e.code})${e.reason}, 时间：${parseTime(
            Date.now()
          )}`
        )
      }
      commit('SET_WS', websocket)
    })
  },

  close({ commit }) {
    return new Promise(() => {
      commit('CLOSE_WS')
      commit('SET_WS', null)
    })
  },

  destory({ commit }) {
    commit('SET_WS', null)
  },

  setAgainConnectTime({ commit }, times) {
    commit('SET_AGAINCONNECTTIME', times)
  },

  setMessage({ commit }, message) {
    commit('SET_MESSAGE', message)
  },
}

export default <Module<WebsocketState, RootState>>{
  namespaced: true,
  state,
  mutations,
  actions,
}
