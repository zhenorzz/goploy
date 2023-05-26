<template>
  <el-link
    v-if="href"
    type="primary"
    :underline="false"
    :href="href"
    target="_blank"
  >
    {{ text }}
  </el-link>
  <span v-else>
    {{ text }}
  </span>
</template>
<script lang="ts">
export default { name: 'RepoURL' }
</script>
<script lang="ts" setup>
import { repoURLToHref } from '@/utils'
import { computed } from 'vue'

const props = defineProps({
  url: {
    type: String,
    default: '',
  },
  suffix: {
    type: String,
    default: '',
  },
  text: {
    type: String,
    default: '',
  },
})

const href = computed(() => {
  return repoURLToHref(props.url, props.suffix)
})
</script>

<style scoped>
.svg-icon {
  width: 1em;
  height: 1em;
  vertical-align: -0.15em;
  fill: currentColor;
  overflow: hidden;
}

.svg-external-icon {
  background-color: currentColor;
  mask-size: cover !important;
  display: inline-block;
}
</style>
