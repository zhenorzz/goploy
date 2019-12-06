<template>
  <el-row id="wrapper" class="app-container">
    json
  </el-row>
</template>
<script>
import './jsonTree.css'
import { jsonTree } from './jsonTree'
export default {

  data() {
    return {
    }
  },
  computed: {

  },
  created() {

  },
  mounted() {
    // Get DOM-element for inserting json-tree
    var wrapper = document.getElementById('wrapper')

    // Get json-data by javascript-object
    var data = {
      'firstName': 'John',
      'lastName': 'Smith',
      'isAlive': true,
      'age': 25,
      'company': null,
      'height_cm': 167.64,
      'address': {
        'streetAddress': '21 2nd Street',
        'city': 'New York',
        'state': 'NY',
        'postalCode': '10021-3100'
      },
      'phoneNumbers': [
        {
          'type': 'home',
          'number': '212 555-1234'
        },
        {
          'type': 'fax',
          'number': '646 555-4567'
        }
      ]
    }

    var tree = jsonTree.create(data, wrapper)

    // Expand all (or selected) child nodes of root (optional)
    tree.expand()
  },
  methods: {
  }
}
</script>
