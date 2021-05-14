<script lang="ts">
import { defineComponent } from 'vue'
export default defineComponent({
  created() {
    const { params, query } = this.$route
    const { path } = params
    this.$router.replace({ path: '/' + path, query })
  },
  render: (h: CallableFunction) => {
    return h() // avoid warning message
  },
})
</script>
