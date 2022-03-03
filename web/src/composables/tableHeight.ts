import { onMounted, ref } from 'vue'
export default () => {
  const tableHeight = ref(window.innerHeight - 160)
  const adjustTableHeight = () => {
    try {
      const outHeight =
        document.getElementsByClassName('app-bar')[0].clientHeight + 180 // 除了table外 查询与按钮的高度 + 其他的高度150
      const maxHeight = window.innerHeight - outHeight
      tableHeight.value = maxHeight
    } catch (e) {
      console.log('缺少节点app-bar')
    }
  }
  onMounted(() => {
    adjustTableHeight()
    window.addEventListener('resize', adjustTableHeight)
  })
  return { tableHeight, adjustTableHeight }
}
