<template>
  <el-row>
    <el-row type="flex" justify="end">
      <el-button type="primary" icon="el-icon-plus" @click="dialogFormVisible = true">添加</el-button>
    </el-row>
    <el-table :data="tableData" style="width: 100%">
      <el-table-column prop="project" label="项目名称"></el-table-column>
      <el-table-column prop="branch" label="分支"></el-table-column>
      <el-table-column prop="commit" label="提交信息" show-overflow-tooltip></el-table-column>
      <el-table-column prop="commitSha" label="sha" show-overflow-tooltip></el-table-column>
      <el-table-column prop="type" label="上线"></el-table-column>
      <el-table-column prop="status" label="状态"></el-table-column>
      <el-table-column prop="creator" label="提交于"></el-table-column>
      <el-table-column prop="editor" label="发布于"></el-table-column>
      <el-table-column prop="createTime" label="创建时间" width="160"></el-table-column>
      <el-table-column prop="updateTime" label="更新时间" width="160"></el-table-column>
      <el-table-column prop="operation" label="操作" width="210">
        <template slot-scope="scope">
          <el-button size="small" type="primary">编辑</el-button>
          <el-button size="small" type="success">发布</el-button>
          <el-button size="small" type="danger">回滚</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-dialog title="提交项目" :visible.sync="dialogFormVisible">
      <el-form ref="form" :rules="form.rules" :model="form">
        <el-form-item label="项目仓库" label-width="120px" prop="projectId">
          <el-select v-model="form.projectId" placeholder="选择项目仓库" @change="selectBranch">
            <el-option
              v-for="(item, index) in projectOption"
              :key="index"
              :label="item.project"
              :value="item.id"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="分支" label-width="120px" prop="branch">
          <el-select v-model="form.branch" placeholder="选择分支" @change="selectCommit">
            <el-option
              v-for="(item, index) in branchOption"
              :key="index"
              :label="item.name"
              :value="item.name"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="提交" label-width="120px" prop="commitSha">
          <el-select v-model="form.commitSha" placeholder="选择Commit">
            <el-option
              v-for="(item, index) in commitOption"
              :key="index"
              :label="item.commit.committer.name + ' : ' + item.commit.message + ' : ' + item.sha.substring(0, 10)"
              :value="item.sha"
              @click.native="selectLabel(item.commit.committer.name + ' : ' + item.commit.message + ' : ' + item.sha.substring(0, 10))"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="部署文件" label-width="120px" prop="type">
          <el-radio v-model="form.type" :label="1">全量上线</el-radio>
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
import {get as getProject, branch, commit} from '@/api/project';
import {get as getDeploy, add} from '@/api/deploy';
import {parseTime} from '@/utils/time';

const TYPE = ['', '全量上线'];
const STATUS = ['提交', '审核', '上线'];
export default {
  data() {
    return {
      dialogFormVisible: false,
      tableData: [],
      projectOption: [],
      branchOption: [],
      commitOption: [],
      form: {
        disabled: false,
        projectId: '',
        branch: '',
        commit: '',
        commitSha: '',
        type: 1,
        rules: {
          projectId: [
            {required: true, message: '请选择项目', trigger: 'change'},
          ],
          branch: [
            {required: true, message: '请选择分支', trigger: 'change'},
          ],
          commitSha: [
            {required: true, message: '请选择提交信息', trigger: 'change'},
          ],
          type: [
            {required: true, message: '请选择类型', trigger: 'change'},
          ],
        },
      },
    };
  },
  created() {
    getProject().then((response) => {
      this.projectOption = response.data.data.projectList;
    });
    getDeploy().then((response) => {
      const deployList = response.data.data.deployList;
      deployList.forEach((element) => {
        element.createTime = parseTime(element.createTime);
        element.updateTime = parseTime(element.updateTime);
        element.type = TYPE[element.type];
        element.status = STATUS[element.status];
      });
      this.tableData = deployList;
    });
  },
  methods: {
    selectBranch() {
      branch().then((response) => {
        this.branchOption = response.data.data.branchList;
      });
    },
    selectCommit() {
      commit().then((response) => {
        this.commitOption = response.data.data.commitList;
      });
    },
    selectLabel(label) {
      this.form.commit = label;
    },
    add() {
      this.$refs.form.validate((valid) => {
        if (valid) {
          this.form.disabled = true;
          add(this.form.projectId, this.form.branch, this.form.commit, this.form.commitSha, this.form.type).then((response) => {
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
