import { Module, MutationTree, ActionTree } from 'vuex'
import { RootState } from '../../types'
import Cookies from 'js-cookie'
import { getLanguage } from '@/lang/index'
import { AppState } from './types'

const state: AppState = {
  sidebar: {
    opened: Cookies.get('sidebarStatus')
      ? !!Cookies.get('sidebarStatus')
      : true,
    withoutAnimation: false,
  },
  device: 'desktop',
  language: getLanguage(),
}

const mutations: MutationTree<AppState> = {
  TOGGLE_SIDEBAR: (state) => {
    state.sidebar.opened = !state.sidebar.opened
    state.sidebar.withoutAnimation = false
    if (state.sidebar.opened) {
      Cookies.set('sidebarStatus', '1')
    } else {
      Cookies.set('sidebarStatus', '0')
    }
  },
  CLOSE_SIDEBAR: (state, withoutAnimation: boolean) => {
    Cookies.set('sidebarStatus', '0')
    state.sidebar.opened = false
    state.sidebar.withoutAnimation = withoutAnimation
  },
  TOGGLE_DEVICE: (state, device: string) => {
    state.device = device
  },
  SET_LANGUAGE: (state, language: string) => {
    state.language = language
    Cookies.set('language', language)
  },
}

const actions: ActionTree<AppState, RootState> = {
  toggleSideBar(context) {
    context.commit('TOGGLE_SIDEBAR')
  },
  closeSideBar({ commit }, { withoutAnimation }) {
    commit('CLOSE_SIDEBAR', withoutAnimation)
  },
  toggleDevice({ commit }, device: string) {
    commit('TOGGLE_DEVICE', device)
  },
  setLanguage({ commit }, language: string) {
    commit('SET_LANGUAGE', language)
  },
}

export default <Module<AppState, RootState>>{
  namespaced: true,
  state,
  mutations,
  actions,
}
