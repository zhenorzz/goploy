<template>
  <div class="login-container">
    <el-row class="login-mark" type="flex" align="middle">
      <img src="@/assets/images/logo.png" width="60" height="60">
      <el-row style="margin-left: 10px;margin-top:5px;">
        <el-row class="main-mark">代码部署平台</el-row>
        <el-row class="sub-mark">Goploy</el-row>
      </el-row>
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
        <h3 class="title">登录</h3>
        <h4 class="sub-title">SIGN IN</h4>
      </div>

      <el-form-item prop="account" class="login-form-input">
        <span class="svg-container">
          <svg-icon icon-class="user" />
        </span>
        <el-input
          ref="account"
          v-model="loginForm.account"
          size="medium"
          placeholder="请输入账号"
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
          placeholder="请输入密码"
          name="password"
          tabindex="2"
          auto-complete="on"
          @keyup.enter.native="handleLogin"
        />
        <span class="show-pwd" @click="showPwd">
          <svg-icon :icon-class="passwordType === 'password' ? 'eye' : 'eye-open'" />
        </span>
      </el-form-item>

      <el-button
        :loading="loading"
        size="medium"
        type="primary"
        class="login-form-btn"
        style="width:100%;margin-bottom:30px;"
        @click.native.prevent="handleLogin"
      >登 录</el-button>
    </el-form>
  </div>
</template>

<script>
import { validUsername, validPassword } from '@/utils/validate'
export default {
  name: 'Login',
  data() {
    const validateUsername = (rule, value, callback) => {
      if (!validUsername(value)) {
        callback(new Error('请输入正确的用户名'))
      } else {
        callback()
      }
    }
    const validatePassword = (rule, value, callback) => {
      if (!validPassword(value)) {
        callback(new Error('8到16个字符，至少包含字母、数字、特殊符号中的两种'))
      } else {
        callback()
      }
    }
    return {
      loginForm: {
        account: process.env.NODE_ENV === 'production' ? '' : 'admin',
        password: process.env.NODE_ENV === 'production' ? '' : 'admin!@#',
        phrase: ''
      },
      loginRules: {
        account: [{ required: true, trigger: 'blur', validator: validateUsername }],
        password: [{ required: true, trigger: 'blur', validator: validatePassword }]
      },
      loading: false,
      passwordType: 'password',
      redirect: undefined
    }
  },
  watch: {
    $route: {
      handler: function(route) {
        this.redirect = route.query && route.query.redirect
      },
      immediate: true
    }
  },
  created() {
  },
  methods: {
    showPwd() {
      if (this.passwordType === 'password') {
        this.passwordType = ''
      } else {
        this.passwordType = 'password'
      }
      this.$nextTick(() => {
        this.$refs.password.focus()
      })
    },
    changePhrase() {
      this.phrase = process.env.VUE_APP_BASE_API + '/user/getPhrase?' + Math.random()
    },
    handleLogin() {
      this.$refs.loginForm.validate(valid => {
        if (valid) {
          this.loading = true
          this.$store.dispatch('user/login', this.loginForm).then(() => {
            this.$router.push({ path: this.redirect || '/' })
            this.loading = false
          }).catch(() => {
            this.loading = false
            this.showPhrase = true
            this.changePhrase()
          })
        } else {
          console.log('error submit!!')
          return false
        }
      })
    }
  }
}
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
  flex-direction:column;
  min-height: 100%;
  width: 100%;
  background-color: $bg;
  background-image: url("~@/assets/images/background.jpg");
  background-position: left bottom;
  background-size: cover;
	width: 100%;
  overflow: hidden;
  .login-mark {
    margin-top: 58px;
    margin-left: 87px;
    .main-mark {
      height: 30px;
      font-size: 30px;
      font-family: PingFang SC;
      font-weight: 500;
      letter-spacing: 8px;
      color: rgba(105, 116, 139, 1);
      line-height: 30px;
    }
    .sub-mark {
      margin-top: 5px;
      margin-left: 2px;
      height: 23px;
      font-size: 12px;
      font-family: PingFang SC;
      font-weight: 400;
      letter-spacing: 4px;
      color: rgba(193, 201, 215, 1);
      line-height: 23px;
      opacity: 0.81;
    }
  }

  .login-form {
    flex: 1;
    align-self: center;
    bottom: 0;
    width: 520px;
    padding: 50px 35px;
    overflow: hidden;
    &-input {
      border: 1px solid rgba(229, 230, 231, 1);
      border-radius: 35px;
    }
    &-btn {
       margin-top:15px;
      height: 40px;
      border-radius:39px;
      font-size: 16px;
      background-color: #2580FA;
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

    .title {
      margin: 20px 0;
      display: inline-block;
      width: 80px;
      font-size: 36px;
      font-family: PingFang SC;
      font-weight: 500;
      color: rgba(47, 47, 47, 1);
    }
    .sub-title {
      margin: 20px 0 20px 5px;
      display: inline-block;
      font-size: 14px;
      font-family: PingFang SC;
      font-weight: 500;
      color: rgba(206, 212, 223, 1);
      opacity: 0.82;
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
