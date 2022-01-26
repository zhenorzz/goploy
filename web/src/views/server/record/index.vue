<template>
  <div ref="record"></div>
</template>
<script lang="ts">
import 'asciinema-player/dist/bundle/asciinema-player.css'
import * as AsciinemaPlayer from 'asciinema-player'
import { NamespaceKey, getNamespaceId } from '@/utils/namespace'
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'ServerTerminal',
  data() {
    return {}
  },

  created() {},
  mounted() {
    const castUrl = `${location.origin}${
      import.meta.env.VITE_APP_BASE_API
    }/server/getTerminalRecord?${NamespaceKey}=${getNamespaceId()}`
    AsciinemaPlayer.create(castUrl, this.$refs['record'], {
      fit: false,
      fontSize: '14px',
    })
  },
  methods: {},
})
</script>
<style lang="scss" scoped>
@import '@/styles/mixin.scss';
.template-dialog {
  padding-right: 10px;
  height: 400px;
  overflow-y: auto;
  @include scrollBar();
}
</style>
