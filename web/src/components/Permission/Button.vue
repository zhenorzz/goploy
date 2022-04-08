<template>
  <el-button v-if="show && content"><slot /></el-button>
  <el-button v-else-if="show"></el-button>
</template>
<script lang="ts" setup>
import { useStore } from 'vuex'
import { useSlots, ref, onMounted } from 'vue'
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

const slots = useSlots()
const content = ref(false)
onMounted(() => {
  if (slots.default && slots.default().length) {
    content.value = true
  }
})
</script>
