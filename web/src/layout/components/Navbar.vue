<template>
  <div class="navbar">
    <img :src="logo" class="navbar-logo" />
    <hamburger
      :is-active="app.sidebar.opened"
      class="hamburger-container"
      @toggleClick="toggleSideBar"
    />
    <el-dropdown
      style="float: left; line-height: 48px; cursor: pointer"
      trigger="click"
      size="medium"
      placement="bottom-start"
      @command="handleNamespaceChange"
    >
      <span class="el-dropdown-link">
        {{ namespace.name }}<i class="el-icon-arrow-down el-icon--right" />
      </span>
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item
            v-for="item in namespaceList"
            :key="item.id"
            :command="item"
          >
            {{ item.name }}
          </el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
    <breadcrumb
      v-show="$store.state.app.device === 'desktop'"
      class="breadcrumb-container"
    />
    <div class="right">
      <div v-show="false" class="github">
        <span class="github-btn">
          <a class="gh-btn" href="https://github.com/zhenorzz" target="_blank">
            <span class="gh-ico" />
            <span class="gh-text">Follow @zhenorzz</span>
          </a>
        </span>
        <span class="github-btn">
          <a
            class="gh-btn"
            href="https://github.com/zhenorzz/goploy/"
            target="_blank"
          >
            <span class="gh-ico" />
            <span class="gh-text">Star</span>
          </a>
          <a
            class="gh-count"
            href="https://github.com/zhenorzz/goploy/stargazers"
            target="_blank"
          >
            {{ starCount }}
          </a>
        </span>
        <span class="github-btn github-forks">
          <a
            class="gh-btn"
            href="https://github.com/zhenorzz/goploy/fork"
            target="_blank"
          >
            <span class="gh-ico" aria-hidden="true" />
            <span class="gh-text">Fork</span>
          </a>
          <a
            class="gh-count"
            href="https://github.com/zhenorzz/goploy/network"
            target="_blank"
          >
            {{ forkCount }}
          </a>
        </span>
      </div>
      <div class="international">
        <el-dropdown
          trigger="click"
          size="medium"
          placement="bottom"
          @command="handleSetLanguage"
        >
          <div>
            <svg-icon class-name="international-icon" icon-class="language" />
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item
                :disabled="app.language === 'zh-cn'"
                command="zh-cn"
              >
                中文
              </el-dropdown-item>
              <el-dropdown-item :disabled="app.language === 'en'" command="en">
                English
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
      <div class="user-menu">
        <div class="user-container">
          <el-dropdown trigger="click" size="medium">
            <div class="user-wrapper">
              <el-row type="flex">
                <el-row>
                  <el-avatar
                    v-if="$store.getters.avatar"
                    :size="40"
                    :src="avatar"
                  />
                  <div
                    v-else
                    class="avatar-box"
                    style="background: rgb(64, 158, 255)"
                  >
                    <span>{{ user.name.substr(0, 1) }}</span>
                  </div>
                </el-row>
                <div style="margin-left: 8px">
                  <div class="user-name">{{ user.name }}</div>
                  <div class="user-title">{{ namespace.role }}</div>
                </div>
              </el-row>
            </div>
            <template #dropdown>
              <el-dropdown-menu class="user-dropdown">
                <router-link to="/user/profile">
                  <el-dropdown-item>
                    {{ $t('navbar.profile') }}
                  </el-dropdown-item>
                </router-link>
                <el-link
                  :underline="false"
                  href="https://docs.goploy.icu/"
                  target="__blank"
                >
                  <el-dropdown-item>
                    {{ $t('navbar.doc') }}
                  </el-dropdown-item>
                </el-link>
                <el-dropdown-item divided>
                  <span
                    style="display: block; text-align: center"
                    @click="logout"
                  >
                    {{ $t('navbar.logout') }}
                  </span>
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import logo from '@/assets/images/logo.png'
import { mapState } from 'vuex'
import Breadcrumb from '@/components/Breadcrumb/index.vue'
import Hamburger from '@/components/Hamburger/index.vue'
import {
  getNamespace,
  getNamespaceList,
  setNamespace,
  removeNamespaceIdCookie,
} from '@/utils/namespace'
import { ElLoading } from 'element-plus'
import { defineComponent } from 'vue'

export default defineComponent({
  components: {
    Breadcrumb,
    Hamburger,
  },
  data() {
    return {
      logo: logo,
      starCount: 0,
      forkCount: 0,
      namespace: getNamespace(),
      namespaceList: getNamespaceList(),
    }
  },
  computed: {
    ...mapState(['app', 'user']),
  },
  created() {
    // fetch('https://api.github.com/repos/zhenorzz/goploy').then(response => response.json()).then(data => {
    //   this.starCount = data.stargazers_count
    //   this.forkCount = data.forks_count
    // })
  },
  methods: {
    toggleSideBar() {
      this.$store.dispatch('app/toggleSideBar')
    },
    handleNamespaceChange(namespace) {
      setNamespace(namespace)
      ElLoading.service({ fullscreen: true })
      location.reload()
    },
    handleSetLanguage(lang) {
      this.$i18n.locale = lang
      this.$store.dispatch('app/setLanguage', lang)
      this.$message.success('Switch language success')
    },
    async logout() {
      await this.$store.dispatch('user/logout')
      await this.$store.dispatch('tagsView/delAllViews')
      removeNamespaceIdCookie()
      this.$router.push(`/login?redirect=${this.$route.fullPath}`)
    },
  },
})
</script>

<style lang="scss" scoped>
.navbar {
  height: 50px;
  overflow: hidden;
  position: relative;
  background: #fff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  &-logo {
    width: 25px;
    float: left;
    margin-left: 15px;
    margin-top: 11px;
  }
  .hamburger-container {
    line-height: 46px;
    height: 100%;
    float: left;
    cursor: pointer;
    transition: background 0.3s;
    -webkit-tap-highlight-color: transparent;

    &:hover {
      background: rgba(0, 0, 0, 0.025);
    }
  }

  .breadcrumb-container {
    float: left;
  }

  .right {
    float: right;
    width: auto;
  }

  .international {
    display: inline-block;
    line-height: 50px;
    cursor: pointer;
    &-icon {
      font-size: 18px;
    }
  }

  .github {
    display: inline-block;
    line-height: 50px;
    .github-btn {
      display: inline-block;
      font: 700 11px/14px 'Helvetica Neue', Helvetica, Arial, sans-serif;
      height: 20px;
      overflow: hidden;
      margin-right: 3px;
      position: relative;
      top: 5px;
      .gh-btn,
      .gh-count {
        padding: 2px 5px 2px 4px;
        color: #333;
        text-decoration: none;
        text-shadow: 0 1px 0 #fff;
        white-space: nowrap;
        cursor: pointer;
        border-radius: 3px;
      }
      .gh-btn,
      .gh-count,
      .gh-ico {
        float: left;
      }
      .gh-btn {
        background-color: #eee;
        background-image: linear-gradient(to bottom, #fcfcfc 0, #eee 100%);
        filter: progid:DXImageTransform.Microsoft.gradient(startColorstr='#fcfcfc', endColorstr='#eeeeee', GradientType=0);
        background-repeat: no-repeat;
        border: 1px solid #d5d5d5;
        position: relative;
        &:focus,
        &:hover {
          text-decoration: none;
          background-color: #ddd;
          background-image: linear-gradient(to bottom, #eee 0, #ddd 100%);
          filter: progid:DXImageTransform.Microsoft.gradient(startColorstr='#eeeeee', endColorstr='#dddddd', GradientType=0);
          border-color: #ccc;
        }
        &:active {
          background-image: none;
          background-color: #dcdcdc;
          border-color: #b5b5b5;
          box-shadow: inset 0 2px 4px rgba(0, 0, 0, 0.15);
        }
      }
      .gh-ico {
        width: 14px;
        height: 14px;
        margin-right: 4px;
        background-image: url(data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHhtbG5zOnhsaW5rPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5L3hsaW5rIiB2ZXJzaW9uPSIxLjEiIGlkPSJMYXllcl8xIiB4PSIwcHgiIHk9IjBweCIgd2lkdGg9IjQwcHgiIGhlaWdodD0iNDBweCIgdmlld0JveD0iMTIgMTIgNDAgNDAiIGVuYWJsZS1iYWNrZ3JvdW5kPSJuZXcgMTIgMTIgNDAgNDAiIHhtbDpzcGFjZT0icHJlc2VydmUiPjxwYXRoIGZpbGw9IiMzMzMzMzMiIGQ9Ik0zMiAxMy40Yy0xMC41IDAtMTkgOC41LTE5IDE5YzAgOC40IDUuNSAxNS41IDEzIDE4YzEgMC4yIDEuMy0wLjQgMS4zLTAuOWMwLTAuNSAwLTEuNyAwLTMuMiBjLTUuMyAxLjEtNi40LTIuNi02LjQtMi42QzIwIDQxLjYgMTguOCA0MSAxOC44IDQxYy0xLjctMS4yIDAuMS0xLjEgMC4xLTEuMWMxLjkgMC4xIDIuOSAyIDIuOSAyYzEuNyAyLjkgNC41IDIuMSA1LjUgMS42IGMwLjItMS4yIDAuNy0yLjEgMS4yLTIuNmMtNC4yLTAuNS04LjctMi4xLTguNy05LjRjMC0yLjEgMC43LTMuNyAyLTUuMWMtMC4yLTAuNS0wLjgtMi40IDAuMi01YzAgMCAxLjYtMC41IDUuMiAyIGMxLjUtMC40IDMuMS0wLjcgNC44LTAuN2MxLjYgMCAzLjMgMC4yIDQuNyAwLjdjMy42LTIuNCA1LjItMiA1LjItMmMxIDIuNiAwLjQgNC42IDAuMiA1YzEuMiAxLjMgMiAzIDIgNS4xYzAgNy4zLTQuNSA4LjktOC43IDkuNCBjMC43IDAuNiAxLjMgMS43IDEuMyAzLjVjMCAyLjYgMCA0LjYgMCA1LjJjMCAwLjUgMC40IDEuMSAxLjMgMC45YzcuNS0yLjYgMTMtOS43IDEzLTE4LjFDNTEgMjEuOSA0Mi41IDEzLjQgMzIgMTMuNHoiLz48L3N2Zz4=);
        background-size: 100% 100%;
        background-repeat: no-repeat;
      }
      .gh-count {
        position: relative;
        display: none;
        margin-left: 4px;
        background-color: #fafafa;
        border: 1px solid #d4d4d4;
        z-index: 1;
        display: block;
        &:focus,
        &:hover {
          color: #4183c4;
        }
        &:after,
        &:before {
          content: '';
          position: absolute;
          display: inline-block;
          width: 0;
          height: 0;
          border-color: transparent;
          border-style: solid;
        }
        &:before {
          top: 50%;
          left: -2px;
          margin-top: -3px;
          border-width: 2px 2px 2px 0;
          border-right-color: #fafafa;
        }

        &:after {
          top: 50%;
          left: -3px;
          z-index: -1;
          margin-top: -4px;
          border-width: 3px 3px 3px 0;
          border-right-color: #d4d4d4;
        }
      }
    }
  }
  .user-menu {
    float: right;
    height: 100%;

    &:focus {
      outline: none;
    }

    &:hover {
      background-color: #f5f7fa;
    }

    .user-container {
      padding: 0 20px;
      height: 50px;
      line-height: 0;
      .user-wrapper {
        position: relative;
        cursor: pointer;
        margin-top: 5px;
      }

      .avatar-box {
        height: 40px;
        width: 40px;
        line-height: 40px;
        border-radius: 50%;
        text-align: center;
        color: #fff;
        font-size: 16px;
      }

      .user-name {
        margin-top: 4px;
        font-size: 16px;
        font-weight: 900;
        color: #9d9d9d;
      }
      .user-title {
        font-size: 13px;
        padding-top: 3px;
        color: #9d9d9d;
      }
    }
  }
}
</style>
