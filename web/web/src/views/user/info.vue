<template>
  <div class="app-container">
    <el-form ref="userInfoForm" :model="userInfoForm" :rules="userInfoForm.rules" label-width="100px">
      <el-form-item label="账号">
        <el-input v-model="userInfoForm.data.account" readonly="readonly" disabled="disabled" style="width: 300px" />
      </el-form-item>
      <el-form-item label="名称">
        <el-input v-model="userInfoForm.data.name" style="width: 300px" placeholder="请输入名称" />
      </el-form-item>
      <el-form-item>
        <el-button :loading="userInfoForm.prop.loading" type="primary" @click="changeUserInfo()">保存</el-button>
      </el-form-item>
    </el-form>
    <el-form ref="pwdForm" :model="pwdForm" :rules="pwdForm.rules" label-width="100px">
      <el-form-item label="原密码" prop="old">
        <el-input v-model="pwdForm.old" :type="pwdForm.type.old" style="width: 300px" placeholder="请输入原密码" />
        <span class="show-pwd" @click="showPwd('old')">
          <svg-icon icon-class="eye" />
        </span>
      </el-form-item>
      <el-form-item label="新密码" prop="new">
        <el-input v-model="pwdForm.new" :type="pwdForm.type.new" style="width: 300px" placeholder="请输入新密码" />
        <span class="show-pwd" @click="showPwd('new')">
          <svg-icon icon-class="eye" />
        </span>
      </el-form-item>
      <el-form-item label="确认新密码" prop="confirm">
        <el-input v-model="pwdForm.confirm" :type="pwdForm.type.confirm" style="width: 300px" placeholder="确认新密码" />
        <span class="show-pwd" @click="showPwd('confirm')">
          <svg-icon icon-class="eye" />
        </span>
      </el-form-item>
      <el-form-item>
        <el-button :loading="pwdForm.loading" type="primary" @click="changePassword()">保存</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
import { validPassword } from '@/utils/validate'
import { changePassword } from '@/api/user'

export default {
  name: 'UserInfo',
  data() {
    const confirmPass = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请再次输入密码'))
      } else if (value !== this.pwdForm.new) {
        callback(new Error('两次输入密码不一致!'))
      } else {
        callback()
      }
    }
    const validatePass = (rule, value, callback) => {
      if (!validPassword(value)) {
        callback(new Error('8到16个字符，至少包含字母、数字、特殊符号中的两种'))
      } else {
        callback()
      }
    }
    return {
      userInfoForm: {
        prop: {
          loading: false
        },
        data: {
          account: this.$store.getters.account,
          name: this.$store.getters.name
        },
        rules: {
        }
      },
      pwdForm: {
        old: '',
        new: '',
        confirm: '',
        loading: false,
        type: {
          old: 'password',
          new: 'password',
          confirm: 'password'
        },
        rules: {
          old: [
            { required: true, message: '请输入旧密码', trigger: ['blur'] }
          ],
          new: [
            { required: true, message: '请输入8到16个字符，至少包含字母、数字、特殊符号中的两种', trigger: ['blur'], validator: validatePass }
          ],
          confirm: [
            { required: true, validator: confirmPass, trigger: ['blur'] }
          ]
        }
      }
    }
  },
  methods: {
    changeUserInfo() {
      this.$refs.userInfoForm.validate((valid) => {
        if (valid) {
          this.userInfoForm.prop.loading = true
          this.$store.dispatch('ChangeInfo', this.userInfoForm.data).then((response) => {
            this.$message({
              message: response.message,
              type: 'success',
              duration: 5 * 1000
            })
          }).finally(() => {
            this.userInfoForm.prop.loading = false
          })
        } else {
          console.log('error submit!!')
          return false
        }
      })
    },
    showPwd(index) {
      if (this.pwdForm.type[index] === 'password') {
        this.pwdForm.type[index] = ''
      } else {
        this.pwdForm.type[index] = 'password'
      }
    },
    changePassword() {
      this.$refs.pwdForm.validate((valid) => {
        if (valid) {
          this.pwdForm.loading = true
          changePassword(this.pwdForm.old, this.pwdForm.new).then(response => {
            this.pwdForm.loading = false
            this.$message({
              message: response.message,
              type: 'success',
              duration: 5 * 1000
            })
          }).catch(() => {
            this.pwdForm.loading = false
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

<style rel="stylesheet/scss" lang="scss" scoped>
.show-pwd {
  position: absolute;
  left: 270px;
  color: #8c939d;
  cursor: pointer;
}
</style>

