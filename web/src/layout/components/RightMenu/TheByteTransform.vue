<template>
  <el-row v-show="modelValue === 'byte'" type="flex" align="middle">
    <span style="width: 40px; font-size: 14px; margin-right: 10px"> Byte </span>
    <el-input
      v-model="byte.input"
      style="width: 130px"
      @keyup.enter="bytesToHumanSize"
    />
    <el-select v-model="byte.unit" style="width: 70px">
      <el-option :value="1" label="B" />
      <el-option :value="1 * 1024" label="KB" />
      <el-option :value="1024 * 1024" label="MB" />
    </el-select>
    <el-button type="primary" @click="bytesToHumanSize">>></el-button>
    <el-input v-model="byte.human" style="width: 200px" />
  </el-row>
</template>

<script lang="ts">
import { defineComponent, reactive } from 'vue'
import { humanSize } from '@/utils'
export default defineComponent({
  props: {
    modelValue: {
      type: String,
      default: '',
    },
  },
  setup() {
    const byte = reactive({
      input: '',
      unit: 1,
      human: '',
    })

    const bytesToHumanSize = () => {
      byte.human = humanSize(Number(byte.input) * byte.unit)
    }
    return {
      byte,
      bytesToHumanSize,
    }
  },
})
</script>
