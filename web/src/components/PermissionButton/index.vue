<template>
  <el-button v-if="show">
    <slot />
  </el-button>
</template>
<script lang="ts">
export default { name: 'PerButton' }
</script>
<script lang="ts" setup>
import { useStore } from 'vuex'
import { ref } from 'vue'
const props = defineProps({
  permissions: {
    type: Array,
    default: () => [],
  },
})
const store = useStore()

const show = ref(false)

show.value = props.permissions.some((permission) =>
  store.state.permission.permissionIds.includes(permission)
)
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
