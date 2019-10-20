// mixin.js
const mixin = {
  methods: {
    hasAdminPermission() {
      return this.$store.getters.role.indexOf(['admin']) !== -1
    },
    hasManagerPermission() {
      return this.$store.getters.role.indexOf(['admin', 'manager']) !== -1
    },
    hasGroupManagerPermission() {
      return this.$store.getters.role.indexOf(['admin', 'manager', 'group-manager']) !== -1
    }
  }
}
export default mixin
