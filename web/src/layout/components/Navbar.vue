<template>
  <el-row class="navbar" align="middle">
    <el-dropdown
      class="navbar-namespace"
      trigger="click"
      placement="bottom"
      @command="handleNamespaceChange"
    >
      <el-row align="middle">
        <img :src="logo" class="navbar-logo" />
        <span>{{ namespace.name }}</span>
        <el-icon class="el-icon--right">
          <arrow-down />
        </el-icon>
      </el-row>
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
    <Sidebar class="sidebar-container" style="flex: 1" />
    <div class="international">
      <el-switch
        v-model="isDark"
        style="
          --el-switch-on-color: var(--el-border-color);
          --el-switch-off-color: var(--el-border-color);
          margin-right: 20px;
          height: 56px;
        "
        inline-prompt
        :active-icon="Moon"
        :inactive-icon="Sunny"
      />
      <el-dropdown
        trigger="click"
        placement="bottom"
        @command="showTransformDialog"
      >
        <div style="height: 100%; padding-top: 2px; margin-right: 20px">
          <svg-icon class-name="international-icon" :icon-class="'toolbox'" />
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="time">Date Transform</el-dropdown-item>
            <el-dropdown-item command="json">JSON Pretty</el-dropdown-item>
            <el-dropdown-item command="password">Random PWD</el-dropdown-item>
            <el-dropdown-item command="unicode">Unicode</el-dropdown-item>
            <el-dropdown-item command="decodeURI">DecodeURI</el-dropdown-item>
            <el-dropdown-item command="md5">MD5 </el-dropdown-item>
            <el-dropdown-item command="cron">Crontab</el-dropdown-item>
            <el-dropdown-item command="qrcode">QRcode</el-dropdown-item>
            <el-dropdown-item command="byte">Byte Transform</el-dropdown-item>
            <el-dropdown-item command="color">Color Transform</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
      <el-dropdown
        trigger="click"
        placement="bottom"
        @command="handleSetLanguage"
      >
        <div style="height: 100%; padding-top: 2px">
          <svg-icon class-name="international-icon" :icon-class="'global'" />
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
      <div class="user-menu">
        <div class="user-container">
          <el-dropdown trigger="click">
            <div class="user-wrapper">
              <div class="user-name">
                {{ user.name }}
              </div>
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
                <el-dropdown-item divided @click="logout">
                  <span style="display: block; text-align: center">
                    {{ $t('navbar.logout') }}
                  </span>
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>
    </div>
    <el-dialog
      v-model="transformVisible"
      width="600px"
      append-to-body
      :close-on-click-modal="false"
      :fullscreen="$store.state.app.device === 'mobile'"
    >
      <el-row class="transform-content">
        <TheDatetransform v-model="transformType" />
        <TheJSONPretty v-model="transformType" />
        <TheRandomPWD v-model="transformType" />
        <TheUnicode v-model="transformType" />
        <el-row v-show="transformType === 'decodeURI'" style="width: 100%">
          <el-input
            v-model="uri.escape"
            type="textarea"
            :autosize="{ minRows: 2 }"
            placeholder="Please enter unescaped URI"
          />
          <el-input
            :value="uri.escape ? decodeURI(uri.escape) : ''"
            style="margin-top: 10px"
            type="textarea"
            :autosize="{ minRows: 2 }"
            readonly
          />
        </el-row>
        <el-row v-show="transformType === 'md5'" style="width: 100%">
          <el-input
            v-model="md5.text"
            type="textarea"
            style="width: 100%"
            :autosize="{ minRows: 3 }"
          />
          <el-input
            :value="hashByMD5(md5.text)"
            style="width: 100%; margin-top: 10px"
            readonly
          />
        </el-row>

        <el-row v-show="transformType === 'qrcode'">
          <el-input
            v-model="qrcode.text"
            type="textarea"
            :autosize="{ minRows: 2 }"
          />
          <el-row style="margin-top: 10px" type="flex" align="middle">
            <span style="width: 30px; font-size: 14px; margin-right: 10px">
              Size
            </span>
            <el-input-number v-model="qrcode.width" />
          </el-row>
          <VueQrcode
            class="text-align:center"
            :value="qrcode.text"
            :options="{ width: qrcode.width }"
          />
        </el-row>
        <TheCronstrue v-model="transformType" />
        <TheByteTransform v-model="transformType" />
        <TheRGBTransform v-model="transformType" />
      </el-row>
    </el-dialog>
  </el-row>
</template>

<script lang="ts" setup>
import logo from '@/assets/images/logo.png'
import Sidebar from './Sidebar/index.vue'
import VueQrcode from '@chenfengyuan/vue-qrcode'
import { md5 as hashByMD5 } from '@/utils/md5'
import TheDatetransform from './Toolbox/TheDatetransform.vue'
import TheJSONPretty from './Toolbox/TheJSONPretty.vue'
import TheRandomPWD from './Toolbox/TheRandomPWD.vue'
import TheUnicode from './Toolbox/TheUnicode.vue'
import TheCronstrue from './Toolbox/TheCronstrue.vue'
import TheByteTransform from './Toolbox/TheByteTransform.vue'
import TheRGBTransform from './Toolbox/TheRGBTransform.vue'
import { ArrowDown } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
import { NamespaceUserData, NamespaceOption } from '@/api/namespace'
import { getNamespace, setNamespace, removeNamespace } from '@/utils/namespace'
import { ref, reactive, computed } from 'vue'
import { useDark } from '@vueuse/core'
import { Sunny, Moon } from '@element-plus/icons-vue'
const isDark = useDark()
const transformVisible = ref(false)
const transformType = ref('')
const qrcode = reactive({
  text: 'https://github.com/zhenorzz/goploy',
  width: 200,
})
const uri = reactive({
  escape: '',
})
const md5 = reactive({
  text: '',
})
const namespaceListLoading = ref(false)
const namespaceList = ref<NamespaceOption['datagram']['list']>()
const namespace = ref(getNamespace())
const { locale } = useI18n({ useScope: 'global' })
const router = useRouter()
const store = useStore()
const app = computed(() => store.state['app'])
const user = computed(() => store.state['user'])

document.title = `Goploy-${namespace.value.name}`

function showTransformDialog(type: string) {
  transformVisible.value = true
  transformType.value = type
}

getNamespaceList()

function getNamespaceList() {
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

function handleNamespaceChange(namespace: NamespaceUserData) {
  setNamespace({
    id: namespace.namespaceId,
    name: namespace.namespaceName,
    roleId: namespace.roleId,
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
  width: 100vw;
  height: 56px;
  overflow: hidden;
  position: relative;
  background-color: var(--el-bg-color);
  &-namespace {
    padding: 0 10px;
    line-height: 56px;
    cursor: pointer;
    &:hover {
      color: var(--el-menu-hover-text-color);
      background-color: var(--el-menu-hover-bg-color);
    }
  }
  // box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  &-logo {
    width: 30px;
    margin-right: 10px;
  }

  .international {
    display: inline-block;
    cursor: pointer;
    margin-left: 20px;
    &-icon {
      font-size: 17px;
    }
    .el-dropdown {
      line-height: 54px;
      height: 56px;
    }
  }
  .user-menu {
    float: right;
    height: 100%;

    &:focus {
      outline: none;
    }

    .user-container {
      padding: 0 20px;
      height: 56px;
      line-height: 0;
      .user-wrapper {
        position: relative;
        cursor: pointer;
        height: 56px;
      }

      .user-name {
        font-size: 15px;
        font-weight: 900;
        line-height: 56px;
        color: var(--el-text-color-regular);
        &:hover {
          color: var(--el-color-primary);
        }
      }
    }
  }
}
</style>
