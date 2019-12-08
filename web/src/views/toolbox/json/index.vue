<template>
  <el-row class="tree-container">
    <el-col :span="10" class="input-container">
      <el-input ref="jsonStringInput" v-model="inputContent" :autosize="false" type="textarea" class="json-string-input" contenteditable="true" placeholder="在此输入json字符串" @input="handleInput" />
    </el-col>
    <el-col :span="14" class="json-container">
      <el-row class="json-helper" type="flex" justify="space-between" align="middle">
        <el-row>
          <el-button type="text" size="medium" @click="expandAll">展开所有</el-button>
          <el-button type="text" size="medium" @click="collapseAll">收起所有</el-button>
          <el-button type="text" size="medium" @click="unmarkAll">取消高亮</el-button>
        </el-row>
        <el-row>
          <el-tooltip class="item" effect="dark" placement="bottom-end">
            <el-button type="text" icon="el-icon-question" />
            <div slot="content">
              1.按住ALT点击label可以实现高亮<br>
              2.按住SHIFT可以查看JSON路径
            </div>
          </el-tooltip>
        </el-row>
      </el-row>
      <div ref="jsonPrettyString" />
    </el-col>
  </el-row>
</template>
<script>
import './jsonTree.css'
import { jsonTree } from './jsonTree'
import { debounce } from '@/utils'
export default {

  data() {
    return {
      inputContent: '',
      tree: undefined
    }
  },
  computed: {

  },
  created() {

  },
  mounted() {
    this.$refs.jsonStringInput.focus()
  },
  methods: {
    handleInput: debounce(function() {
      const wrapper = this.$refs.jsonPrettyString
      const text = this.inputContent
      if (text.length === 0) {
        wrapper.innerText = ''
        return
      }
      try {
        const data = JSON.parse(text)

        if (this.tree) {
          this.tree.loadData(data)
        } else {
          this.tree = jsonTree.create(data, wrapper)
        }
        this.tree.expand()
      } catch (error) {
        wrapper.innerText = error.message
      }
    }),

    expandAll() {
      this.tree && this.tree.expand()
    },

    collapseAll() {
      this.tree && this.tree.collapse()
    },

    unmarkAll() {
      this.tree && this.tree.unmarkAll()
    }
  }
}
</script>
<style lang="scss" scoped>
.tree-container {
  height: calc(100vh - 50px);
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
    border-bottom: solid 1px #EBEEF5;
  }

  .json-container {
    height: 100%;
    border-left: solid 1px #EBEEF5;
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
