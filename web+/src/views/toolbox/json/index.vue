<template>
  <el-row class="tree-container">
    <el-col :span="10" class="input-container">
      <el-input
        ref="jsonStringInput"
        v-model="inputContent"
        :autosize="false"
        type="textarea"
        class="json-string-input"
        placeholder="JSON string"
        contenteditable="true"
        @input="handleInput"
      />
    </el-col>
    <el-col :span="14" class="json-container">
      <el-row
        class="json-helper"
        type="flex"
        justify="space-between"
        align="middle"
      >
        <el-row>
          <el-button type="text" size="medium" @click="expandAll">
            {{ $t('JSONPage.expandAll') }}
          </el-button>
          <el-button type="text" size="medium" @click="collapseAll">
            {{ $t('JSONPage.collapseAll') }}
          </el-button>
          <el-button type="text" size="medium" @click="unmarkAll">
            {{ $t('JSONPage.unmarkAll') }}
          </el-button>
        </el-row>
        <el-row>
          <el-tooltip class="item" effect="dark" placement="bottom-end">
            <el-button type="text" icon="el-icon-question" />
            <template #content>
              <span v-html="$t('JSONPage.tips')"></span>
            </template>
          </el-tooltip>
        </el-row>
      </el-row>
      <div ref="jsonPrettyString" class="json-pretty-string" />
    </el-col>
  </el-row>
</template>
<script>
import './jsonTree.css'
import { jsonTree } from './jsonTree'
import { defineComponent } from 'vue'
export default defineComponent({
  name: 'JSONFormatter',
  data() {
    return {
      inputContent: '',
      tree: undefined,
    }
  },
  computed: {},
  mounted() {
    this.$refs.jsonStringInput.focus()
  },
  methods: {
    handleInput() {
      const wrapper = this.$refs.jsonPrettyString
      const text = this.inputContent
      wrapper.innerText = ''
      try {
        const data = JSON.parse(text)
        this.tree = jsonTree.create(data, wrapper)
        this.tree.expand()
      } catch (error) {
        wrapper.innerText = error.message
      }
    },

    expandAll() {
      this.tree && this.tree.expand()
    },

    collapseAll() {
      this.tree && this.tree.collapse()
    },

    unmarkAll() {
      this.tree && this.tree.unmarkAll()
    },
  },
})
</script>
<style lang="scss" scoped>
@import '@/styles/mixin.scss';
.tree-container {
  height: calc(100vh - 90px);
  .input-container {
    height: 100%;
  }
  .json-string-input {
    height: 100%;
    textarea {
      font-size: 16px;
      border-radius: 0;
      padding: 20px 30px;
      height: 100%;
      resize: none;
    }
  }

  .json-helper {
    padding: 0 10px;
    height: 35px;
    border-bottom: solid 1px #ebeef5;
  }

  .json-container {
    height: 100%;
    border-left: solid 1px #ebeef5;
    .json-pretty-string {
      overflow: auto;
      height: 100%;
      @include scrollBar();
    }
  }
}
</style>
<style lang="scss">
.json-string-input {
  .el-textarea__inner {
    font-size: 16px;
    border-radius: 0;
    padding: 20px 30px;
    height: 100%;
    resize: none;
  }
}
</style>
