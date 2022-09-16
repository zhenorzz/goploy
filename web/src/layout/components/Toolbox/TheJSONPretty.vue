<template>
  <el-row v-show="modelValue === 'json'" style="width: 100%">
    <el-row
      class="json-helper"
      type="flex"
      justify="space-between"
      align="middle"
    >
      <el-row>
        <el-button type="primary" text @click="switchJSONView">
          <svg-icon icon-class="switch" /> JSON Pretty
        </el-button>
        <el-button
          v-show="json.view === 'pretty'"
          type="primary"
          text
          @click="expandAll"
        >
          {{ $t('JSONPage.expandAll') }}
        </el-button>
        <el-button
          v-show="json.view === 'pretty'"
          type="primary"
          text
          @click="collapseAll"
        >
          {{ $t('JSONPage.collapseAll') }}
        </el-button>
        <el-button
          v-show="json.view === 'pretty'"
          type="primary"
          text
          @click="unmarkAll"
        >
          {{ $t('JSONPage.unmarkAll') }}
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
    <el-input
      v-show="json.view === 'raw'"
      v-model="json.input"
      style="width: 100%"
      :autosize="{ minRows: 25, maxRows: 25 }"
      type="textarea"
      class="json-string-input"
      placeholder="JSON string"
      @input="handleInput"
      @paste="onPaste"
    />
    <div
      v-show="json.view === 'pretty'"
      ref="jsonPrettyString"
      class="json-pretty-string"
    />
  </el-row>
</template>

<script lang="ts" setup>
import '@/components/JSONTree/index.css'
import { jsonTree } from '@/components/JSONTree'
import { QuestionFilled } from '@element-plus/icons-vue'
import { ref, reactive } from 'vue'
defineProps({
  modelValue: {
    type: String,
    default: '',
  },
})
const json = reactive({
  view: 'raw',
  input: '',
  tree: {} as any,
})
const jsonPrettyString = ref()
const handleInput = () => {
  const text = json.input
  jsonPrettyString.value.innerText = ''
  try {
    const data = JSON.parse(text)
    json.tree = jsonTree.create(data, jsonPrettyString.value)
    json.tree.expand()
  } catch (error) {
    jsonPrettyString.value.innerText = error.message
  }
}
const onPaste = (e: ClipboardEvent) => {
  json.view = 'pretty'
  const clip = e.clipboardData ? e.clipboardData.getData('Text') : ''
  json.input = clip
  handleInput()
  return true
}

const switchJSONView = () => {
  if (json.view === 'raw') {
    json.view = 'pretty'
  } else {
    json.view = 'raw'
  }
}
const expandAll = () => {
  json.tree && json.tree.expand()
}

const collapseAll = () => {
  json.tree && json.tree.collapse()
}

const unmarkAll = () => {
  json.tree && json.tree.unmarkAll()
}
</script>
<style lang="scss" scoped>
.json-helper {
  padding: 0 10px;
  height: 35px;
  width: 100%;
}
</style>
