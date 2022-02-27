<template>
  <el-config-provider :locale="locale[$i18n.locale]">
    <router-view />
  </el-config-provider>
</template>
<script lang="ts">
import { ElConfigProvider } from 'element-plus'
import { elementLocale } from '@/lang'
import { defineComponent } from 'vue'
export default defineComponent({
  name: 'App',
  components: {
    ElConfigProvider,
  },
  data: function () {
    return {
      locale: elementLocale,
    }
  },
})
</script>
