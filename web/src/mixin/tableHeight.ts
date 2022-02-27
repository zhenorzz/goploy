export default {
  data() {
    return {
      tableHeight: window.innerHeight - 160, // table高度
    }
  },
  mounted() {
    try {
      const outHeight =
        document.getElementsByClassName('app-bar')[0].clientHeight + 140 // 除了table外 查询与按钮的高度 + 其他的高度150
      const maxHeight = window.innerHeight - outHeight
      this.tableHeight = maxHeight
    } catch (e) {
      console.log('缺少节点app-bar')
    }
  },
}
