<template>
  <el-row v-show="modelValue === 'color'">
    <el-row type="flex" align="middle">
      <span style="width: 40px; font-size: 14px; margin-right: 10px">
        HEX
      </span>
      <el-input
        v-model="cHexExchange.hex"
        style="width: 200px"
        placeholder="#FFFFFF"
        clearable
        @keyup.enter="hexToRGB"
      />
      <el-button type="primary" @click="hexToRGB">>></el-button>
      <el-input v-model="cHexExchange.rgb" style="width: 200px" />
    </el-row>
    <el-row style="margin-top: 10px" type="flex" align="middle">
      <span style="width: 40px; font-size: 14px; margin-right: 10px">
        RGB
      </span>
      <el-input
        v-model="rgbExchange.rgb"
        style="width: 200px"
        placeholder="(255,255,255)"
        clearable
        @keyup.enter="rgbToHex"
      />
      <el-button type="primary" @click="rgbToHex">>></el-button>
      <el-input v-model="rgbExchange.hex" style="width: 200px" />
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
