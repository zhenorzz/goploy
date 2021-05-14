import { ModuleTree, createStore, createLogger } from 'vuex'
import { RootState } from './types'
const files = import.meta.globEager('./modules/*/index.ts')
const modules: ModuleTree<RootState> = {}
for (const path in files) {
  modules[path.split('/')[2]] = files[path].default
}
const store = createStore({
  strict: true,
  modules: { ...modules },
  plugins: process.env.NODE_ENV !== 'production' ? [createLogger()] : [],
})
export default store
