const getters = {
  sidebar: state => state.app.sidebar,
  device: state => state.app.device,
  language: state => state.app.language,
  uid: state => state.user.id,
  name: state => state.user.name,
  account: state => state.user.account,
  superManager: state => state.user.superManager,
  permission_routes: state => state.permission.routes,
  ws_message: state => state.websocket.message,
  ws: state => state.websocket.ws
}
export default getters
