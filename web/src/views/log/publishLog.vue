<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex" justify="space-between">
      <el-row>
        <el-input
          v-model="searchParam.username"
          style="width: 200px"
          placeholder="Filter the username"
        />
        <el-input
          v-model="searchParam.projectName"
          style="width: 200px"
          placeholder="Filter the project name"
        />
        <el-button
          :loading="tableLoading"
          type="primary"
          :icon="Search"
          @click="searchList"
        />
      </el-row>
    </el-row>
    <el-row class="app-table">
      <el-table
        v-loading="tableLoading"
        highlight-current-row
        height="100%"
        :data="tableData"
      >
        <el-table-column prop="token" label="Token" width="300" />
        <el-table-column prop="publisherName" label="Username" width="100" />
        <el-table-column prop="projectName" label="Project Name" width="160" />
        <el-table-column prop="state" label="State" align="center" width="80">
          <template #default="scope">
            <span v-if="scope.row.state === 1" style="color: #67c23a">
              {{ $t('success') }}
            </span>
            <span v-else style="color: #f56c6c">{{ $t('fail') }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="detail" label="Reason" show-overflow-tooltip />
        <el-table-column prop="insertTime" label="insertTime" width="160" />
        <el-table-column
          prop="operation"
          :label="$t('op')"
          width="100"
          align="center"
          :fixed="$store.state.app.device === 'mobile' ? false : 'right'"
        >
          <template #default="scope">
            <Button
              type="primary"
              link
              :permissions="[pms.DeployDetail]"
              @click="handleDetail(scope.row)"
            >
              {{ $t('detail') }}
            </Button>
          </template>
        </el-table-column>
      </el-table>
    </el-row>
    <el-row type="flex" justify="end" class="app-page">
      <el-pagination
        :total="pagination.total"
        :page-size="pagination.rows"
        background
        layout="total, prev, pager, next"
        @current-change="handlePageChange"
      />
    </el-row>
    <TheDetailDialog
      v-model="dialogVisible"
      :show-preivew="false"
      :project-row="projectItem"
    />
  </el-row>
</template>

<script lang="ts">
export default { name: 'PublishLog' }
</script>
<script lang="ts" setup>
import pms from '@/permission'
import Button from '@/components/Permission/Button.vue'
import { Search } from '@element-plus/icons-vue'
import TheDetailDialog from '@/views/deploy/TheDetailDialog.vue'
import { PublishLogData, PublishLogList, PublishLogTotal } from '@/api/log'
import { ref } from 'vue'
import { ProjectData } from '@/api/project'

const dialogVisible = ref(false)
const searchParam = ref({ username: '', projectName: '' })
const tableLoading = ref(false)
const tableData = ref<PublishLogList['datagram']['list']>([])
const pagination = ref({ page: 1, rows: 20, total: 0 })
const projectItem = ref({} as ProjectData)

getList()
getTotal()

function searchList() {
  pagination.value.page = 1
  getList()
  getTotal()
}
function getList() {
  tableLoading.value = true
  tableData.value = []
  new PublishLogList(searchParam.value, pagination.value)
    .request()
    .then((response) => {
      tableData.value = response.data.list
    })
    .finally(() => {
      tableLoading.value = false
    })
}
function getTotal() {
  new PublishLogTotal(searchParam.value).request().then((response) => {
    pagination.value.total = response.data.total
  })
}
function handlePageChange(val = 1) {
  pagination.value.page = val
  getList()
}
function handleDetail(data: PublishLogData) {
  projectItem.value.id = data.projectId
  projectItem.value.lastPublishToken = data.token
  projectItem.value.deployState = data.state === 1 ? 2 : 3
  projectItem.value.repoType = 'repository'
  projectItem.value.transferType = 'transfer'
  dialogVisible.value = true
}
</script>

<style scoped lang="scss">
@import '@/styles/mixin.scss';
.icon-success {
  color: #67c23a;
  font-size: 15px;
  &::before {
    content: '\2713';
  }
}

.icon-fail {
  color: #f56c6c;
  font-size: 15px;
  &::before {
    content: '\2717';
  }
}
.project {
  &-detail {
    padding-left: 5px;
    height: 470px;
    overflow-y: auto;
    @include scrollBar();
  }
  &-title {
    display: flex;
    flex-direction: row;
    width: 100%;
    &:before,
    &:after {
      content: '';
      flex: 1 1;
      border-bottom: 1px solid;
      margin: auto;
    }
    &:before {
      margin-right: 10px;
    }
    &:after {
      margin-left: 10px;
    }
  }
}
</style>
