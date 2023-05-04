<template>
  <section class="app-main">
    <router-view v-slot="{ Component }">
      <transition name="fade-transform" mode="out-in">
        <keep-alive :include="cachedViews">
          <component :is="Component" :key="key" />
        </keep-alive>
      </transition>
    </router-view>
  </section>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
export default defineComponent({
  name: 'AppMain',
  computed: {
    cachedViews() {
      return import.meta.env.PROD === true
        ? this.$store.state.tagsView.cachedViews.join(',').split(',')
        : ['ServerAgent']
    },
    key() {
      return this.$route.path
    },
  },
})
</script>

<style scoped>
.app-main {
  /*50 = navbar  */
  min-height: calc(100vh - 97px);
  width: 100%;
  position: relative;
  overflow: hidden;
}
.fixed-header + .app-main {
  padding-top: 50px;
}
</style>

<style lang="scss">
// fix css style bug in open el-dialog
.el-popup-parent--hidden {
  .fixed-header {
    padding-right: 15px;
  }
}
</style>
