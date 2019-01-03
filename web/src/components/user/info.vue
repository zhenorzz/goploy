<template>
  <el-row>
    <h1 class="app-title">密码设置</h1>
    <el-form ref="pwdForm" :model="pwdForm" :rules="pwdForm.rules" label-width="100px">
      <el-form-item label="原密码" prop="old">
        <el-input v-model="pwdForm.old" :type="pwdForm.type.old" style="width: 300px" placeholder="请输入原密码"/>
        <span class="show-pwd" @click="showPwd('old')">
          <svg-icon icon-class="eye" />
        </span>
      </el-form-item>
      <el-form-item label="新密码" prop="new">
        <el-input v-model="pwdForm.new" :type="pwdForm.type.new" style="width: 300px" placeholder="请输入新密码"/>
        <span class="show-pwd" @click="showPwd('new')">
          <svg-icon icon-class="eye" />
        </span>
      </el-form-item>
      <el-form-item label="确认新密码" prop="confirm">
        <el-input v-model="pwdForm.confirm" :type="pwdForm.type.confirm" style="width: 300px" placeholder="确认新密码"/>
        <span class="show-pwd" @click="showPwd('confirm')">
          <svg-icon icon-class="eye" />
        </span>
      </el-form-item>
      <el-form-item>
        <el-button :loading="pwdForm.loading" type="primary" @click="changePassword">保存</el-button>
      </el-form-item>
    </el-form>
  </el-row>
</template>
<script>

import {changePassword} from '@/api/user';

export default {
  data() {
    const confirmPass = (rule, value, callback) => {
      if (value === '') {
        callback(new Error('请再次输入密码'));
      } else if (value !== this.pwdForm.new) {
        callback(new Error('两次输入密码不一致!'));
      } else {
        callback();
      }
    };
    return {
      pwdForm: {
        old: '',
        new: '',
        confirm: '',
        loading: false,
        type: {
          old: 'password',
          new: 'password',
          confirm: 'password',
        },
        rules: {
          old: [
            {required: true, min: 5, message: '不能少于5位数', trigger: ['blur']},
          ],
          new: [
            {required: true, min: 5, message: '不能少于5位数', trigger: ['blur']},
          ],
          confirm: [
            {required: true, validator: confirmPass, trigger: ['blur']},
          ],
        },
      },
    };
  },
  methods: {
    showPwd(index) {
      if (this.pwdForm.type[index] === 'password') {
        this.pwdForm.type[index] = '';
      } else {
        this.pwdForm.type[index] = 'password';
      }
    },
    changePassword() {
      this.$refs.pwdForm.validate((valid) => {
        if (valid) {
          this.pwdForm.loading = true;
          changePassword(this.$store.getters.account, this.pwdForm.old, this.pwdForm.new).then((response) => {
            this.pwdForm.loading = false;
            this.$message({
              message: response.data.message,
              type: 'success',
              duration: 5 * 1000,
            });
          }).catch(() => {
            this.pwdForm.loading = false;
          });
        } else {
          console.log('error submit!!');
          return false;
        }
      });
    },
  },
};
</script>
<style rel="stylesheet/scss" lang="scss" scoped>
.show-pwd {
  position: absolute;
  left: 270px;
  color: #8c939d;
  cursor: pointer;
}
</style>
