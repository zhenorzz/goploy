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
            <el-popover
              v-model:visible="serverOptionVisible"
              placement="bottom-start"
              :width="200"
              trigger="click"
              popper-class="tabs-popover"
              :show-arrow="false"
              :offset="-2"
            >
              <el-input
                v-model="serverFilterInput"
                placeholder="Filter server name"
                @input="filterServer"
              />
              <el-scrollbar class="server-list">
                <div
                  v-for="server in serverFilteredOption"
                  :key="server.id"
                  style="padding: 4px 0"
                >
                  <el-button link @click="selectServer(server)">
                    <span :title="server.label">
                      {{ server.label }}
                    </span>
                  </el-button>
                </div>
              </el-scrollbar>
              <template #reference>
                <el-button
                  link
                  style="height: 45px; width: 100%; font-size: 24px"
                  @click="serverOptionVisible = !serverOptionVisible"
                >
                  +
                </el-button>
              </template>
            </el-popover>
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
export default { name: 'ServerSFTP' }
</script>
<script lang="ts" setup>
import explorer from './explorer.vue'
import { ServerOption, ServerData } from '@/api/server'
import { ref } from 'vue'
interface sftp {
  uuid: number
  server: ServerData
  dir: string
}
const currentUUID = ref(0)
const serverOptionVisible = ref(false)
const serverOption = ref<ServerOption['datagram']['list']>([])
const serverFilteredOption = ref<ServerOption['datagram']['list']>([])
const serverFilterInput = ref('')
const serverList = ref<sftp[]>([])

getServerOption()

function getServerOption() {
  new ServerOption().request().then((response) => {
    serverOption.value = serverFilteredOption.value = response.data.list
  })
}

function filterServer(value: string) {
  serverFilteredOption.value = serverOption.value.filter((server) =>
    server.name.includes(value)
  )
}

function selectTab(sftp: sftp) {
  currentUUID.value = sftp.uuid
}

function deleteTab(sftp: sftp, index: number) {
  serverList.value.splice(index, 1)
  if (currentUUID.value === sftp.uuid) {
    currentUUID.value =
      serverList.value.length === 0
        ? 0
        : serverList.value[serverList.value.length - 1].uuid
  }
}

function selectServer(server: ServerData) {
  if (serverList.value.length === 0) {
    currentUUID.value = 0
  } else {
    currentUUID.value = serverList.value[serverList.value.length - 1].uuid + 1
  }
  serverList.value.push({ uuid: currentUUID.value, server, dir: '' })
  serverOptionVisible.value = false
}

function handleDirChange(dir: string) {
  const index = serverList.value.findIndex(
    (server) => server.uuid === currentUUID.value
  )
  if (index > -1) {
    serverList.value[index].dir = dir
  }
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
        cursor: pointer;
        text-overflow: ellipsis;
        overflow: hidden;
        width: 130px;
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
      width: 45px;
      border-bottom: 1px solid var(--el-border-color);
    }
  }
}
</style>
<style lang="scss">
.input-with-select .el-input-group__prepend {
  background-color: var(--el-bg-color);
}
.server-list {
  height: 216px;
  margin-top: 10px;
}
</style>
