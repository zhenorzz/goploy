<template>
  <el-row v-show="modelValue === 'color'">
    <el-row type="flex" align="middle" style="width: 100%">
      <el-input
        v-model="cHexExchange.hex"
        style="flex: 1"
        placeholder="#FFFFFF"
        clearable
        @keyup.enter="hexToRGB"
      >
        <template #prepend>HEX</template>
      </el-input>
      <el-button type="primary" @click="hexToRGB">>></el-button>
      <el-input v-model="cHexExchange.rgb" style="flex: 1" readonly>
        <template #prepend>RGB</template>
      </el-input>
    </el-row>
    <el-row style="width: 100%; margin-top: 10px" type="flex" align="middle">
      <el-input
        v-model="rgbExchange.rgb"
        style="flex: 1"
        placeholder="(255,255,255)"
        clearable
        @keyup.enter="rgbToHex"
      >
        <template #prepend>RGB</template>
      </el-input>
      <el-button type="primary" @click="rgbToHex">>></el-button>
      <el-input v-model="rgbExchange.hex" style="flex: 1" readonly>
        <template #prepend>HEX</template>
      </el-input>
    </el-row>
  </el-row>
</template>

<script lang="ts" setup>
import { reactive } from 'vue'
defineProps({
  modelValue: {
    type: String,
    default: '',
  },
})
const cHexExchange = reactive({
  hex: '',
  rgb: '',
})

const rgbExchange = reactive({
  hex: '',
  rgb: '',
})

const hexToRGB = () => {
  const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(
    cHexExchange.hex
  )

  if (result) {
    const r = parseInt(result[1], 16)
    const g = parseInt(result[2], 16)
    const b = parseInt(result[3], 16)
    cHexExchange.rgb = 'rgb(' + r + ', ' + g + ', ' + b + ')'
  } else {
    cHexExchange.rgb = 'rgb(0, 0, 0)'
  }
}
const rgbToHex = () => {
  const color = rgbExchange.rgb.replace(/\(|\)|rgb/g, '')
  const rgb = color.split(',')
  const r = parseInt(rgb[0])
  const g = parseInt(rgb[1])
  const b = parseInt(rgb[2])
  const hex = '#' + ((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1)
  rgbExchange.hex = hex.toLocaleUpperCase()
}
</script>
