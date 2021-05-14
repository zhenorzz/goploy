import { Module } from 'vuex'
import { SettingState } from './types'
import { RootState } from '../../types'

const state: SettingState = {
  fixedHeader: false,
}

const setting: Module<SettingState, RootState> = {
  namespaced: true,
  state,
}

export default setting
