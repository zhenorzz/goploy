<template>
  <el-row
    v-show="modelValue === 'byte'"
    type="flex"
    align="middle"
    style="width: 100%"
  >
    <el-input
      v-model="byte.input"
      style="width: 250px"
      @keyup.enter="bytesToHumanSize"
    >
      <template #prepend>Byte</template>
      <template #append>
        <el-select v-model="byte.unit" style="width: 70px">
          <el-option :value="1" label="B" />
          <el-option :value="1 * 1024" label="KB" />
          <el-option :value="1024 * 1024" label="MB" />
        </el-select>
      </template>
    </el-input>
    <el-button type="primary" @click="bytesToHumanSize">>></el-button>
    <el-input v-model="byte.human" style="flex: 1" />
  </el-row>
</template>

<script lang="ts" setup>
import { reactive } from 'vue'
import { humanSize } from '@/utils'
defineProps({
  modelValue: {
    type: String,
    default: '',
  },
})
const byte = reactive({
  input: '',
  unit: 1,
  human: '',
})

const bytesToHumanSize = () => {
  byte.human = humanSize(Number(byte.input) * byte.unit)
}
</script>
