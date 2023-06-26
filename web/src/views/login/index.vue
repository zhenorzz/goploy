<template>
  <div v-if="!query['code'] && !query['authCode']" class="login-container">
    <el-row class="login-mark" type="flex" align="middle" justify="center">
      <img src="@/assets/images/logo.png" width="120" height="120" />
    </el-row>
    <el-form
      ref="form"
      :model="loginForm"
      :rules="loginRules"
      class="login-form"
      auto-complete="on"
      label-position="left"
    >
      <div class="title-container">
        <h3 class="title">
          Sign in to Goploy <sub>{{ version }}</sub>
        </h3>
      </div>

      <el-form-item prop="account" class="login-form-input">
        <span class="svg-container">
          <svg-icon icon-class="user" />
        </span>
        <input
          v-model="loginForm.account"
          placeholder="account"
          name="account"
          type="text"
          tabindex="1"
          auto-complete="on"
        />
      </el-form-item>

      <el-form-item prop="password" class="login-form-input">
        <span class="svg-container">
          <svg-icon icon-class="password" />
        </span>
        <input
          :key="passwordType"
          v-model="loginForm.password"
          :type="passwordType"
          placeholder="password"
          name="password"
          tabindex="2"
          auto-complete="on"
          @keyup.enter="handleLogin"
        />
        <span class="show-pwd" @click="showPwd">
          <svg-icon
            :icon-class="passwordType === 'password' ? 'eye' : 'eye-open'"
          />
        </span>
      </el-form-item>

      <el-button
        :loading="loading"
        type="primary"
        class="login-form-btn"
        style="width: 100%; margin-bottom: 30px"
        @click.prevent="handleLogin"
      >
        Sign in
      </el-button>
      <el-divider v-if="Object.keys(mediaLoginUrl).length > 0" class="divider">
        <span class="media-logo">
          <span
            v-for="(item, key) in mediaLoginUrl"
            :key="mediaMap[key].media"
            @click="handleMediaWindow(item)"
          >
            <svg-icon :icon-class="mediaMap[key].icon" class="icon" />
          </span>
        </span>
      </el-divider>
    </el-form>
  </div>
</template>
<script lang="ts">
export default { name: 'Login' }
</script>
<script lang="ts" setup>
import { param2Obj } from '@/utils'
import { validUsername, validPassword } from '@/utils/validate'
import { useRoute, useRouter } from 'vue-router'
import { useStore } from 'vuex'
import type { ElForm } from 'element-plus'
import { ref, watch, nextTick, reactive } from 'vue'
import { MediaLoginUrl } from '@/api/user'
const version = import.meta.env.VITE_APP_VERSION
const store = useStore()
const router = useRouter()
const form = ref<InstanceType<typeof ElForm>>()
const loginForm = ref({
  account: import.meta.env.PROD === true ? '' : 'admin',
  password: import.meta.env.PROD === true ? '' : 'admin!@#',
  phrase: '',
})
const loginRules: InstanceType<typeof ElForm>['rules'] = {
  account: [
    {
      required: true,
      trigger: 'blur',
      validator: (_, value) => {
        if (!validUsername(value)) {
          return new Error('Greater than 5 characters')
        } else {
          return true
        }
      },
    },
  ],
  password: [
    {
      required: true,
      trigger: 'blur',
      validator: (_, value) => {
        if (!validPassword(value)) {
          return new Error(
            '8 to 16 characters and a minimum of 2 character sets from these classes: [letters], [numbers], [special characters]'
          )
        } else {
          return true
        }
      },
    },
  ],
}

const redirectUri = window.location.origin + '/#/login'
const passwordType = ref('password')
const loading = ref(false)
const redirect = ref()
const query = ref()
watch(
  useRoute(),
  (route) => {
    redirect.value = route.query?.redirect
    if (route.query['account'] && route.query['time'] && route.query['token']) {
      handleExtLogin(
        route.query['account'] as string,
        Number(route.query['time']),
        route.query['token'] as string
      )
    }

    query.value = param2Obj(window.location.href)

    if (query.value.code) {
      handleMediaLogin(
        query.value.code.toString(),
        query.value.state.toString()
      )
    } else if (query.value.authCode) {
      handleMediaLogin(
        query.value.authCode.toString(),
        query.value.state.toString()
      )
    }
  },
  { immediate: true }
)

const mediaLoginUrl = ref<Record<string, string>>({})
const mediaMap = reactive<Record<string, any>>({
  dingtalk: {
    media: 'dingtalk',
    icon: 'dingtalk',
  },
  feishu: {
    media: 'feishu',
    icon: 'feishu',
  },
})
getMediaLoginUrl()

const password = ref<HTMLInputElement>()
function showPwd() {
  if (passwordType.value === 'password') {
    passwordType.value = ''
  } else {
    passwordType.value = 'password'
  }
  nextTick(() => {
    password.value?.focus()
  })
}

function handleExtLogin(account: string, time: number, token: string) {
  store
    .dispatch('user/extLogin', { account, time, token })
    .then(() => {
      router.push({
        path: redirect.value || '/',
        query: redirect.value ? param2Obj(redirect.value) : {},
      })
      loading.value = false
    })
    .catch(() => {
      loading.value = false
    })
}

function handleLogin() {
  form.value?.validate((valid) => {
    if (valid) {
      loading.value = true
      store
        .dispatch('user/login', loginForm.value)
        .then(() => {
          router.push({
            path: redirect.value || '/',
            query: redirect.value ? param2Obj(redirect.value) : {},
          })
          loading.value = false
        })
        .catch(() => {
          loading.value = false
        })
      return Promise.resolve(true)
    } else {
      return Promise.reject(false)
    }
  })
}

function handleMediaLogin(authCode: string, state: string) {
  loading.value = true
  store
    .dispatch('user/mediaLogin', { authCode, state, redirectUri })
    .then(() => {
      window.opener.location.reload()
      window.close()
    })
    .catch(() => {
      loading.value = false
    })
}

function getMediaLoginUrl() {
  new MediaLoginUrl({ redirectUri: redirectUri }).request().then((response) => {
    for (const media in response.data) {
      if (response.data[media] != '') {
        mediaLoginUrl.value[media] = response.data[media]
      }
    }
  })
}

function handleMediaWindow(url: string) {
  window.open(url, 'loginPopup', 'left=200, top=200, width=800,height=600')
}
</script>

<style lang="scss">
/* 修复input 背景不协调 和光标变色 */
/* Detail see https://github.com/PanJiaChen/vue-element-admin/pull/927 */

$bg: #fff;
$light_gray: #fff;
$cursor: #2f2f2f;

@supports (-webkit-mask: none) and (not (cater-color: $cursor)) {
  .login-container input {
    color: $cursor;
  }
}

/* reset element-ui css */
.login-container {
  input {
    width: 85%;
    background: $bg;
    border: none;
    padding: 12px 5px 12px 15px;
    color: #2f2f2f;
    height: 40px;
    caret-color: $cursor;

    &:-webkit-autofill {
      box-shadow: 0 0 0px 1000px $bg inset !important;
      -webkit-text-fill-color: $cursor !important;
    }
    &:focus {
      outline: none;
    }
  }
  .el-form-item__error {
    padding-top: 4px;
  }
  .el-form-item {
    background: #fff;
  }
}
</style>

<style lang="scss" scoped>
$bg: #fff;
$dark_gray: #889aa4;
$light_gray: #eee;

.login-container {
  display: flex;
  flex-direction: column;
  min-height: 100%;
  width: 100%;
  background-color: $bg;
  background-image: url('@/assets/images/background.jpg');
  background-position: left bottom;
  background-size: cover;
  width: 100%;
  overflow: hidden;
  .login-mark {
    margin-top: 58px;
  }

  .login-form {
    flex: 1;
    align-self: center;
    bottom: 0;
    width: 520px;
    padding: 0 35px;
    overflow: hidden;
    &-input {
      border: 1px solid rgba(229, 230, 231, 1);
      border-radius: 4px;
    }
    &-btn {
      margin-top: 15px;
      height: 40px;
      border-radius: 4px;
      font-size: 16px;
      background-color: #2580fa;
    }
  }

  .tips {
    font-size: 14px;
    color: #fff;
    margin-bottom: 10px;

    span {
      &:first-of-type {
        margin-right: 16px;
      }
    }
  }

  .svg-container {
    padding: 6px 5px 6px 20px;
    color: #c1c9d7;
    font-size: 16px;
    vertical-align: middle;
    width: 40px;
    display: inline-block;
  }

  .title-container {
    position: relative;
    padding-left: 5px;
    text-align: center;
    .title {
      margin: 20px 0;
      display: inline-block;
      font-size: 20px;
      font-family: PingFang SC;
      color: #586069;
    }
  }

  .show-pwd {
    position: absolute;
    right: 20px;
    top: 10px;
    font-size: 16px;
    color: #c1c9d7;
    cursor: pointer;
    user-select: none;
  }
}
.media-logo {
  display: flex;
  width: 100%;
  align-items: center;
  justify-content: center;
  .icon {
    width: 36px;
    height: 36px;
    margin: 0 5px;
    cursor: pointer;
  }
}

.divider {
  --el-bg-color: #ffffff;
  --el-border-color: #dcdfe6;
}
</style>
