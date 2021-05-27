export interface WebsocketState {
  ws: WebSocket | null
  message: any
  againConnectTime: number
}
