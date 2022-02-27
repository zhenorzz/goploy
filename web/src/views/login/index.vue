<template>
  <div class="login-container">
    <el-row class="login-mark" type="flex" align="middle" justify="center">
      <img src="@/assets/images/logo.png" width="120" height="120" />
    </el-row>
    <el-form
      ref="loginForm"
      :model="loginForm"
      :rules="loginRules"
      class="login-form"
      auto-complete="on"
      label-position="left"
    >
      <div class="title-container">
        <h3 class="title">Sign in to Goploy</h3>
      </div>

      <el-form-item prop="account" class="login-form-input">
        <span class="svg-container">
          <svg-icon icon-class="user" />
        </span>
        <el-input
          ref="account"
          v-model="loginForm.account"
          size="medium"
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
        <el-input
          :key="passwordType"
          ref="password"
          v-model="loginForm.password"
          :type="passwordType"
          size="medium"
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
        size="medium"
        type="primary"
        class="login-form-btn"
        style="width: 100%; margin-bottom: 30px"
        @click.prevent="handleLogin"
      >
        Sign in
      </el-button>
    </el-form>
  </div>
</template>

<script lang="ts">
import { param2Obj } from '@/utils'
import { validUsername, validPassword } from '@/utils/validate'
import Validator, { RuleItem } from 'async-validator'
import { defineComponent } from 'vue'
export default defineComponent({
  name: 'Login',
  data() {
    return {
      loginForm: {
        account: import.meta.env.PROD === true ? '' : 'admin',
        password: import.meta.env.PROD === true ? '' : 'admin!@#',
        phrase: '',
      },
      loginRules: {
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
          } as RuleItem,
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
          } as RuleItem,
        ],
      },
      loading: false,
      passwordType: 'password',
      redirect: '',
    }
  },
  watch: {
    $route: {
      handler: function (route) {
        this.redirect = route.query && route.query.redirect
      },
      immediate: true,
    },
  },
  methods: {
    showPwd() {
      if (this.passwordType === 'password') {
        this.passwordType = ''
      } else {
        this.passwordType = 'password'
      }
      this.$nextTick(() => {
        ;(this.$refs.password as HTMLInputElement).focus()
      })
    },
    handleLogin() {
      ;(this.$refs.loginForm as Validator).validate((valid: boolean) => {
        if (valid) {
          this.loading = true
          this.$store
            .dispatch('user/login', this.loginForm)
            .then(() => {
              this.$router.push({
                path: this.redirect || '/',
                query: this.redirect ? param2Obj(this.redirect) : {},
              })
              this.loading = false
            })
            .catch(() => {
              this.loading = false
            })
        } else {
          console.log('error submit!!')
          return false
        }
      })
    },
  },
})
</script>

<style lang="scss">
/* 修复input 背景不协调 和光标变色 */
/* Detail see https://github.com/PanJiaChen/vue-element-admin/pull/927 */

$bg: #fff;
$light_gray: #fff;
$cursor: #2f2f2f;

@supports (-webkit-mask: none) and (not (cater-color: $cursor)) {
  .login-container .el-input input {
    color: $cursor;
  }
}

/* reset element-ui css */
.login-container {
  .el-input {
    display: inline-block;
    height: 42px;
    width: 85%;

    input {
      background: $bg;
      border: 0px;
      -webkit-appearance: none;
      border-radius: 0px;
      padding: 12px 5px 12px 15px;
      color: #2f2f2f;
      height: 40px;
      caret-color: $cursor;

      &:-webkit-autofill {
        box-shadow: 0 0 0px 1000px $bg inset !important;
        -webkit-text-fill-color: $cursor !important;
      }
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
      border-radius: 35px;
    }
    &-btn {
      margin-top: 15px;
      height: 40px;
      border-radius: 39px;
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
</style>
