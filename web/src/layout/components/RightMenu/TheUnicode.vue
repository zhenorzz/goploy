<template>
  <el-row v-show="modelValue === 'unicode'" style="width: 100%">
    <el-input
      v-model="unicode.escape"
      type="textarea"
      :autosize="{ minRows: 2 }"
      placeholder="Please enter unescaped unicode encoding"
    />
    <el-input
      :value="unicodeUnescapeStr"
      style="margin-top: 10px"
      type="textarea"
      :autosize="{ minRows: 2 }"
      readonly
    />
  </el-row>
</template>

<script lang="ts" setup>
import { computed, reactive } from 'vue'
defineProps({
  modelValue: {
    type: String,
    default: '',
  },
})
const unicode = reactive({
  escape: '',
})
const unicodeUnescapeStr = computed(() =>
  unescape(unicode.escape.replace(/\\u/g, '%u'))
)
</script>
