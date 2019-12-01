<template>
  <div class="navbar">
    <hamburger :is-active="sidebar.opened" class="hamburger-container" @toggleClick="toggleSideBar" />

    <breadcrumb class="breadcrumb-container" />

    <div class="right-menu">
      <iframe src="https://ghbtns.com/github-btn.html?user=zhenorzz&type=follow&count=true" frameborder="0" scrolling="0" width="130px" height="20px" style="position: relative; top: 5px;" />
      <iframe src="https://ghbtns.com/github-btn.html?user=zhenorzz&repo=goploy&type=star&count=true" frameborder="0" scrolling="0" width="80px" height="20px" style="position: relative; top: 5px;" />
      <iframe src="https://ghbtns.com/github-btn.html?user=zhenorzz&repo=goploy&type=fork&count=true" frameborder="0" scrolling="0" width="80px" height="20px" style="position: relative; top: 5px;" />
      <el-dropdown class="user-container" trigger="click" size="medium">
        <div class="user-wrapper">
          {{ name }}
          <i class="el-icon-caret-bottom" />
        </div>
        <el-dropdown-menu slot="dropdown" class="user-dropdown">
          <router-link to="/user/info">
            <el-dropdown-item>
              个人信息设置
            </el-dropdown-item>
          </router-link>
          <el-dropdown-item divided>
            <span style="display:block;" @click="logout">退出</span>
          </el-dropdown-item>
        </el-dropdown-menu>
      </el-dropdown>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'
import Breadcrumb from '@/components/Breadcrumb'
import Hamburger from '@/components/Hamburger'

export default {
  components: {
    Breadcrumb,
    Hamburger
  },
  computed: {
    ...mapGetters([
      'sidebar',
      'name'
    ])
  },
  methods: {
    toggleSideBar() {
      this.$store.dispatch('app/toggleSideBar')
    },
    async logout() {
      await this.$store.dispatch('user/logout')
      this.$router.push(`/login?redirect=${this.$route.fullPath}`)
    }
  }
}
</script>

<style lang="scss" scoped>
.navbar {
  height: 50px;
  overflow: hidden;
  position: relative;
  background: #fff;
  box-shadow: 0 1px 4px rgba(0,21,41,.08);

  .hamburger-container {
    line-height: 46px;
    height: 100%;
    float: left;
    cursor: pointer;
    transition: background .3s;
    -webkit-tap-highlight-color:transparent;

    &:hover {
      background: rgba(0, 0, 0, .025)
    }
  }

  .breadcrumb-container {
    float: left;
  }

  .right-menu {
    float: right;
    height: 100%;
    line-height: 50px;

    &:focus {
      outline: none;
    }
    .star-btn {
      padding: 3px 5px;
      font-size: 12px;
      line-height: 20px;
      color: #24292e;
      background-color: #eff3f6;
      background-image: linear-gradient(-180deg,#fafbfc,#eff3f6 90%);
      border: 1px solid rgba(2i7,31,35,.2);
      border-radius: .25em;
      span {
        font-size: 14px !important;
      }
    }
    .right-menu-item {
      display: inline-block;
      padding: 0 8px;
      height: 100%;
      font-size: 18px;
      color: #5a5e66;
      vertical-align: text-bottom;

      &.hover-effect {
        cursor: pointer;
        transition: background .3s;

        &:hover {
          background: rgba(0, 0, 0, .025)
        }
      }
    }

    .user-container {
      margin-right: 30px;

      .user-wrapper {
        position: relative;
      }
    }
  }
}
</style>
