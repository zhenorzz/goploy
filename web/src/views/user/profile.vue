<template>
  <div class="app-container">
    <el-form ref="pwdForm" :model="pwdForm" :rules="pwdForm.rules" label-width="100px">
      <el-form-item :label="$t('userPage.oldPassword')" prop="old">
        <el-input v-model="pwdForm.old" :type="pwdForm.type.old" style="width: 300px" />
        <span class="show-pwd" @click="showPwd('old')">
          <svg-icon icon-class="eye" />
        </span>
      </el-form-item>
      <el-form-item :label="$t('userPage.newPassword')" prop="new">
        <el-input v-model="pwdForm.new" :type="pwdForm.type.new" style="width: 300px" />
        <span class="show-pwd" @click="showPwd('new')">
          <svg-icon icon-class="eye" />
        </span>
      </el-form-item>
      <el-form-item :label="$t('userPage.rePassword')" prop="confirm">
        <el-input v-model="pwdForm.confirm" :type="pwdForm.type.confirm" style="width: 300px" />
        <span class="show-pwd" @click="showPwd('confirm')">
          <svg-icon icon-class="eye" />
        </span>
      </el-form-item>
      <el-form-item>
        <el-button :loading="pwdForm.loading" type="primary" @click="changePassword()">{{ $t('save') }}</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
import { validPassword } from '@/utils/validate'
import { changePassword } from '@/api/user'

export default {
  name: 'UserProfile',
  data() {
    const confirmPass = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('Please enter the password again'))
      } else if (value !== this.pwdForm.new) {
        callback(new Error('The two passwords do not match!'))
      } else {
        callback()
      }
    }
    const validatePass = (rule, value, callback) => {
      if (!validPassword(value)) {
        callback(new Error('8 to 16 characters and a minimum of 2 character sets from these classes: [letters], [numbers], [special characters]'))
      } else {
        callback()
      }
    }
    return {
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
            { required: true, message: 'Old password required', trigger: ['blur'] }
          ],
          new: [
            { required: true, message: '8 to 16 characters and a minimum of 2 character sets from these classes: [letters], [numbers], [special characters]', trigger: ['blur'], validator: validatePass }
          ],
          confirm: [
            { required: true, validator: confirmPass, trigger: ['blur'] }
          ]
        }
      }
    }
  },
  methods: {
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
            this.$message.success('Success')
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

