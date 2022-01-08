<template>
  <el-dialog
    v-model="dialogVisible"
    :title="$t('deployPage.processManager')"
    :fullscreen="$store.state.app.device === 'mobile'"
    :close-on-click-modal="false"
  >
    <el-row type="flex">
      <el-select
        v-model="name"
        v-loading="processLoading"
        filterable
        style="flex: 1"
        @change="handleProcessChange"
      >
        <el-option
          v-for="item in processOption"
          :key="item"
          :label="item"
          :value="item"
        />
      </el-select>
      <el-button type="primary" icon="el-icon-plus" @click="handleAdd" />
    </el-row>
    <el-table
      v-loading="tableLoading"
      border
      stripe
      highlight-current-row
      max-height="447px"
      style="margin-top: 10px; width: 100%"
      :data="tableData"
    >
      <el-table-column prop="ip" label="IP" />
      <el-table-column
        prop="status"
        label="Status"
        width="200"
        show-overflow-tooltip
      />
      <el-table-column
        prop="operation"
        :label="$t('op')"
        width="80"
        align="center"
        :fixed="$store.state.app.device === 'mobile' ? false : 'right'"
      >
        <template #default="scope">
          <el-button
            v-if="scope.row.isModified"
            type="primary"
            @click="handleDiff(scope.row)"
          >
            diff
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    <template #footer class="dialog-footer">
      <el-button @click="dialogVisible = false">
        {{ $t('cancel') }}
      </el-button>
    </template>
  </el-dialog>
  <el-dialog
    v-model="fileVisible"
    title="File"
    :fullscreen="$store.state.app.device === 'mobile'"
  >
    <el-table
      v-loading="fileLoading"
      border
      stripe
      highlight-current-row
      :data="fileList"
      style="width: 100%"
      :max-height="460"
    >
      <el-table-column prop="name" :label="$t('name')" min-width="100">
        <template #default="scope">
          <i v-if="scope.row.isDir" class="el-icon-folder-opened"></i>
          {{ scope.row.name }}
        </template>
      </el-table-column>
      <el-table-column prop="size" :label="$t('size')" width="100">
        <template #default="scope">
          {{ humanSize(scope.row.size) }}
        </template>
      </el-table-column>
      <el-table-column prop="mode" label="mode" width="100" />
      <el-table-column
        prop="modTime"
        :label="$t('modifiedTime')"
        width="135"
        align="center"
      />
      <el-table-column
        prop="operation"
        :label="$t('op')"
        width="100"
        align="center"
        :fixed="$store.state.app.device === 'mobile' ? false : 'right'"
      >
        <template #default="scope">
          <template v-if="scope.row.uploading">
            <i class="el-icon-loading"></i>
          </template>
          <template v-else>
            <el-button
              v-if="scope.row.isDir"
              style="margin-right: 10px"
              type="text"
              icon="el-icon-right"
              @click="handleSelectPath(filePath + scope.row.name + '/')"
            />
            <el-button
              v-else
              style="margin-right: 10px"
              type="text"
              icon="el-icon-check"
              @click="handleSelectFile(filePath + scope.row.name)"
            />
          </template>
        </template>
      </el-table-column>
    </el-table>
    <template #footer class="dialog-footer">
      <el-button @click="fileVisible = false">
        {{ $t('cancel') }}
      </el-button>
    </template>
  </el-dialog>
  <el-dialog
    v-model="fileDiffVisible"
    title="File diff"
    :fullscreen="$store.state.app.device === 'mobile'"
  >
    <div class="file-content">
      <div v-loading="diffLoading">
        <div
          v-for="(item, index) in changeLines"
          :key="index"
          class="file-line"
        >
          <div
            class="file-line-number"
            :style="{
              'background-color': item.type
                ? item.type === '+'
                  ? '#ccFFD8'
                  : '#ffd7d5'
                : '',
            }"
          >
            {{ item.lineNumber }}
          </div>
          <div
            class="file-line-type"
            :style="{
              'background-color': item.type
                ? item.type === '+'
                  ? '#e6ffec'
                  : '#ffebe9'
                : '',
            }"
          >
            {{ item.type }}
          </div>
          <div
            class="file-line-text"
            :style="{
              'background-color': item.type
                ? item.type === '+'
                  ? '#e6ffec'
                  : '#ffebe9'
                : '',
            }"
          >
            {{ item.text }}
          </div>
        </div>
      </div>
    </div>
    <template #footer class="dialog-footer">
      <el-button @click="fileDiffVisible = false">
        {{ $t('cancel') }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script lang="ts">
import { DeployProcessList } from '@/api/deploy'
import { ElMessage } from 'element-plus'
import { computed, defineComponent, ref, watch } from 'vue'

export default defineComponent({
  props: {
    modelValue: {
      type: Boolean,
      default: false,
    },
    projectRow: {
      type: Object,
      required: true,
    },
  },
  emits: ['update:modelValue'],
  setup(props, { emit }) {
    const dialogVisible = computed({
      get: () => props.modelValue,
      set: (val) => {
        emit('update:modelValue', val)
      },
    })

    watch(
      () => props.modelValue,
      (val: typeof props['modelValue']) => {
        if (val === true) {
          getProcessList()
        }
      }
    )

    const processLoading = ref(false)
    const processOption = ref<DeployProcessList['datagram']['list']>([])
    const getProcessList = () => {
      processLoading.value = true
      new DeployProcessList({ projectId: props.projectRow.id }, {})
        .request()
        .then((response) => {
          processOption.value = response.data.list
        })
        .finally(() => {
          processLoading.value = false
        })
    }
    const tableLoading = ref(false)
    const tableData = ref<DeployProcessList['datagram']['list']>([])

    return {
      dialogVisible,
      processOption,
      processLoading,
      tableLoading,
      tableData,
    }
  },
})
</script>

<style rel="stylesheet/scss" lang="scss" scoped>
@import '@/styles/mixin.scss';
.file {
  &-content {
    height: 470px;
    overflow-y: auto;
    @include scrollBar();
  }
  &-line {
    display: flex;
    flex-direction: row;
    font-size: 12px;
    line-height: 20px;
  }
  &-line-number {
    width: 50px;
    text-align: right;
    padding-right: 10px;
    padding-left: 10px;
    color: #6e7781;
  }
  &-line-type {
    width: 22px;
    text-align: center;
  }
  &-line-text {
    flex: 1;
    white-space: pre-wrap;
  }
}
</style>
