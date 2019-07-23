<template>
  <div class="app-container">
    <el-form ref="form" :model="formData" :rules="formRules" :validate-on-rule-change="false" label-width="80px">
      <el-form-item label="角色" prop="roleId">
        <el-radio-group v-model="formData.roleId" @change="handleRoleChange">
          <el-radio v-for="role in roleList" v-show="role.id !== 1" :key="role.id" :label="role.id">{{ role.name }}</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="名称" prop="name">
        <el-input v-model="formData.name" autocomplete="off" style="width: 250px" />
      </el-form-item>
      <el-form-item label="描述" prop="remark">
        <el-input v-model="formData.remark" :rows="3" type="textarea" autocomplete="off" style="width: 250px" />
      </el-form-item>
      <el-form-item label="权限" prop="permissionList">
        <el-row v-for="(item, index) in permissionTree" :key="index">
          <el-row>{{ item.title }}</el-row>
          <el-row style="margin-left: 10px">
            <el-checkbox
              v-for="(children, cindex) in item.children"
              :key="cindex"
              v-model="formData.permissionList"
              :label="children"
            >
              {{ children.title }}
            </el-checkbox>
          </el-row>
        </el-row>
      </el-form-item>
      <el-form-item>
        <el-button :disabled="formProps.disabled" type="primary" @click="submitForm()">提交</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
import { getPermissionList, edit } from '@/api/role'
export default {
  name: 'Permission',
  data() {
    return {
      roleList: [],
      permissionTree: [],
      formProps: {
        disabled: false
      },
      formData: {
        roleId: 2,
        name: '',
        remark: '',
        permissionList: []
      },
      formRules: {
        roleId: [
          { required: true, message: '请选择所需要更新权限的角色', trigger: 'blur' }
        ],
        name: [
          { required: true, message: '请输入角色名称', trigger: 'blur' }
        ],
        permissionList: [
          { required: true, type: 'array', message: '请选择所需要更新的权限', trigger: 'blur' }
        ]
      }
    }
  },
  created() {
    this.getPermissionList()
  },
  methods: {
    handleRoleChange(roleId) {
      const role = this.roleList.find(element => element.id === roleId)
      this.formData.name = role['name']
      this.formData.remark = role['remark']
      const permissionList = role['permissionList'].split(',').map(element => parseInt(element))
      this.formData.permissionList = []
      if (permissionList.length === 0) {
        return
      }

      this.permissionTree.forEach(father => {
        father.children.forEach(child => {
          const childIndex = permissionList.indexOf(child.id)
          if (childIndex !== -1) {
            this.formData.permissionList.push(child)
          }
        })
      })
    },
    submitForm() {
      this.$refs.form.validate((valid) => {
        if (valid) {
          const permissionList = []
          this.formData.permissionList.forEach((item) => {
            permissionList.push(item.id, item.pid)
          })

          this.formProps.disabled = true
          edit(this.formData.roleId, this.formData.name, this.formData.remark, Array.from(new Set(permissionList)).sort((x, y) => { return x - y }).join(',')).then(response => {
            this.getPermissionList()
            this.$message({
              message: response.message,
              type: 'success',
              duration: 5 * 1000
            })
          }).finally(_ => {
            this.formProps.disabled = false
          })
        } else {
          return false
        }
      })
    },
    getPermissionList() {
      getPermissionList().then(response => {
        const data = response.data
        this.permissionTree = data.permissionTree
        this.roleList = data.roleList
        this.handleRoleChange(this.formData.roleId)
      })
    }
  }
}
</script>

<style rel="stylesheet/scss" lang="scss" scoped>
.el-form-item--mini.el-form-item {
    margin-bottom: 14px;
}

</style>
<style rel="stylesheet/scss" lang="scss">
.el-dialog__body {
    padding: 5px 20px;
}
</style>
