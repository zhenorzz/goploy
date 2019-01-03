<template>
  <el-container>
    <sidebar class="sidebar-container"/>
    <el-container style="margin-left:180px;">
      <el-header>
        <el-row type="flex" justify="end">
          <el-dropdown class="user-container" trigger="click">
            <div class="user-wrapper">
              {{ name }}
              <i class="el-icon-caret-bottom"/>
            </div>
            <el-dropdown-menu slot="dropdown" class="user-dropdown">
              <router-link class="inlineBlock" to="/userInfo">
                <el-dropdown-item>个人信息设置</el-dropdown-item>
              </router-link>
              <el-dropdown-item divided>
                <span style="display:block;" @click="logout">退出</span>
              </el-dropdown-item>
            </el-dropdown-menu>
          </el-dropdown>
        </el-row>
      </el-header>
      <el-main>
        <router-view></router-view>
      </el-main>
      <el-footer></el-footer>
    </el-container>
  </el-container>
</template>
<script>
import {mapGetters} from 'vuex';
import Sidebar from './Sidebar';
export default {
  name: 'Layout',
  components: {
    Sidebar,
  },
  computed: {
    ...mapGetters([
      'name',
    ]),
  },
  methods: {
    logout() {
      this.$store.dispatch('FedLogOut').then(() => {
        location.reload(); // 为了重新实例化vue-router对象 避免bug
      });
    },
  },
};
</script>
<style rel="stylesheet/scss" lang="scss">
#app {
  .el-header {
    border-bottom: solid 1px #e6e6e6;
    padding: 0, 20px;
    line-height: 60px;
    .user-wrapper {
      cursor: pointer;
    }
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
