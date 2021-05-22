<template>
  <section class="app-main">
    <router-view v-slot="{ Component }" :key="key">
      <transition name="fade-transform" mode="out-in" tag="div">
        <keep-alive :include="cachedViews">
          <component :is="Component" />
        </keep-alive>
      </transition>
    </router-view>
  </section>
</template>

<script>
import { defineComponent } from 'vue'
export default defineComponent({
  name: 'AppMain',
  computed: {
    cachedViews() {
      return import.meta.env.NODE_ENV === 'production'
        ? this.$store.state.tagsView.cachedViews
        : []
      // return this.$store.state.tagsView.cachedViews
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
  min-height: calc(100vh - 100px);
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
