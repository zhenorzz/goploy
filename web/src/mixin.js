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
    }
  }
}
export default mixin
