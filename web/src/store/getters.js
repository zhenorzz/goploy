const getters = {
  sidebar: state => state.app.sidebar,
  device: state => state.app.device,
  token: state => state.user.token,
  uid: state => state.user.id,
  name: state => state.user.name,
  account: state => state.user.account,
  permission_routes: state => state.permission.routes,
  permission_uri: state => state.permission.permissionUri
}
export default getters
