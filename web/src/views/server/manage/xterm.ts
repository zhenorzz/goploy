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
  }
  public close(): void {
    this.websocket.close()
    this.terminal.dispose()
  }

  public send(message: string): void {
    this.websocket.send(message)
  }
}
