import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import { AttachAddon } from 'xterm-addon-attach'
export class xterm {
  private serverId: number
  private element: HTMLDivElement
  private websocket!: WebSocket
  private terminal: Terminal | null = null
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
    })
    const fitAddon = new FitAddon()
    this.terminal.loadAddon(fitAddon)
    this.terminal.open(this.element)
    fitAddon.fit()
    this.websocket = new WebSocket(
      `${location.protocol.replace('http', 'ws')}//${
        window.location.host + import.meta.env.VITE_APP_BASE_API
      }/ws/xterm?serverId=${this.serverId}&rows=${this.terminal.rows}&cols=${
        this.terminal.cols
      }`
    )
    const attachAddon = new AttachAddon(this.websocket)
    this.terminal.loadAddon(attachAddon)
  }
  public close(): void {
    this.terminal = null
    this.websocket.close()
  }
}
