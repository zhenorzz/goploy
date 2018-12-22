const getters = {
  token: (state) => state.user.token,
  id: (state) => state.user.id,
  name: (state) => state.user.name,
  role: (state) => state.user.role,
  routers: (state) => state.permission.routers,
  addRouters: (state) => state.permission.addRouters,
};
export default getters;
