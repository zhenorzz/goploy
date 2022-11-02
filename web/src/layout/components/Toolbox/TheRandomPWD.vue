<template>
  <el-row v-show="modelValue === 'password'">
    <el-checkbox-group v-model="password.checkbox">
      <el-checkbox :label="1">A-Z</el-checkbox>
      <el-checkbox :label="2">a-z</el-checkbox>
      <el-checkbox :label="3">0-9</el-checkbox>
      <el-checkbox :label="4">!@#$%^&*</el-checkbox>
    </el-checkbox-group>
    <el-row style="margin-top: 10px" type="flex" align="middle">
      <span style="width: 60px; font-size: 14px; margin-right: 5px">
        Length
      </span>
      <el-input-number
        v-model="password.length"
        :min="1"
        :max="40"
        placeholder="Please enter the password length"
      />
      <el-button type="primary" @click="createPassword">Gen</el-button>
    </el-row>

    <el-input :value="password.text" style="margin-top: 10px" readonly />
  </el-row>
</template>

<script lang="ts" setup>
import { reactive } from 'vue'
defineProps({
  modelValue: {
    type: String,
    default: '',
  },
})
const password = reactive({
  checkbox: [1, 2, 3],
  length: 8,
  text: '',
})
const createPassword = () => {
  let randArr: string[] = []
  for (const num of password.checkbox) {
    if (num === 1) {
      for (let i = 0; i < 26; i++) {
        randArr.push(String.fromCharCode(65 + i))
      }
    } else if (num === 2) {
      for (let i = 0; i < 26; i++) {
        randArr.push(String.fromCharCode(97 + i))
      }
    } else if (num === 3) {
      for (let i = 0; i < 10; i++) {
        randArr.push(String(i))
      }
    } else {
      randArr = randArr.concat(['!', '@', '#', '$', '%', '^', '&', '*'])
    }
  }
  let tmpPWD = ''
  for (let i = 0; i < password.length; i++) {
    tmpPWD += randArr[Math.floor(Math.random() * randArr.length)]
  }
  password.text = tmpPWD
}
</script>
