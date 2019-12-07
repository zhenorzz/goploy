<template>
  <el-row class="tree-container">
    <el-col :span="10" class="input-container">
      <div ref="jsonStringInput" class="json-string-input" contenteditable="true" placeholder="在此输入json字符串" @input="handleInput" />
    </el-col>
    <el-col :span="14" class="json-container">
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
      tree: undefined
    }
  },
  computed: {

  },
  created() {

  },
  mounted() {
    this.$refs.jsonStringInput.focus()
    // Get DOM-element for inserting json-tree

    // var wrapper = this.$refs.jsonPrettyString

    // var tree = jsonTree.create({}, wrapper)
    // tree.loadData
    // // Expand all (or selected) child nodes of root (optional)
    // tree.expand()
  },
  methods: {
    handleInput: debounce(function(event) {
      const wrapper = this.$refs.jsonPrettyString
      const text = event.target.innerText.replace(/<[^>]+>/g, '')
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
    })
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
    padding: 20px 30px;
    height: 100%;
    &:empty:before{
      content: attr(placeholder);
      color:#bbb;
    }
    &:focus:before{
      content:none;
    }
  }
  .json-container {
    height: 100%;
    border-left: solid 1px #f6f6f6;
  }
}

</style>
