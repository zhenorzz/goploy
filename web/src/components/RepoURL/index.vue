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
  let url = props.url.trim()
  if (props.url === '') {
    return ''
  }

  url = props.url.split(' ').shift() || ''
  const lastDotGitIndex = url.lastIndexOf('.git')
  if (lastDotGitIndex !== -1) {
    url = url.substring(0, lastDotGitIndex)
  }

  if (url.startsWith('git@')) {
    return '//' + url.substring(4).replace(':', '/') + props.suffix
  } else if (url.startsWith('http')) {
    const urlObj = new URL(url)
    if (urlObj.password !== '') {
      urlObj.password = '******'
    }
    return urlObj.toString() + props.suffix
  } else if (url.startsWith('svn')) {
    return url
  } else {
    return ''
  }
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
