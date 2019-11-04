// mixin.js
const mixin = {
  methods: {
    hasAdminPermission() {
      return ['admin'].indexOf(this.$store.getters.role) !== -1
    },
    hasManagerPermission() {
      return ['admin', 'manager'].indexOf(this.$store.getters.role) !== -1
    },
    hasGroupManagerPermission() {
      return ['admin', 'manager', 'group-manager'].indexOf(this.$store.getters.role) !== -1
    }
  }
}
export default mixin
