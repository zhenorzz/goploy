<template>
  <el-container>
    <sidebar class="sidebar-container"/>
    <el-container style="margin-left:180px;">
      <el-header></el-header>
      <el-main>
        <router-view></router-view>
      </el-main>
      <el-footer></el-footer>
    </el-container>
  </el-container>
</template>
<script>
import Sidebar from './Sidebar';
export default {
  name: 'Layout',
  components: {
    Sidebar,
  },
};
</script>
<style rel="stylesheet/scss" lang="scss">
  #app {
    .el-header {
      border-bottom: solid 1px #e6e6e6;
      padding: 0;
    }
    .el-main {
      height: calc(100vh - 120px);
    }
    .el-footer {
      position: fixed;
      bottom: 0;
      width: 100%;
      border-top: solid 1px #e6e6e6;
    }
  }
</style>
