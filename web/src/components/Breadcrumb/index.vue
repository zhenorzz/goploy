<template>
  <el-breadcrumb class="app-breadcrumb" separator="/">
    <transition-group name="breadcrumb" tag="span">
      <el-breadcrumb-item v-for="(item, index) in levelList" :key="item.path">
        <span
          v-if="item.redirect === 'noRedirect' || index == levelList.length - 1"
          class="no-redirect"
          @click.prevent="handleLink(item)"
        >
          {{ $t(`route.${item.meta.title}`) }}
        </span>
        <a v-else @click.prevent="handleLink(item)">
          {{ $t(`route.${item.meta.title}`) }}
        </a>
      </el-breadcrumb-item>
    </transition-group>
  </el-breadcrumb>
</template>

<script lang="ts" setup>
import { compile } from 'path-to-regexp'
import { ref, watch } from 'vue'
import {
  RouteRecordRaw,
  RouteLocationRaw,
  useRoute,
  useRouter,
} from 'vue-router'

const levelList = ref()
const route = useRoute()
const router = useRouter()
watch(route, () => {
  getBreadcrumb()
})
getBreadcrumb()
function getBreadcrumb() {
  // only show routes with meta.title
  const matched = route.matched.filter((item) => item.meta && item.meta.title)
  levelList.value = matched.filter(
    (item) => item.meta && item.meta.title && item.meta.breadcrumb !== false
  )
}

function pathCompile(path: string) {
  // To solve this problem https://github.com/PanJiaChen/vue-element-admin/issues/561
  const { params } = route
  var toPath = compile(path)
  return toPath(params)
}

function handleLink(item: RouteRecordRaw) {
  const { redirect, path } = item
  if (redirect) {
    router.push(<RouteLocationRaw>redirect)
    return
  }
  router.push(pathCompile(path))
}
</script>

<style lang="scss" scoped>
.app-breadcrumb.el-breadcrumb {
  display: inline-block;
  font-size: 14px;
  line-height: 50px;
  margin-left: 8px;

  .no-redirect {
    color: #97a8be;
    cursor: text;
  }
  .breadcrumb-leave-active {
    position: absolute;
  }

  .breadcrumb-enter-active,
  .breadcrumb-leave-active {
    transition: all 0.5s;
  }

  .breadcrumb-leave-active {
    opacity: 0;
    transform: translateX(20px);
  }

  .breadcrumb-enter-active {
    opacity: 0;
    transform: translateX(20px);
  }

  .breadcrumb-enter-to {
    opacity: 1;
    transform: translateX(0px);
  }

  .breadcrumb-move {
    transition: all 0.5s;
  }
}
</style>
