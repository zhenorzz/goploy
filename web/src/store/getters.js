const getters = {
  sidebar: state => state.app.sidebar,
  device: state => state.app.device,
  uid: state => state.user.id,
  name: state => state.user.name,
  account: state => state.user.account,
  role: state => state.user.role,
  permission_routes: state => state.permission.routes,
  ws_message: state => state.websocket.message
}
export default getters
