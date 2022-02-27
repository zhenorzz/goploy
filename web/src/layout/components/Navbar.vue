<template>
  <div class="navbar">
    <img :src="logo" class="navbar-logo" />
    <hamburger
      :is-active="app.sidebar.opened"
      class="hamburger-container"
      @toggle-click="toggleSideBar"
    />
    <el-dropdown
      style="float: left; line-height: 48px; cursor: pointer"
      trigger="click"
      size="medium"
      placement="bottom-start"
      @visible-change="handleNamespaceVisible"
      @command="handleNamespaceChange"
    >
      <span class="el-dropdown-link">
        {{ namespace.name }}<i class="el-icon-arrow-down el-icon--right" />
      </span>
      <template #dropdown>
        <el-dropdown-menu v-loading="namespaceListLoading">
          <el-dropdown-item
            v-for="item in namespaceList"
            :key="item.namespaceId"
            :command="item"
          >
            {{ item.namespaceName }}
          </el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
    <breadcrumb
      v-show="$store.state.app.device === 'desktop'"
      class="breadcrumb-container"
    />
    <div class="right">
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
                    :src="$store.getters.avatar"
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

<script lang="ts" setup>
import logo from '@/assets/images/logo.png'
import Breadcrumb from '@/components/Breadcrumb/index.vue'
import Hamburger from '@/components/Hamburger/index.vue'
import { useI18n } from 'vue-i18n'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
import { NamespaceUserData, NamespaceOption } from '@/api/namespace'
import { getNamespace, setNamespace, removeNamespace } from '@/utils/namespace'
import { ElLoading, ElMessage } from 'element-plus'
import { ref, computed } from 'vue'

const namespaceListLoading = ref(false)
const namespaceList = ref<NamespaceOption['datagram']['list']>()
const namespace = ref(getNamespace())
const { locale } = useI18n({ useScope: 'global' })
const router = useRouter()
const store = useStore()
const app = computed(() => store.state['app'])
const user = computed(() => store.state['user'])

document.title = `Goploy-${namespace.value.name}`

function toggleSideBar() {
  store.dispatch('app/toggleSideBar')
}
function handleNamespaceVisible(visible: boolean) {
  if (visible === true) {
    namespaceListLoading.value = true
    new NamespaceOption()
      .request()
      .then((response) => {
        namespaceList.value = response.data.list
      })
      .finally(() => {
        namespaceListLoading.value = false
      })
  }
}
function handleNamespaceChange(namespace: NamespaceUserData['datagram']) {
  setNamespace({
    id: namespace.namespaceId,
    name: namespace.namespaceName,
    role: namespace.role,
  })
  ElLoading.service({ fullscreen: true })
  location.reload()
}
function handleSetLanguage(lang: string) {
  locale.value = lang
  store.dispatch('app/setLanguage', lang)
  ElMessage.success('Switch language success')
}

async function logout() {
  await store.dispatch('user/logout')
  await store.dispatch('tagsView/delAllViews')
  removeNamespace()
  router.push(`/login`)
}
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
