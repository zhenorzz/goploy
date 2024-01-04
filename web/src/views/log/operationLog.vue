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
          v-model="searchParam.router"
          style="width: 200px"
          placeholder="Filter the router"
        />
        <el-input
          v-model="searchParam.api"
          style="width: 200px"
          placeholder="Filter the api"
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
        stripe
        highlight-current-row
        height="100%"
        :data="tableData"
      >
        <el-table-column prop="id" label="ID" width="100" />
        <el-table-column prop="username" label="Username" width="150" />
        <el-table-column prop="routerName" label="Router" min-width="230">
          <template #default="scope">
            {{ scope.row.routerName }}({{ scope.row.router }})
          </template>
        </el-table-column>
        <el-table-column prop="api" label="API" min-width="200" />
        <el-table-column label="Request Data" width="160" align="center">
          <template #default="scope">
            <el-button
              type="primary"
              text
              @click="showJSON(scope.row.requestData)"
            >
              {{ $t('detail') }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column label="Response Data" width="160" align="center">
          <template #default="scope">
            <el-button
              type="primary"
              text
              @click="showJSON(scope.row.responseData)"
            >
              {{ $t('detail') }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column
          prop="requestTime"
          label="Request time"
          width="160"
          align="center"
        />
        <el-table-column
          prop="responseTime"
          label="Response time"
          width="160"
          align="center"
        />
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
    <el-dialog v-model="json.dialog" :width="'70%'">
      <el-row
        class="json-helper"
        type="flex"
        justify="space-between"
        align="middle"
      >
        <el-row>
          <el-button type="primary" text @click="expandAll">
            {{ $t('JSONPage.expandAll') }}
          </el-button>
          <el-button type="primary" text @click="collapseAll">
            {{ $t('JSONPage.collapseAll') }}
          </el-button>
          <el-button type="primary" text @click="unmarkAll">
            {{ $t('JSONPage.unmarkAll') }}
          </el-button>
          <el-button type="primary" text @click="copyAll">
            {{ $t('JSONPage.copyAll') }}
          </el-button>
        </el-row>
        <el-row>
          <el-tooltip class="item" effect="dark" placement="bottom-end">
            <el-button type="primary" text :icon="QuestionFilled" />
            <template #content>
              <span style="white-space: pre-line">
                {{ $t('JSONPage.tips') }}
              </span>
            </template>
          </el-tooltip>
        </el-row>
      </el-row>
      <div ref="jsonPrettyString" v-loading="json.loading" />
    </el-dialog>
  </el-row>
</template>

<script lang="ts">
export default { name: 'OperationLog' }
</script>
<script lang="ts" setup>
import '@/components/JSONTree/index.css'
import { jsonTree } from '@/components/JSONTree'
import { QuestionFilled, Search } from '@element-plus/icons-vue'
import { copy } from '@/utils'
import { OperationLogList, OperationLogTotal } from '@/api/log'
import { ref, reactive, nextTick } from 'vue'
import { useStore } from 'vuex'
import { useI18n } from 'vue-i18n'
const store = useStore()
const { t } = useI18n()
const searchParam = ref({ username: '', router: '', api: '' })
const tableLoading = ref(false)
const tableData = ref<OperationLogList['datagram']['list']>([])
const pagination = ref({ page: 1, rows: 20, total: 0 })
const jsonPrettyString = ref()
const json = reactive({
  str: '',
  tree: {} as any,
  dialog: false,
  loading: false,
})
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
  new OperationLogList(searchParam.value, pagination.value)
    .request()
    .then((response) => {
      const routerObj: any = {}
      for (const father of store.state.permission.routes) {
        try {
          for (const children of father.children) {
            routerObj[father.path + '/' + children.path] = t(
              'route.' + children.meta.title
            )
          }
        } catch (error) {
          console.log(error)
        }
      }
      tableData.value = response.data.list.map((item) => {
        item.routerName = routerObj[item.router] || item.router
        return item
      })
    })
    .finally(() => {
      tableLoading.value = false
    })
}
function getTotal() {
  new OperationLogTotal(searchParam.value).request().then((response) => {
    pagination.value.total = response.data.total
  })
}

function handlePageChange(val = 1) {
  pagination.value.page = val
  getList()
}

function showJSON(jsonStr: string) {
  json.str = jsonStr
  json.dialog = true
  json.loading = true
  nextTick(() => {
    jsonPrettyString.value.innerText = ''
    try {
      const data = JSON.parse(jsonStr)
      json.tree = jsonTree.create(data, jsonPrettyString.value)
    } catch (error) {
      jsonPrettyString.value.innerText = error.message
    }
    json.loading = false
  })
}

function expandAll() {
  json.tree.expand()
}

function collapseAll() {
  json.tree.collapse()
}
function unmarkAll() {
  json.tree.unmarkAll()
}

function copyAll() {
  copy(json.tree.toSourceJSON())
  ElMessage.success('Success')
}
</script>
