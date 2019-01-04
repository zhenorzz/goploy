<template>
  <el-row>
    <el-form ref="pwdForm" :model="form" :rules="form.rules" label-width="100px" style="margin-left:50px">
      <el-form-item label="" prop="old">
        <el-input
          v-model="form.projectId"
          style="width: 300px"
          placeholder="请输入原密码"
        />
      </el-form-item>
      <el-form-item label="" prop="old">
        <el-input
          v-model="form.owner"
          style="width: 300px"
          placeholder="请输入原密码"
        />
      </el-form-item>
      <el-form-item label="新密码" prop="new">
        <el-input
          v-model="form.repository"
          style="width: 300px"
          placeholder="请输入新密码"
        />
      </el-form-item>
      <el-form-item label="确认新密码" prop="confirm">
        <el-input
          v-model="form.status"
          style="width: 300px"
          placeholder="确认新密码"
        />
      </el-form-item>
      <el-form-item label="确认新密码" prop="confirm">
        <el-input
          v-model="form.createTime"
          style="width: 300px"
          placeholder="确认新密码"
        />
      </el-form-item>
      <el-form-item label="确认新密码" prop="confirm">
        <el-input
          v-model="form.updateTime"
          style="width: 300px"
          placeholder="确认新密码"
        />
      </el-form-item>
      <el-form-item>
        <el-button :disabled="form.disabled" type="primary">保存</el-button>
      </el-form-item>
    </el-form>
  </el-row>
</template>
<script>

import {getDetail} from '@/api/project';

const STATUS = ['未初始化', '初始化中', '初始化成功', '初始化失败'];
export default {
  data() {
    return {
      projectId: this.$route.query.project_id,
      serverOption: [],
      userOption: [],
      form: {
        disabled: false,
        projectId: this.$route.query.project_id,
        project: '',
        owner: '',
        repository: '',
        status: '',
        createTime: '',
        updateTime: '',
        serverIds: [],
        userIds: [],
        rules: {

        },
      },
    };
  },
  created() {
  },
  methods: {
  },
};
</script>
