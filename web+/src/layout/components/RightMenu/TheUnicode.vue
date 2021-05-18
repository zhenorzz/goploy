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

<script lang="ts">
import { computed, defineComponent, reactive } from 'vue'

export default defineComponent({
  props: {
    modelValue: {
      type: String,
      default: '',
    },
  },
  setup() {
    const unicode = reactive({
      escape: '',
    })
    const unicodeUnescapeStr = computed(() =>
      unescape(unicode.escape.replace(/\\u/g, '%u'))
    )
    return {
      unicode,
      unicodeUnescapeStr,
    }
  },
})
</script>
