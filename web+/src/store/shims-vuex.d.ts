import { Store } from 'vuex'
import { RootState } from './types'
declare module '@vue/runtime-core' {
  // Declare your own store states.
  interface ComponentCustomProperties {
    $store: Store<RootState>
  }
}
