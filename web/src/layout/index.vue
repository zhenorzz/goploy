<template>
  <div :class="classObj" class="app-wrapper">
    <div>
      <Navbar />
    </div>
    <div class="main-container">
      <TagsView />
      <AppMain />
    </div>
  </div>
</template>
<script lang="ts">
export default { name: 'Layout' }
</script>
<script lang="ts" setup>
import { Navbar, AppMain, TagsView } from './components'
import ResizeHandler from './mixin/ResizeHandler'
import { computed } from 'vue'
import { useStore } from 'vuex'
const store = useStore()
const sidebar = computed(() => store.state['app'].sidebar)
const device = computed(() => store.state['app'].device)
const classObj = computed(() => {
  return {
    openSidebar: sidebar.value.opened,
    withoutAnimation: sidebar.value.withoutAnimation,
    mobile: device.value === 'mobile',
  }
})
ResizeHandler()
store.dispatch('websocket/init')
</script>

<style lang="scss" scoped>
@import '@/styles/mixin.scss';
@import '@/styles/variables.module.scss';

.app-wrapper {
  @include clearfix;
  position: relative;
  height: 100%;
  width: 100%;
  &.mobile.openSidebar {
    position: fixed;
    top: 0;
  }
}

.hideSidebar .fixed-header {
  width: calc(100% - 54px);
}

.mobile .fixed-header {
  width: 100%;
}
</style>
