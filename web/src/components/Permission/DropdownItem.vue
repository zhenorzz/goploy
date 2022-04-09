<template>
  <el-dropdown-item v-if="show"><slot /></el-dropdown-item>
</template>
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
