<template>
  <el-row class="app-container main">
    <el-row class="header">
      <el-scrollbar>
        <el-row class="nav" justify="start" align="middle">
          <div
            v-for="(item, index) in terminalList"
            :key="item.uuid"
            class="nav-item"
            :class="
              item.uuid === currentTerminalUUID ? 'nav-item-selected' : ''
            "
          >
            <el-row>
              <el-row class="nav-item-name">
                <el-button
                  type="text"
                  style="color: #bfcbd9; font-size: 14px"
                  :title="`${item.server.name}(${item.server.description})`"
                  @click="selectTerminal(item)"
                >
                  {{ item.server.name }}({{ item.server.description }})
                </el-button>
              </el-row>
              <el-button
                type="text"
                style="
                  color: #bfcbd9;
                  font-size: 14px;
                  padding-top: 20px;
                  padding-left: 8px;
                "
                @click="deleteTerminal(item, index)"
              >
                x
              </el-button>
            </el-row>
          </div>
          <div class="nav-plus">
            <el-popover
              v-model:visible="serverOptionVisible"
              placement="bottom-start"
              :width="200"
              trigger="click"
              popper-class="nav-popover"
              :show-arrow="false"
              :offset="-2"
              style="background: #0e0f12"
            >
              <el-input
                v-model="serverFilterInput"
                placeholder="Filter server name"
                :input-style="{
                  background: '#0e0f12',
                  borderColor: '#304156',
                  color: '#bfcbd9',
                }"
                @input="filterServer"
              />
              <div class="server-list">
                <div v-for="server in serverFilteredOption" :key="server.id">
                  <el-button
                    link
                    class="server-item"
                    @click="selectServer(server)"
                  >
                    <span :title="server.name + '(' + server.description + ')'">
                      {{ server.name }}({{ server.description }})
                    </span>
                  </el-button>
                </div>
              </div>
              <template #reference>
                <el-button
                  link
                  style="
                    color: #bfcbd9;
                    height: 45px;
                    width: 100%;
                    font-size: 24px;
                  "
                  @click="serverOptionVisible = !serverOptionVisible"
                >
                  +
                </el-button>
              </template>
            </el-popover>
          </div>
        </el-row>
      </el-scrollbar>
    </el-row>
    <el-row class="terminal">
      <div
        v-for="item in terminalList"
        v-show="item.uuid === currentTerminalUUID"
        :key="item.uuid"
        :ref="
          (el) => {
            if (el) {
              terminalRefs[item.uuid] = el
            }
          }
        "
        style="width: 100%; height: 100%"
      ></div>
    </el-row>
    <el-row class="footer">
      <el-input
        v-model="command"
        :disabled="terminalList.length === 0"
        placeholder="Click here send to all windows"
        :input-style="{
          background: '#0e0f12',
          borderColor: '#304156',
          color: '#bfcbd9',
        }"
        @keyup.enter="enterCommand"
      />
    </el-row>
  </el-row>
</template>
<script lang="ts">
export default { name: 'ServerTerminal' }
</script>
<script lang="ts" setup>
import 'xterm/css/xterm.css'
import { ServerOption, ServerData } from '@/api/server'
import { xterm } from './xterm'
import { ref, nextTick, ComponentPublicInstance } from 'vue'
interface terminal {
  uuid: number
  xterm?: xterm
  server: ServerData
}
const serverOptionVisible = ref(false)
const terminalList = ref<terminal[]>([])
const currentTerminalUUID = ref(0)
const serverOption = ref<ServerOption['datagram']['list']>([])
const serverFilteredOption = ref<ServerOption['datagram']['list']>([])
const serverFilterInput = ref('')
const command = ref('')
const terminalRefs = ref<Record<string, Element | ComponentPublicInstance>>({})

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
function selectServer(server: ServerData) {
  if (terminalList.value.length === 0) {
    currentTerminalUUID.value = 0
  } else {
    currentTerminalUUID.value =
      terminalList.value[terminalList.value.length - 1].uuid + 1
  }
  terminalList.value.push({ uuid: currentTerminalUUID.value, server })
  serverOptionVisible.value = false
  nextTick(() => {
    const x = new xterm(
      terminalRefs.value[currentTerminalUUID.value] as HTMLDivElement,
      server.id
    )
    x.connect()
    terminalList.value[terminalList.value.length - 1].xterm = x
  })
}
function selectTerminal(terminal: terminal) {
  currentTerminalUUID.value = terminal.uuid
}
function deleteTerminal(terminal: terminal, index: number) {
  terminal.xterm?.close()
  terminalList.value.splice(index, 1)
  if (currentTerminalUUID.value === terminal.uuid) {
    currentTerminalUUID.value =
      terminalList.value.length === 0
        ? 0
        : terminalList.value[terminalList.value.length - 1].uuid
  }
}
function enterCommand() {
  terminalList.value.forEach((terminal) => {
    terminal.xterm?.send(command.value + '\n')
  })
  command.value = ''
}
</script>

<style lang="scss" scoped>
.main {
  flex-direction: column;
  height: calc(100vh - 84px);
}

.header {
  width: 100%;
  background: #0e0f12;
}
.nav {
  height: auto;
  flex-wrap: nowrap;
  &-item {
    color: #bfcbd9;
    font-size: 14px;
    flex-shrink: 0;
    width: 150px;
    height: 45px;
    line-height: 45px;
    padding: 0 10px;
    &-name {
      cursor: pointer;
      text-align: center;
      text-overflow: ellipsis;
      overflow: hidden;
      width: 110px;
      display: inline-block;
    }
    &-selected {
      background: #1d2935;
    }
  }
  &-plus {
    text-align: center;
    width: 40px;
  }
}

.terminal {
  flex: 1;
  height: 100%;
  padding: 10px;
  background: #1d2935;
}

.footer {
  height: 28px;
  width: 100%;
  background: #0e0f12;
}

@media only screen and (max-device-width: 992px) {
  .terminal {
    width: calc(100vw - 40px);
  }
}
</style>
<style lang="scss">
@import '@/styles/mixin.scss';
.xterm .xterm-viewport {
  overflow: auto;
  @include scrollBar();
}
.nav-popover.el-popper {
  background: #0e0f12;
  border-color: #0e0f12;
}
.server-list {
  height: 216px;
  overflow-x: auto;
  margin-top: 10px;
  @include scrollBar();
}
.server-item {
  color: #bfcbd9;
  width: 150px;
  text-align: left;
  overflow-x: scroll;
  display: inline-block;
  &::-webkit-scrollbar {
    display: none;
  }
}
</style>
