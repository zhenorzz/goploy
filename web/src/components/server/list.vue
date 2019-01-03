<template>
  <el-row>
    <el-row type="flex" justify="end">
      <el-button type="primary" icon="el-icon-plus" @click="dialogFormVisible = true">添加</el-button>
    </el-row>
    <el-table :data="tableData" style="width: 100%">
      <el-table-column prop="name" label="服务器"></el-table-column>
      <el-table-column prop="ip" label="IP"></el-table-column>
      <el-table-column prop="path" label="部署目录" show-overflow-tooltip></el-table-column>
      <el-table-column prop="owner" label="sshKey所有者" show-overflow-tooltip></el-table-column>
      <el-table-column prop="createTime" label="创建时间" width="160"></el-table-column>
      <el-table-column prop="updateTime" label="更新时间" width="160"></el-table-column>
      <el-table-column prop="operation" label="操作" width="150">
        <template slot-scope="scope">
          <el-button size="small" type="primary">编辑</el-button>
          <el-button size="small" type="danger">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-dialog title="新增服务器" :visible.sync="dialogFormVisible">
      <el-form ref="form" :rules="form.rules" :model="form">
        <el-form-item label="服务器名称" label-width="120px" prop="name">
          <el-input v-model="form.name" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="IP" label-width="120px" prop="ip">
          <el-input v-model="form.ip" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="部署目录" label-width="120px" prop="path">
          <el-input v-model="form.path" autocomplete="off" placeholder="绝对路径"></el-input>
        </el-form-item>
        <el-form-item label="sshKey所有者" label-width="120px" prop="owner">
          <el-input v-model="form.owner" autocomplete="off" placeholder="root"></el-input>
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
import {get, add} from '@/api/server';
import {parseTime} from '@/utils/time';

export default {
  data() {
    return {
      dialogFormVisible: false,
      tableData: [],
      form: {
        disabled: false,
        name: '',
        ip: '',
        path: '',
        owner: '',
        rules: {
          name: [
            {required: true, message: '请输入服务器名称', trigger: 'blur'},
          ],
          ip: [
            {required: true, message: '请输入服务器ip', trigger: 'blur'},
          ],
          path: [
            {required: true, message: '请输入部署目录', trigger: 'blur'},
          ],
          owner: [
            {required: true, message: '请输入SSH-KEY的所有者', trigger: 'blur'},
          ],
        },
      },
    };
  },
  created() {
    get().then((response) => {
      const serverList = response.data.data.serverList;
      serverList.forEach((element) => {
        element.createTime = parseTime(element.createTime);
        element.updateTime = parseTime(element.updateTime);
      });
      this.tableData = serverList;
    });
  },
  methods: {
    add() {
      this.$refs.form.validate((valid) => {
        if (valid) {
          this.form.disabled = true;
          add(this.form.name, this.form.ip, this.form.path, this.form.owner).then((response) => {
            this.form.disabled = this.dialogFormVisible = false;
            this.$message({
              message: response.data.message,
              type: 'success',
              duration: 5 * 1000,
            });
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
