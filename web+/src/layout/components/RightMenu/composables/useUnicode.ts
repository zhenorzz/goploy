import { reactive, computed } from 'vue'

export default function useUnicode() {
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
}
