<template>
  <el-container>
    <el-aside width="200px">
      <el-menu
        default-active="2"
        mode="vertical"
        :router="true"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409EFF"
        style="height:100%">
        <el-menu-item index="dashboard">
          <i class="el-icon-menu"></i>
          <span slot="title">dashboard</span>
        </el-menu-item>
        <el-submenu index="1">
          <template slot="title">
            <i class="el-icon-location"></i>
            <span>路由测试</span>
          </template>
          <el-menu-item index="1-1">选项1</el-menu-item>
          <el-menu-item index="admin">选项2</el-menu-item>
        </el-submenu>

        <el-menu-item index="3" disabled>
          <i class="el-icon-document"></i>
          <span slot="title">导航三</span>
        </el-menu-item>
        <el-menu-item index="4">
          <i class="el-icon-setting"></i>
          <span slot="title">导航四</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header></el-header>
      <el-main>
        <router-view></router-view>
      </el-main>
      <el-footer></el-footer>
    </el-container>
  </el-container>
</template>
<style rel="stylesheet/scss" lang="scss">
  #app {
    .el-aside {
      height: 100vh;
    }

    .el-submenu .el-menu-item {
      min-width: 0;
    }

    .el-header {
      border-bottom: solid 1px #e6e6e6;
      padding: 0;
    }

    .el-footer {
      border-top: solid 1px #e6e6e6;
    }
  }
</style>
