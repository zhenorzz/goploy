<template>
  <a v-if="isExternal(to)" :href="to" target="_blank" rel="noopener">
    <slot />
  </a>
  <router-link v-else :to="to">
    <slot />
  </router-link>
</template>

<script lang="ts">
import { isExternal } from '@/utils/validate'
import { defineComponent } from 'vue'
export default defineComponent({
  props: {
    to: {
      type: String,
      required: true,
    },
  },
  methods: {
    isExternal,
  },
})
</script>
