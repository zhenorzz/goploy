<template>
  <el-row>
    <el-row type="flex" justify="end">
      <el-button type="primary" icon="el-icon-plus" @click="dialogFormVisible = true">添加</el-button>
    </el-row>
    <el-table :data="tableData" style="width: 100%">
      <el-table-column prop="account" label="账号"></el-table-column>
      <el-table-column prop="name" label="名称"></el-table-column>
      <el-table-column prop="email" label="邮箱" show-overflow-tooltip></el-table-column>
      <el-table-column prop="role" label="角色" show-overflow-tooltip></el-table-column>
      <el-table-column prop="createTime" label="创建时间" width="160"></el-table-column>
      <el-table-column prop="updateTime" label="更新时间" width="160"></el-table-column>
      <el-table-column prop="operation" label="操作" width="210">
        <template slot-scope="scope">
          <el-button size="small" type="primary">编辑</el-button>
          <el-button size="small" type="danger">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-row type="flex" justify="end" style="margin-top: 10px;">
      <el-pagination
        v-show="pagination.total>pagination.rows"
        :total="pagination.total"
        :page-size="pagination.rows"
        background
        layout="prev, pager, next"
        @current-change="handleCurrentChange"
      />
    </el-row>
    <el-dialog title="新增成员" :visible.sync="dialogFormVisible">
      <el-form ref="form" :rules="form.rules" :model="form">
        <el-form-item label="账号" label-width="120px" prop="account">
          <el-input v-model="form.account" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="密码" label-width="120px" prop="password">
          <el-input v-model="form.password" autocomplete="off" placeholder="请输入初始密码"></el-input>
        </el-form-item>
        <el-form-item label="名称" label-width="120px" prop="name">
          <el-input v-model="form.name" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="邮箱" label-width="120px" prop="email">
          <el-input v-model="form.email" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="角色" label-width="120px" prop="role">
          <el-select v-model="form.role" placeholder="请选择角色">
            <el-option label="管理员" value="manager"></el-option>
            <el-option label="普通成员" value="member"></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">取 消</el-button>
        <el-button :disabled="form.disabled" type="primary" @click="add">确 定</el-button>
      </div>
    </el-dialog>
  </el-row>
</template>
<script>
import {get, add} from '@/api/user';
import {parseTime} from '@/utils/time';

export default {
  data() {
    return {
      dialogFormVisible: false,
      tableData: [],
      pagination: {
        page: 1,
        rows: 11,
        total: 0,
      },
      form: {
        disabled: false,
        account: '',
        password: '',
        name: '',
        email: '',
        role: '',
        rules: {
          account: [
            {required: true, message: '请输入账号', trigger: 'blur'},
          ],
          password: [
            {required: true, message: '请输入初始密码', trigger: 'blur'},
          ],
          name: [
            {required: true, message: '请输入名称', trigger: 'blur'},
          ],
          email: [
            {type: 'email', message: '邮箱格式不正确', trigger: 'blur'},
          ],
          role: [
            {required: true, message: '请选择角色', trigger: 'change'},
          ],
        },
      },
    };
  },
  created() {
    this.getUserList();
  },
  methods: {
    getUserList() {
      get(this.pagination).then((response) => {
        const userList = response.data.data.userList;
        userList.forEach((element) => {
          element.createTime = parseTime(element.createTime);
          element.updateTime = parseTime(element.updateTime);
        });
        this.tableData = userList;
        this.pagination = response.data.data.pagination;
      });
    },
    // 分页事件
    handleCurrentChange(val) {
      this.pagination.page = val;
      this.getUserList();
    },
    add() {
      this.$refs.form.validate((valid) => {
        if (valid) {
          this.form.disabled = true;
          add(this.form.account, this.form.password, this.form.name, this.form.email, this.form.role).then((response) => {
            this.form.disabled = this.dialogFormVisible = false;
            this.$message({
              message: response.data.message,
              type: 'success',
              duration: 5 * 1000,
            });
            this.getUserList();
          }).catch(() => {
            this.form.disabled = this.dialogFormVisible = false;
          });
        } else {
          this.form.disabled = this.dialogFormVisible = false;
          return false;
        }
      });
    },
  },
};
</script>
