<template>
  <el-row>
    <el-row type="flex" justify="end">
      <el-button type="primary" icon="el-icon-plus" @click="dialogFormVisible = true">添加</el-button>
    </el-row>
    <el-table
      :data="tableData"
      style="width: 100%">
      <el-table-column
        prop="project"
        label="项目名称">
      </el-table-column>
      <el-table-column
        prop="time"
        label="更新时间">
      </el-table-column>
      <el-table-column
        prop="status"
        label="状态">
      </el-table-column>
      <el-table-column
        prop="commit"
        label="提交于">
      </el-table-column>
      <el-table-column
        prop="deploy"
        label="发布于">
      </el-table-column>
      <el-table-column
          prop="operation"
          label="操作"
          width="210">
          <template slot-scope="scope">
            <el-button size="small" type="primary">编辑</el-button>
            <el-button size="small" type="success">发布</el-button>
            <el-button size="small" type="danger">回滚</el-button>
          </template>
        </el-table-column>
    </el-table>
    <el-dialog title="提交项目" :visible.sync="dialogFormVisible">
      <el-form :model="form">
        <el-form-item label="项目仓库" label-width="120px">
          <el-select v-model="form.repository" placeholder="选择项目仓库">
            <el-option label="区域一" value="shanghai"></el-option>
            <el-option label="区域二" value="beijing"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="分支" label-width="120px">
          <el-select v-model="form.branch" placeholder="选择分支">
            <el-option label="区域一" value="shanghai"></el-option>
            <el-option label="区域二" value="beijing"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="部署文件" label-width="120px">
          <el-select v-model="form.file" placeholder="请选择部署文件">
            <el-option label="区域一" value="shanghai"></el-option>
            <el-option label="区域二" value="beijing"></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">取 消</el-button>
        <el-button type="primary" @click="dialogFormVisible = false">确 定</el-button>
      </div>
    </el-dialog>
  </el-row>
</template>
<script>
export default {
  data() {
    return {
      dialogFormVisible: false,
      tableData: [{
        project: '测试项目',
        time: '2016-05-02',
        status: '成功',
        commit: 'xxx',
        deploy: 'xxx',
      }],
      form: {
        repository: '',
        branch: '',
        file: '',
      },
    };
  },
};
</script>
