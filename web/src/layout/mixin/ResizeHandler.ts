import { computed, watch, onBeforeMount, onBeforeUnmount, onMounted } from 'vue'
import { useStore } from 'vuex'
import { useRoute } from 'vue-router'
const { body } = document
const WIDTH = 992 // refer to Bootstrap's responsive design
export default () => {
  const store = useStore()
  const route = useRoute()
  const app = computed(() => store.state['app'])
  watch(route, () => {
    if (app.value.device === 'mobile' && app.value.sidebar.opened) {
      store.dispatch('app/closeSideBar', { withoutAnimation: false })
    }
  })
  onBeforeMount(() => window.addEventListener('resize', $_resizeHandler))
  onMounted(() => {
    const isMobile = $_isMobile()
    if (isMobile) {
      store.dispatch('app/toggleDevice', 'mobile')
      store.dispatch('app/closeSideBar', { withoutAnimation: true })
      import.meta.env.DEV === true &&
        import('vconsole').then((module) => {
          new module.default()
        })
    }
  })
  onBeforeUnmount(() => window.addEventListener('resize', $_resizeHandler))

  function $_isMobile() {
    const rect = body.getBoundingClientRect()
    return rect.width - 1 < WIDTH
  }

  function $_resizeHandler() {
    if (!document.hidden) {
      const isMobile = $_isMobile()
      store.dispatch('app/toggleDevice', isMobile ? 'mobile' : 'desktop')

      if (isMobile) {
        store.dispatch('app/closeSideBar', { withoutAnimation: true })
      }
    }
  }
}
