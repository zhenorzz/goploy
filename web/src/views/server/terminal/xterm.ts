import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import { AttachAddon } from 'xterm-addon-attach'
import { NamespaceKey, getNamespaceId } from '@/utils/namespace'
export class xterm {
  private serverId: number
  private element: HTMLDivElement
  private websocket!: WebSocket
  private terminal!: Terminal
  constructor(element: HTMLDivElement, serverId: number) {
    this.element = element
    this.serverId = serverId
  }
  public connect(): void {
    const isWindows =
      ['Windows', 'Win16', 'Win32', 'WinCE'].indexOf(navigator.platform) >= 0
    this.terminal = new Terminal({
      fontSize: 14,
      cursorBlink: true,
      windowsMode: isWindows,
      theme: {
        foreground: '#ebeef5',
        background: '#1d2935',
        cursor: '#e6a23c',
        black: '#000000',
        brightBlack: '#555555',
        red: '#ef4f4f',
        brightRed: '#ef4f4f',
        green: '#67c23a',
        brightGreen: '#67c23a',
        yellow: '#e6a23c',
        brightYellow: '#e6a23c',
        blue: '#409eff',
        brightBlue: '#409eff',
        magenta: '#ef4f4f',
        brightMagenta: '#ef4f4f',
        cyan: '#17c0ae',
        brightCyan: '#17c0ae',
        white: '#bbbbbb',
        brightWhite: '#ffffff',
      },
    })
    const fitAddon = new FitAddon()
    this.terminal.open(this.element)
    this.terminal.loadAddon(fitAddon)
    fitAddon.fit()
    this.websocket = new WebSocket(
      `${location.protocol.replace('http', 'ws')}//${
        window.location.host + import.meta.env.VITE_APP_BASE_API
      }/ws/xterm?${NamespaceKey}=${getNamespaceId()}&serverId=${
        this.serverId
      }&rows=${this.terminal.rows}&cols=${this.terminal.cols}`
    )
    this.terminal.loadAddon(new AttachAddon(this.websocket))
    this.websocket.onclose = function (evt) {
      if (evt.reason !== '') {
        ElMessage.error(evt.reason)
      }
    }
  }
  public close(): void {
    this.terminal.dispose()
    this.websocket.close()
  }

  public send(message: string): void {
    this.websocket.send(message)
  }
}
