<template>
  <el-scrollbar wrap-class="scrollbar-wrapper">
    <el-menu
      :default-active="activeMenu"
      :background-color="variables.menuBg"
      :text-color="variables.menuText"
      :unique-opened="false"
      :active-text-color="variables.menuActiveText"
      :collapse-transition="true"
      mode="horizontal"
    >
      <sidebar-item
        v-for="item in permission.routes"
        :key="item.path"
        :item="item"
        :base-path="item.path"
      />
    </el-menu>
  </el-scrollbar>
</template>

<script lang="ts" setup>
import { useStore } from 'vuex'
import { useRoute } from 'vue-router'
import SidebarItem from './SidebarItem.vue'
import variables from '@/styles/variables.module.scss'
import { computed } from 'vue'
const route = useRoute()
const store = useStore()
const permission = computed(() => store.state['permission'])
const activeMenu = computed(() => {
  const { meta, path } = route
  // if set path, the sidebar will highlight the path you set
  if (meta.activeMenu) {
    return <string>meta.activeMenu
  }
  return path
})
</script>
