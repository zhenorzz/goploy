<template>
  <el-dialog
    v-model="dialogVisible"
    title="WebSFTP"
    :close-on-click-modal="false"
    :fullscreen="$store.state.app.device === 'mobile'"
  >
    <el-select
      v-model="serverId"
      v-loading="serverLoading"
      filterable
      default-first-option
      placeholder="please select server"
      style="width: 100%; margin-bottom: 10px"
      @change="handleSelectServer"
    >
      <el-option
        v-for="item in serverOption"
        :key="item.id"
        :label="`${item.name}(${item.description})`"
        :value="item.id"
      />
    </el-select>
    <el-input
      v-model="dir"
      style="width: 100%"
      :disabled="!wsConnected"
      :readonly="!wsConnected"
      @keyup.enter="goto(dir)"
    >
      <template #prepend>Path</template>
      <template #append>
        <i v-if="!wsConnected" class="el-icon-loading"></i>
        <el-button
          v-if="wsConnected"
          type="primary"
          icon="el-icon-right"
          @click="goto(dir)"
        />
        <el-button
          v-if="wsConnected"
          type="primary"
          icon="el-icon-refresh-left"
          @click="back"
        />
        <el-button
          v-if="wsConnected"
          type="primary"
          icon="el-icon-refresh"
          @click="refresh"
        />
      </template>
    </el-input>
    <el-table
      v-loading="tableLoading"
      border
      stripe
      highlight-current-row
      :data="tableData"
      style="width: 100%"
      :max-height="460"
    >
      <el-table-column prop="name" :label="$t('name')" min-width="100">
        <template #default="scope">
          <i v-if="scope.row.isDir" class="el-icon-folder-opened"></i>
          {{ scope.row.name }}
        </template>
      </el-table-column>
      <el-table-column prop="size" :label="$t('size')" width="100">
        <template #default="scope">
          {{ humanSize(scope.row.size) }}
        </template>
      </el-table-column>
      <el-table-column prop="mode" label="mode" width="100" />
      <el-table-column
        prop="modTime"
        :label="$t('modifiedTime')"
        width="135"
        align="center"
      />
      <el-table-column
        prop="operation"
        :label="$t('op')"
        width="100"
        align="center"
        :fixed="$store.state.app.device === 'mobile' ? false : 'right'"
      >
        <template #default="scope">
          <el-button
            style="margin-right: 10px"
            :disabled="!scope.row.isDir"
            type="text"
            icon="el-icon-right"
            @click="goto(`${dir}/${scope.row.name}`)"
          />
          <el-link
            :disabled="scope.row.isDir"
            :href="downloadHref(`${dir}/${scope.row.name}`)"
            target="_blank"
            :underline="false"
          >
            <i class="el-icon-download"></i>
          </el-link>
        </template>
      </el-table-column>
    </el-table>
  </el-dialog>
</template>

<script lang="ts">
import path from 'path-browserify'
import { ElMessage } from 'element-plus'
import { computed, watch, defineComponent, ref } from 'vue'
import { ServerOption } from '@/api/server'
import { humanSize } from '@/utils'
export default defineComponent({
  props: {
    modelValue: {
      type: Boolean,
      default: false,
    },
  },
  emits: ['update:modelValue'],
  setup(props, { emit }) {
    const dialogVisible = computed({
      get: () => props.modelValue,
      set: (val) => {
        emit('update:modelValue', val)
      },
    })

    const serverLoading = ref(false)
    const serverOption = ref<ServerOption['datagram']['list']>([])
    watch(dialogVisible, (val: boolean) => {
      if (val === true) {
        serverLoading.value = true
        new ServerOption()
          .request()
          .then((response) => {
            serverOption.value = response.data.list
          })
          .finally(() => {
            serverLoading.value = false
          })
      }
    })
    const serverId = ref('')
    const tableData = ref([])
    const wsConnected = ref(false)
    let ws: WebSocket
    const handleSelectServer = () => {
      const selectedServer = serverOption.value.find(
        (_) => _.id === Number(serverId.value)
      )
      if (!selectedServer) {
        return
      }
      ws = new WebSocket(
        `${location.protocol.replace('http', 'ws')}//${
          window.location.host + import.meta.env.VITE_APP_BASE_API
        }/ws/sftp?serverId=${selectedServer.id}`
      )
      ws.onopen = (event) => {
        console.log(event)
        wsConnected.value = true
      }
      ws.onerror = (error) => {
        console.log(error)
      }
      ws.onclose = (event) => {
        console.log(event)
      }
      ws.onmessage = (event) => {
        const responseData = JSON.parse(event.data)
        console.log(responseData)
        tableLoading.value = false
        if (responseData['code'] > 0) {
          ElMessage.error(responseData['message'])
          tableData.value = []
        } else {
          tableData.value = responseData['data']
          if (dir.value !== '/') {
            tableData.value.unshift({
              name: '..',
              size: 0,
              mode: '',
              modTime: '',
              isDir: true,
            })
          }
        }
      }
    }
    const tableLoading = ref(false)
    const dir = ref('/')
    const dirHistory: string[] = []
    const goto = (target: string) => {
      tableLoading.value = true
      dir.value = path.normalize(target)
      dirHistory.push(dir.value)
      ws.send(target)
    }
    const back = () => {
      if (dirHistory.length > 1) {
        dirHistory.pop()
        dir.value = dirHistory[dirHistory.length - 1]
        ws.send(dir.value)
      }
    }
    const refresh = () => {
      tableLoading.value = true
      ws.send(dir.value)
    }
    const downloadHref = (file: string) => {
      file = path.normalize(file)
      return `${import.meta.env.VITE_APP_BASE_API}/server/remoteFile?id=${
        serverId.value
      }&file=${file}`
    }
    return {
      dialogVisible,
      serverLoading,
      serverOption,
      serverId,
      handleSelectServer,
      wsConnected,
      dir,
      tableLoading,
      tableData,
      humanSize,
      goto,
      back,
      refresh,
      downloadHref,
    }
  },
})
</script>
<style lang="scss" scoped>
.xterm {
  width: 100%;
  height: 100%;
}
</style>
