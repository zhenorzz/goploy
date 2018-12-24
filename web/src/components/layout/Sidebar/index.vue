<template>
  <el-scrollbar wrap-class="scrollbar-wrapper">
    <el-row class="logo" type="flex">
      <svg-icon icon-class="logo" class="logo-pic"/>
      <span class="logo-text">goploy</span>
    </el-row>
    <el-menu
      :show-timeout="200"
      :default-active="$route.path"
      mode="vertical"
      background-color="#304156"
      text-color="#bfcbd9"
      active-text-color="#409EFF"
    >
      <sidebar-item v-for="route in routers" :key="route.path" :item="route" :base-path="route.path"/>
    </el-menu>
  </el-scrollbar>
</template>

<script>
import {mapGetters} from 'vuex';
import SidebarItem from './SidebarItem';

export default {
  components: {SidebarItem},
  computed: {
    ...mapGetters([
      'routers',
    ]),
  },
};
</script>
<style rel="stylesheet/scss" lang="scss" scoped>
.logo {
  position: relative;
  width: 180px;
  height: 49px;
  padding: 0;
  font-size: 12px;
  color: #fff;
  line-height: 49px;
  background-color: #304156;
  padding: 0 16px;
  &-pic {
    margin-top: 5px;
    margin-right: 10px !important;
    font-size:40px;
    color:#bfcbd9;
  }
  &-text {
    font-size:20px;
    color:#bfcbd9;
  }
}
</style>
