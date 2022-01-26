<template>
  <el-row class="app-container main">
    <el-row class="header">
      <el-scrollbar>
        <el-row class="nav" justify="start" align="middle">
          <div
            v-for="(terminal, index) in terminalList"
            :key="terminal.uuid"
            class="nav-item"
            :class="
              terminal.uuid === currentTerminalUUID ? 'nav-item-selected' : ''
            "
          >
            <el-row>
              <el-row class="nav-item-name">
                <el-button
                  type="text"
                  style="color: #bfcbd9; font-size: 14px"
                  @click="selectTerminal(terminal)"
                >
                  {{ terminal.server.name }}
                </el-button>
              </el-row>
              <el-button
                type="text"
                style="color: #bfcbd9"
                icon="el-icon-close"
                @click="deleteTerminal(terminal, index)"
              ></el-button>
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
              <el-row class="server-list" type="flex" align="middle">
                <el-row v-for="server in serverFilteredOption" :key="server.id">
                  <el-button
                    type="text"
                    style="
                      color: #bfcbd9;
                      width: 150px;
                      text-align: left;
                      text-overflow: ellipsis;
                      overflow: hidden;
                    "
                    size="medium"
                    @click="selectServer(server)"
                  >
                    {{ server.name }}
                  </el-button>
                </el-row>
              </el-row>
              <template #reference>
                <el-button
                  type="text"
                  style="color: #bfcbd9; height: 45px; width: 100%"
                  icon="el-icon-plus"
                  @click="serverOptionVisible = !serverOptionVisible"
                ></el-button>
              </template>
            </el-popover>
          </div>
        </el-row>
      </el-scrollbar>
    </el-row>
    <el-row class="terminal">
      <div
        v-for="terminal in terminalList"
        v-show="terminal.uuid === currentTerminalUUID"
        :key="terminal.uuid"
        :ref="`terminal${terminal.uuid}`"
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
import 'xterm/css/xterm.css'
import { ServerOption, ServerData } from '@/api/server'
import { xterm } from './xterm'
import { defineComponent } from 'vue'
interface terminal {
  uuid: number
  xterm?: xterm | void
  server: ServerData['datagram']
}
export default defineComponent({
  name: 'ServerTerminal',
  data() {
    return {
      serverOptionVisible: false,
      terminalList: [] as terminal[],
      currentTerminalUUID: 0,
      serverOption: [] as ServerOption['datagram']['list'],
      serverFilteredOption: [] as ServerOption['datagram']['list'],
      serverFilterInput: '',
      command: '',
    }
  },
  created() {
    this.getServerOption()
  },
  methods: {
    getServerOption() {
      new ServerOption().request().then((response) => {
        this.serverOption = this.serverFilteredOption = response.data.list
      })
    },
    filterServer(value: string) {
      this.serverFilteredOption = this.serverOption.filter((server) =>
        server.name.includes(value)
      )
    },
    selectServer(server: ServerData['datagram']) {
      if (this.terminalList.length === 0) {
        this.currentTerminalUUID = 0
      } else {
        this.currentTerminalUUID =
          this.terminalList[this.terminalList.length - 1].uuid + 1
      }
      this.terminalList.push({ uuid: this.currentTerminalUUID, server })
      this.serverOptionVisible = false
      this.$nextTick(() => {
        const x = new xterm(
          this.$refs[`terminal${this.currentTerminalUUID}`] as HTMLDivElement,
          server.id
        )
        x.connect()
        this.terminalList[this.terminalList.length - 1].xterm = x
      })
    },
    selectTerminal(terminal: terminal) {
      this.currentTerminalUUID = terminal.uuid
    },
    deleteTerminal(terminal: terminal, index: number) {
      terminal.xterm?.close()
      this.terminalList.splice(index, 1)
      if (this.currentTerminalUUID === terminal.uuid) {
        this.currentTerminalUUID =
          this.terminalList.length === 0
            ? 0
            : this.terminalList[this.terminalList.length - 1].uuid
      }
    },
    enterCommand() {
      this.terminalList.forEach((terminal) => {
        terminal.xterm?.send(this.command + '\n')
      })
      this.command = ''
    },
  },
})
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
  flex-wrap: nowrap;
  flex-direction: column;
  margin-left: 10px;
  margin-top: 10px;
  @include scrollBar();
}
</style>
