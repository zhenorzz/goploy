<template>
  <div class="app-container">
    <el-form
      ref="form"
      :rules="formRules"
      :model="formData"
      label-position="top"
      style="margin-left: 40px"
    >
      <el-form-item :label="$t('userPage.oldPassword')" prop="old">
        <el-input
          v-model="formData.old"
          :type="formProps.type.old"
          style="width: 300px"
        />
        <span class="show-pwd" @click="showPwd(inputElem.old)">
          <svg-icon icon-class="eye" />
        </span>
      </el-form-item>
      <el-form-item :label="$t('userPage.newPassword')" prop="new">
        <el-input
          v-model="formData.new"
          :type="formProps.type.new"
          style="width: 300px"
          autocomplete="off"
        />
        <span class="show-pwd" @click="showPwd(inputElem.new)">
          <svg-icon icon-class="eye" />
        </span>
      </el-form-item>
      <el-form-item :label="$t('userPage.rePassword')" prop="confirm">
        <el-input
          v-model="formData.confirm"
          :type="formProps.type.confirm"
          style="width: 300px"
          autocomplete="off"
        />
        <span class="show-pwd" @click="showPwd(inputElem.confirm)">
          <svg-icon icon-class="eye" />
        </span>
      </el-form-item>
      <el-form-item>
        <el-button
          :loading="formProps.loading"
          type="primary"
          @click="changePassword()"
        >
          {{ $t('save') }}
        </el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script lang="ts">
export default { name: 'UserProfile' }
</script>
<script lang="ts" setup>
import type { ElForm } from 'element-plus'
import { ElMessage } from 'element-plus'
import { validPassword } from '@/utils/validate'
import { UserChangePassword } from '@/api/user'
import { ref } from 'vue'
enum inputElem {
  old = 'old',
  new = 'new',
  confirm = 'confirm',
}

const form = ref<InstanceType<typeof ElForm>>()
const formData = ref({
  old: '',
  new: '',
  confirm: '',
})
const formProps = ref({
  loading: false,
  type: {
    old: 'password',
    new: 'password',
    confirm: 'password',
  },
})

const formRules = <InstanceType<typeof ElForm>['rules']>{
  old: [
    {
      required: true,
      message: 'Old password required',
      trigger: ['blur'],
    },
  ],
  new: [
    {
      required: true,
      message:
        '8 to 16 characters and a minimum of 2 character sets from these classes: [letters], [numbers], [special characters]',
      trigger: ['blur'],
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
  confirm: [
    {
      required: true,
      validator: (_, value) => {
        if (value === '') {
          return new Error('Please enter the password again')
        } else if (value !== formData.value.new) {
          return new Error('The two passwords do not match!')
        } else {
          return true
        }
      },
      trigger: ['blur'],
    },
  ],
}

function showPwd(index: inputElem) {
  if (formProps.value.type[index] === 'password') {
    formProps.value.type[index] = ''
  } else {
    formProps.value.type[index] = 'password'
  }
}
function changePassword() {
  form.value?.validate((valid) => {
    if (valid) {
      formProps.value.loading = true
      new UserChangePassword({
        oldPwd: formData.value.old,
        newPwd: formData.value.new,
      })
        .request()
        .then(() => {
          formProps.value.loading = false
          ElMessage.success('Success')
        })
        .catch(() => {
          formProps.value.loading = false
        })
      return Promise.resolve(true)
    } else {
      return Promise.reject(false)
    }
  })
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
