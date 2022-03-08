import { Module } from 'vuex'
import { SettingState } from './types'
import { RootState } from '../../types'

const state: SettingState = {
  fixedHeader: false,
}

export default <Module<SettingState, RootState>>{
  namespaced: true,
  state,
}
