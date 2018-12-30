<template>
    <el-table
      :data="tableData"
      style="width: 100%">
      <el-table-column
        prop="project"
        label="项目名称">
      </el-table-column>
      <el-table-column
        prop="time"
        label="最近运行">
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
    </el-table>
  </template>

<script>
export default {
  data() {
    return {
      tableData: [{
        project: '测试项目',
        time: '2016-05-02',
        status: '成功',
        commit: 'xxx',
        deploy: 'xxx',
      }],
    };
  },
};
</script>
