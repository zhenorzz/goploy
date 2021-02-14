// mixin.js
import { getNamespace } from '@/utils/namespace'

const mixin = {
  methods: {
    hasAdminPermission() {
      return [this.$global.Admin].indexOf(getNamespace()['role']) !== -1
    },
    hasManagerPermission() {
      return [this.$global.Admin, this.$global.Manager].indexOf(getNamespace()['role']) !== -1
    },
    hasGroupManagerPermission() {
      return [this.$global.Admin, this.$global.Manager, this.$global.GroupManager].indexOf(getNamespace()['role']) !== -1
    },
    isMember() {
      return this.$global.member === getNamespace()['role']
    },
    copy(content, message = '复制成功') {
      const el = document.createElement('textarea')
      el.value = content
      el.setAttribute('readonly', '')
      el.style.position = 'absolute'
      el.style.left = '-9999px'
      document.body.appendChild(el)
      const selected =
        document.getSelection().rangeCount > 0
          ? document.getSelection().getRangeAt(0)
          : false
      el.select()
      document.execCommand('copy')
      document.body.removeChild(el)
      if (selected) {
        document.getSelection().removeAllRanges()
        document.getSelection().addRange(selected)
      }
      this.$message.success(message)
    }
  }
}
export default mixin
