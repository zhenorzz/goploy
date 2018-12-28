<template>
  <el-row>
    <el-row type="flex" justify="end">
      <el-button type="primary" icon="el-icon-plus" @click="dialogFormVisible = true">添加</el-button>
    </el-row>
    <el-table :data="tableData" style="width: 100%">
      <el-table-column prop="project" label="项目名称"></el-table-column>
      <el-table-column prop="owner" label="仓库拥有者"></el-table-column>
      <el-table-column prop="repository" label="仓库名称"></el-table-column>
      <el-table-column prop="status" label="状态"></el-table-column>
      <el-table-column prop="createTime" label="创建时间"></el-table-column>
      <el-table-column prop="updateTime" label="更新时间"></el-table-column>
      <el-table-column prop="operation" label="操作" width="230">
        <template slot-scope="scope">
          <el-button size="small" type="success" @click="create(scope.row.id)">初始化</el-button>
          <el-button size="small" type="primary">编辑</el-button>
          <el-button size="small" type="danger">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-dialog title="提交项目" :visible.sync="dialogFormVisible">
      <el-form ref="form" :rules="form.rules" :model="form">
        <el-form-item label="项目名称" label-width="120px" prop="project">
          <el-input v-model="form.project" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="仓库拥有者" label-width="120px" prop="owner">
          <el-input v-model="form.owner" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="仓库名称" label-width="120px" prop="repository">
          <el-input v-model="form.repository" autocomplete="off"></el-input>
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
import {get, add, edit, remove, create} from '@/api/project';
import {parseTime} from '@/utils/time';

const STATUS = ['未初始化', '初始化中', '初始化成功', '初始化失败'];
export default {
  data() {
    return {
      dialogFormVisible: false,
      tableData: [],
      form: {
        repository: '',
        project: '',
        owner: '',
        disabled: false,
        rules: {
          project: [
            {required: true, message: '请输入项目名称', trigger: ['blur']},
          ],
          owner: [
            {required: true, message: '请输入仓库拥有者', trigger: ['blur']},
          ],
          repository: [
            {required: true, message: '请输入仓库名称', trigger: ['blur']},
          ],
        },
      },
    };
  },
  created() {
    this.get();
  },
  methods: {
    add() {
      this.$refs.form.validate((valid) => {
        if (valid) {
          this.form.disabled = true;
          add(this.form.project, this.form.owner, this.form.repository).then((response) => {
            this.form.disabled = false;
            this.dialogFormVisible = false;
            this.$message({
              message: response.data.message,
              type: 'success',
              duration: 5 * 1000,
            });
          }).catch(() => {
            this.form.disabled = false;
          });
        } else {
          return false;
        }
      });
    },
    get() {
      get().then((response) => {
        const projectList = response.data.data.projectList;
        projectList.forEach((element) => {
          element.createTime = parseTime(element.createTime);
          element.updateTime = parseTime(element.updateTime);
          element.status = STATUS[element.status];
        });
        this.tableData = projectList;
      }).catch(() => {
      });
    },
    create(projectId) {
      create(projectId).then((response) => {
        this.$message({
          message: response.data.message,
          type: 'success',
          duration: 5 * 1000,
        });
      });
    },
  },
};
</script>
