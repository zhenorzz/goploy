import { onBeforeMount, onBeforeUnmount, onMounted } from 'vue'
import { useStore } from 'vuex'
const { body } = document
const WIDTH = 992 // refer to Bootstrap's responsive design
export default () => {
  const store = useStore()

  onBeforeMount(() => window.addEventListener('resize', $_resizeHandler))
  onMounted(() => {
    const isMobile = $_isMobile()
    if (isMobile) {
      store.dispatch('app/toggleDevice', 'mobile')
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
    }
  }
}
