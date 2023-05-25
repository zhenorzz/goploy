<template>
  <el-row class="app-container">
    <el-row class="tabs-container">
      <el-scrollbar style="width: 100%">
        <el-row class="tabs" justify="start" align="middle">
          <div
            v-for="(item, index) in serverList"
            :key="item.uuid"
            class="tabs-item"
            :class="item.uuid === currentUUID ? 'tabs-item-selected' : ''"
          >
            <el-row>
              <div class="tabs-item-serial" @click="selectTab(item)">
                {{ index + 1 }}
              </div>
              <div
                class="tabs-item-name"
                :title="`${item.server.name} : ${item.dir}`"
                @click="selectTab(item)"
              >
                {{ item.server.name }} : {{ item.dir }}
              </div>
              <el-button
                link
                style="font-size: 14px; padding-left: 8px"
                @click="deleteTab(item, index)"
              >
                x
              </el-button>
            </el-row>
          </div>

          <div class="tabs-plus">
            <el-select
              v-model="serverId"
              style="width: 200px"
              filterable
              clearable
              @change="selectServer"
            >
              <el-option
                v-for="item in serverOption"
                :key="item.id"
                :label="item.label"
                :value="item.id"
              />
            </el-select>
          </div>
          <div class="tabs-placeholder"></div>
        </el-row>
      </el-scrollbar>
    </el-row>
    <explorer
      v-for="item in serverList"
      v-show="item.uuid === currentUUID"
      :key="item.uuid"
      :server="item.server"
      @dir-change="handleDirChange"
    ></explorer>
  </el-row>
</template>

<script lang="ts">
export default { name: 'ServerNginx' }
</script>
<script lang="ts" setup>
import explorer from './explorer.vue'
import { ServerOption, ServerData } from '@/api/server'
import { ref } from 'vue'
interface nginx {
  uuid: number
  server: ServerData
  dir: string
}
const currentUUID = ref(0)
const serverOption = ref<ServerOption['datagram']['list']>([])
const serverId = ref('')
const serverList = ref<nginx[]>([])
const selectedNginx = ref<nginx>({} as nginx)

function handleDirChange(dir: string) {
  selectedNginx.value.dir = dir
}

getServerOption()

function getServerOption() {
  new ServerOption().request().then((response) => {
    serverOption.value = response.data.list
  })
}

function selectTab(nginx: nginx) {
  currentUUID.value = nginx.uuid
  selectedNginx.value = nginx
}

function deleteTab(nginx: nginx, index: number) {
  serverList.value.splice(index, 1)
  if (currentUUID.value === nginx.uuid) {
    currentUUID.value =
      serverList.value.length === 0
        ? 0
        : serverList.value[serverList.value.length - 1].uuid
  }
}

function selectServer(value: number) {
  const server =
    serverOption.value.find((_) => _.id === value) || ({} as ServerData)
  if (serverList.value.length === 0) {
    currentUUID.value = 0
  } else {
    currentUUID.value = serverList.value[serverList.value.length - 1].uuid + 1
  }
  const serverTab = { uuid: currentUUID.value, server, dir: '' }
  serverList.value.push(serverTab)
  selectTab(serverTab)
  serverId.value = ''
}
</script>
<style lang="scss" scoped>
@import '@/styles/mixin.scss';
.app-container {
  background-color: var(--el-bg-color);
  border: 1px solid var(--el-border-color);
}
.tabs-container {
  width: 100%;
  height: 46px;
  .tabs {
    height: auto;
    flex-wrap: nowrap;
    &-item {
      font-size: 14px;
      flex-shrink: 0;
      width: 195px;
      height: 46px;
      line-height: 46px;
      padding: 0 10px;
      border-bottom: 1px solid var(--el-border-color);
      background-color: var(--el-disabled-bg-color);
      &:not(:last-child) {
        border-right: 1px solid var(--el-border-color);
      }
      &-serial {
        cursor: pointer;
        color: var(--el-text-color-secondary);
        text-align: center;
        font-weight: 600;
        width: 20px;
      }
      &-name {
        flex: 1;
        cursor: pointer;
        text-overflow: ellipsis;
        overflow: hidden;
        padding-left: 5px;
        white-space: nowrap;
        display: inline-block;
      }
      &-selected {
        border-bottom: none;
        background-color: var(--el-bg-color);
      }
    }
    &-placeholder {
      height: 46px;
      flex: 1;
      border-bottom: 1px solid var(--el-border-color);
    }
    &-plus {
      text-align: center;
      border-bottom: 1px solid var(--el-border-color);
      border-right: 1px solid var(--el-border-color);
    }
  }
}
</style>
<style lang="scss">
.tabs-plus {
  .el-input__wrapper {
    border-radius: 0px !important;
    box-shadow: none;
  }
  .el-input__inner {
    height: 43px !important;
  }
}
</style>
